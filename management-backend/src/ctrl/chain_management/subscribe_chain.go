/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"github.com/gin-gonic/gin"
	"management_backend/src/ctrl/ca"
	"management_backend/src/ctrl/log_report"
	"management_backend/src/db/chain"
	"management_backend/src/db/chain_participant"
	"management_backend/src/sync"
	"strings"

	"management_backend/src/ctrl/common"
	"management_backend/src/entity"
)

const (
	CONNECT_ERR = "all client connections are busy"
	AUTH_ERR    = "authentication error"
	TLS_ERR     = "handshake failure"
	CHAIN_ERR   = "not found"
)

type SubscribeChainHandler struct{}

func (subscribeChainHandler *SubscribeChainHandler) LoginVerify() bool {
	return false
}

func (subscribeChainHandler *SubscribeChainHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindSubscribeChainHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	nodeInfo, err := chain_participant.GetNodeByNodeName(params.NodeName)
	if err != nil {
		log.Error("GetNodeByNodeName err : " + err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorGetNode)
		return
	}

	orgCa, err := chain_participant.GetOrgCaCert(params.OrgId)
	if err != nil {
		log.Error("GetOrgCaCert err : " + err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorGetOrgCaCert)
		return
	}

	userInfo, err := chain_participant.GetUserTlsCert(params.UserName)
	if err != nil {
		log.Error("GetUserTlsCert err : " + err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorGetUserTlsCert)
		return
	}

	Tls := true
	if params.Tls == NO_TLS {
		Tls = false
	}

	sdkConfig := &entity.SdkConfig{
		ChainId:     params.ChainId,
		OrgId:       params.OrgId,
		UserName:    params.UserName,
		Tls:         Tls,
		TlsHost:     ca.TLS_HOST,
		Remote:      nodeInfo.NodeIp + ":" + nodeInfo.NodePort,
		CaCert:      []byte(orgCa.Cert),
		UserCert:    []byte(userInfo.Cert),
		UserPrivKey: []byte(userInfo.PrivateKey),
	}

	err = sync.SubscribeChain(sdkConfig)
	if err != nil {
		log.Error("SubscribeChain err : " + err.Error())
		if strings.Contains(err.Error(), CONNECT_ERR) {
			common.ConvergeFailureResponse(ctx, common.ErrorSubscribeChainConnectNode)
			return
		}
		if strings.Contains(err.Error(), AUTH_ERR) {
			common.ConvergeFailureResponse(ctx, common.ErrorSubscribeChainCert)
			return
		}
		if strings.Contains(err.Error(), TLS_ERR) {
			common.ConvergeFailureResponse(ctx, common.ErrorSubscribeChainTls)
			return
		}
		if strings.Contains(err.Error(), CHAIN_ERR) {
			common.ConvergeFailureResponse(ctx, common.ErrorSubscribeChainId)
			return
		}
		common.ConvergeFailureResponse(ctx, common.ErrorSubscribeChain)
		return
	}

	chainInfo, err := chain.GetChainByChainId(params.ChainId)
	if err != nil {
		log.Error("GetChainInfoByChainId err : " + err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorGetChain)
		return
	}

	if chainInfo.AutoReport == log_report.AUTO {
		tickerMap := log_report.TickerMap
		_, ok := tickerMap[params.ChainId]
		if !ok {
			err := sync.ReportChainData(params.ChainId)
			if err != nil {
				log.Error(err)
			}
			ticker := log_report.NewTicker(24)
			ticker.Start(params.ChainId)
		}
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
