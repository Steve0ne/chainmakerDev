/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_invoke

import (
	"chainmaker.org/chainmaker/common/v2/evmutils"
	"chainmaker.org/chainmaker/common/v2/random/uuid"
	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"management_backend/src/ctrl/contract_management"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/contract"
	"management_backend/src/entity"
	loggers "management_backend/src/logger"
	"management_backend/src/sync"
)

const (
	INVOKE_FAIL    = 2
	INVOKE_SUCCESS = 1
)

const (
	TxStatusCode_SUCCESS = 0
	TxStatusCode_FAIL    = 1
)

const contractFail = 1

var log = loggers.GetLogger(loggers.ModuleWeb)

type InvokeContractHandler struct{}

func (invokeContractHandler *InvokeContractHandler) LoginVerify() bool {
	return true
}

func (invokeContractHandler *InvokeContractHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindInvokeContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	chainId := params.ChainId
	sdkClientPool := sync.GetSdkClientPool()
	if sdkClientPool == nil {
		common.ConvergeFailureResponse(ctx, common.ErrorChainNotSub)
		return
	}

	txId := uuid.GetUUID() + uuid.GetUUID()
	sdkClient := sdkClientPool.SdkClients[chainId]
	if sdkClient == nil {
		common.ConvergeFailureResponse(ctx, common.ErrorChainNotSub)
		return
	}

	contractInfo, err := contract.GetContractByName(params.ChainId, params.ContractName)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorContractNotExist)
		return
	}

	var kvPair []*pbcommon.KeyValuePair
	var methodName string
	var contractName string

	if contractInfo.RuntimeType == contract_management.EVM {
		if len(contractInfo.EvmAbiSaveKey) < 1 {
			common.ConvergeFailureResponse(ctx, common.ErrorContractNotExist)
		}
		kvPair, methodName, err = GetEvmKv(contractInfo.EvmAbiSaveKey, params.MethodName, params.Parameters, sdkClient.SdkConfig.UserCert)
		contractName = hex.EncodeToString(evmutils.Keccak256([]byte(params.ContractName)))[24:]
	} else {
		kvPair = convertToPbKeyValues(params)
		methodName = params.MethodName
		contractName = params.ContractName
	}

	resp, err := sdkClient.ChainClient.InvokeContract(contractName,
		methodName, txId, kvPair, -1, true)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorInvokeContract)
		return
	}

	var status = INVOKE_SUCCESS
	var txStatus = TxStatusCode_SUCCESS
	if resp.Code != pbcommon.TxStatusCode_SUCCESS {
		status = INVOKE_FAIL
		log.Infof("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}

	if resp.ContractResult.Code == contractFail {
		txStatus = TxStatusCode_FAIL
	}

	orgName, err := chain_participant.GetOrgNameByOrgId(sdkClient.SdkConfig.OrgId)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorGetOrgName)
		return
	}

	invokeRecords := &dbcommon.InvokeRecords{
		ChainId:      params.ChainId,
		OrgId:        sdkClient.SdkConfig.OrgId,
		OrgName:      orgName,
		ContractName: params.ContractName,
		TxId:         txId,
		TxStatus:     txStatus,
		Status:       status,
		UserName:     sdkClient.SdkConfig.UserName,
		Method:       params.MethodName,
	}
	err = contract.CreateInvokeRecords(invokeRecords)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorCreateRecordFailed)
		return
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

func convertToPbKeyValues(body *InvokeContractListParams) []*pbcommon.KeyValuePair {
	keyValues := body.Parameters
	if len(keyValues) > 0 {
		pbKvs := make([]*pbcommon.KeyValuePair, 0)
		for _, kv := range keyValues {
			pbKvs = append(pbKvs, &pbcommon.KeyValuePair{
				Key:   kv.Key,
				Value: []byte(kv.Value),
			})
		}
		return pbKvs
	}
	return []*pbcommon.KeyValuePair{}
}

type ReInvokeContractHandler struct{}

func (reInvokeContractHandler *ReInvokeContractHandler) LoginVerify() bool {
	return true
}

func (reInvokeContractHandler *ReInvokeContractHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindReInvokeContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	invokeRecord, err := contract.GetInvokeRecords(params.InvokeRecordId)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorQueryInvokeRecord)
		return
	}
	if invokeRecord.Status == INVOKE_SUCCESS {
		common.ConvergeFailureResponse(ctx, common.ErrorAlreadyOnChain)
		return
	}

	sdkClientPool := sync.GetSdkClientPool()
	if sdkClientPool == nil {
		common.ConvergeFailureResponse(ctx, common.ErrorChainNotSub)
		return
	}
	txId := uuid.GetUUID() + uuid.GetUUID()
	sdkClient := sdkClientPool.SdkClients[invokeRecord.ChainId]
	resp, err := sdkClient.ChainClient.InvokeContract(invokeRecord.ContractName,
		invokeRecord.Method, txId, nil, -1, true)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorInvokeContract)
		return
	}
	if resp.Code == pbcommon.TxStatusCode_SUCCESS {
		invokeRecord.TxStatus = int(resp.ContractResult.Code)
		invokeRecord.Status = INVOKE_SUCCESS
		invokeRecord.TxId = txId
		err := contract.UpdateInvokeRecordsStatus(invokeRecord)
		if err != nil {
			common.ConvergeFailureResponse(ctx, common.ErrorUpdateRecordFailed)
			return
		}
	} else {
		log.Infof("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
