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
	dbcontract "management_backend/src/db/contract"
	"management_backend/src/entity"
)

type UnFreezeContractHandler struct{}

func (unFreezeContractHandler *UnFreezeContractHandler) LoginVerify() bool {
	return true
}

func (unFreezeContractHandler *UnFreezeContractHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindUnFreezeContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	contract, err := dbcontract.GetContractByName(params.ChainId, params.ContractName)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorGetContract)
		return
	}

	if contract.MultiSignStatus == dbcommon.VOTING {
		common.ConvergeFailureResponse(ctx, common.ErrorContractBeingVoting)
		return
	}

	if !contract.CanUnfreeze() {
		// 不可以进行解冻操作
		common.ConvergeFailureResponse(ctx, common.ErrorContractCanNotUnfreeze)
		return
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		log.Error(err.Error())
	}

	_, err = SaveVote(params.ChainId, params.Reason, string(jsonBytes), UNFREEZE_CONTRACT, UPDATE_OTHER)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}

	contractInfo := &dbcommon.Contract{
		Name:            params.ContractName,
		MultiSignStatus: dbcommon.VOTING,
	}
	err = UpdateMultiSignStatus(contractInfo)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorUpdateVotingStatus)
		return
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)

}
