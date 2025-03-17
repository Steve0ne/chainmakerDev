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

type RevokeContractHandler struct{}

func (revokeContractHandler *RevokeContractHandler) LoginVerify() bool {
	return true
}

func (revokeContractHandler *RevokeContractHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindRevokeContractHandler(ctx)
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

	if !contract.CanRevoke() {
		// 不可以进行注销操作
		common.ConvergeFailureResponse(ctx, common.ErrorContractCanNotRevoke)
		return
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		log.Error(err.Error())
	}

	_, err = SaveVote(params.ChainId, params.Reason, string(jsonBytes), REVOKE_CONTRACT, UPDATE_OTHER)
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
