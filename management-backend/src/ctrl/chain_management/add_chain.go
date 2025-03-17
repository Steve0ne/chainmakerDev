/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"management_backend/src/db/connection"
	"strings"

	"github.com/gin-gonic/gin"

	"chainmaker.org/chainmaker/pb-go/v2/consensus"

	"management_backend/src/ctrl/common"
	dbchain "management_backend/src/db/chain"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/relation"
	"management_backend/src/entity"
)

type AddChainHandler struct{}

func (addChainHandler *AddChainHandler) LoginVerify() bool {
	return true
}

func (addChainHandler *AddChainHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindAddChainHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	_, err := dbchain.GetChainByChainIdOrName(params.ChainId, params.ChainName)
	if err == nil {
		log.Error("Chain has existed")
		common.ConvergeFailureResponse(ctx, common.ErrorChainExisted)
		return
	}

	chain := &dbcommon.Chain{
		ChainId:           params.ChainId,
		ChainName:         params.ChainName,
		Consensus:         consensus.ConsensusType_name[params.Consensus],
		TxTimeout:         params.TxTimeout,
		BlockTxCapacity:   params.BlockTxCapacity,
		BlockInterval:     params.BlockInterval,
		Status:            connection.NO_START,
		Monitor:           params.Monitor,
		ChainmakerImprove: params.ChainmakerImprove,
		AutoReport:        params.ChainmakerImprove,
		Address:           params.Address,
		TLS:               params.Tls,
		DockerVm:          params.DockerVm,
	}
	err = dbchain.CreateChain(chain)
	if err != nil {
		log.Error("CreateChain err : " + err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorCreateChain)
		return
	}

	for _, node := range params.Nodes {
		org, err := chain_participant.GetOrgByOrgId(node.OrgId)
		if err != nil {
			common.ConvergeFailureResponse(ctx, common.ErrorGetOrg)
			return
		}
		chainOrg := &dbcommon.ChainOrg{
			ChainId: chain.ChainId,
			OrgId:   org.OrgId,
			OrgName: org.OrgName,
		}
		err = relation.CreateChainOrg(chainOrg)
		if err != nil {
			log.Error("CreateChainOrg err : " + err.Error())
			common.ConvergeFailureResponse(ctx, common.ErrorCreateChainOrg)
			return
		}
		nodeNames := strings.Split(node.NodeList, ",")
		for _, nodeName := range nodeNames {
			nodeInfo, err := chain_participant.GetNodeByNodeName(nodeName)
			if err != nil {
				common.ConvergeFailureResponse(ctx, common.ErrorGetNode)
				return
			}
			chainOrgNode := &dbcommon.ChainOrgNode{
				ChainId:  chain.ChainId,
				OrgId:    org.OrgId,
				OrgName:  org.OrgName,
				NodeId:   nodeInfo.NodeId,
				NodeName: nodeInfo.NodeName,
			}
			err = relation.CreateChainOrgNode(chainOrgNode)
			if err != nil {
				log.Error("CreateChainOrgNode err : " + err.Error())
				common.ConvergeFailureResponse(ctx, common.ErrorCreateChainOrgNode)
				return
			}
		}
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
