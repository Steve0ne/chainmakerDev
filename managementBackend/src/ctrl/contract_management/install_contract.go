/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_management

import (
	"chainmaker.org/chainmaker/common/v2/evmutils"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"

	sdkcommon "chainmaker.org/chainmaker/pb-go/v2/common"

	"management_backend/src/ctrl/common"
	dbcommon "management_backend/src/db/common"
	dbcontract "management_backend/src/db/contract"
	"management_backend/src/entity"
)

type InstallContractHandler struct{}

func (installContractHandler *InstallContractHandler) LoginVerify() bool {
	return true
}

func (installContractHandler *InstallContractHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindInstallContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	_, err := dbcontract.GetContractByName(params.ChainId, params.ContractName)
	if err == nil {
		common.ConvergeFailureResponse(ctx, common.ErrorContractExist)
		return
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorMarshalParameters)
		return
	}

	orgId, err := SaveVote(params.ChainId, params.Reason, string(jsonBytes), INIT_CONTRACT, UPDATE_OTHER)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}

	paramJson, err := json.Marshal(params.Parameters)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorMarshalParameters)
		return
	}
	paramStr := string(paramJson)
	if paramStr == NULL {
		paramStr = ""
	}

	var methodStr string
	var functionType int
	if params.RuntimeType == EVM {
		methodStr, functionType, err = GetEvmMethodsByAbi(params.EvmAbiSaveKey)
		if err != nil {
			common.ConvergeFailureResponse(ctx, common.ErrorAbiMethods)
			return
		}
	} else {
		methodJson, err := json.Marshal(params.Methods)
		if err != nil {
			common.ConvergeFailureResponse(ctx, common.ErrorMarshalMethods)
			return
		}

		methodStr = string(methodJson)
		if methodStr == NULL {
			methodStr = ""
		}
	}

	contract := &dbcommon.Contract{
		ChainId:          params.ChainId,
		Name:             params.ContractName,
		Version:          params.ContractVersion,
		RuntimeType:      params.RuntimeType,
		SourceSaveKey:    params.CompileSaveKey,
		EvmAbiSaveKey:    params.EvmAbiSaveKey,
		EvmAddress:       hex.EncodeToString(evmutils.Keccak256([]byte(params.ContractName)))[24:],
		EvmFunctionType:  functionType,
		ContractOperator: user.Name,
		MgmtParams:       paramStr,
		Methods:          methodStr,
		ContractStatus:   int(dbcommon.ContractInitStored),
		MultiSignStatus:  int(sdkcommon.MultiSignStatus_PROCESSING),
		OrgId:            orgId,
		Timestamp:        time.Now().Unix(),
	}
	err = dbcontract.CreateContract(contract)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorInstallContract)
		return
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
