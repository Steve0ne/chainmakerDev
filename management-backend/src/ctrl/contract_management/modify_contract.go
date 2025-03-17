/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_management

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/contract"
	"management_backend/src/entity"
)

type ModifyContractHandler struct{}

func (modifyContractHandler *ModifyContractHandler) LoginVerify() bool {
	return true
}

func (modifyContractHandler *ModifyContractHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindModifyContractParamsHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	methodJson, err := json.Marshal(params.Methods)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorMarshalMethods)
		return
	}

	methodStr := string(methodJson)
	if methodStr == NULL {
		methodStr = ""
	}
	var functionType int
	if len(params.EvmAbiSaveKey) > 0 {
		methodStr, functionType, err = GetEvmMethodsByAbi(params.EvmAbiSaveKey)
		if err != nil {
			common.ConvergeFailureResponse(ctx, common.ErrorAbiMethods)
			return
		}
	}

	contractInfo := &dbcommon.Contract{
		Methods:         methodStr,
		EvmAbiSaveKey:   params.EvmAbiSaveKey,
		EvmFunctionType: functionType,
	}

	contractInfo.Id = params.Id

	err = contract.UpdateContractMethod(contractInfo)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorUpdateMethod)
		return
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
