/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package overview

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/ctrl/contract_management"
	"management_backend/src/entity"
)

const (
	MAJORITY    = 0
	ROLE_ADMIN  = 0
	ROLE_CLIENT = 1
)

type ModifyChainAuthHandler struct {
}

func (handler *ModifyChainAuthHandler) LoginVerify() bool {
	return true
}

func (handler *ModifyChainAuthHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindModifyChainAuthHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	if params.Rule == MAJORITY {
		if len(params.OrgList) > 0 {
			common.ConvergeFailureResponse(ctx, common.ErrorMajorityPolicy)
			return
		}
	}

	for _, role := range params.RoleList {
		if params.Rule == MAJORITY && role.Role == ROLE_CLIENT {
			common.ConvergeFailureResponse(ctx, common.ErrorMajorityPolicy)
			return
		}
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorMarshalParameters)
		return
	}

	_, err = contract_management.SaveVote(params.ChainId, "", string(jsonBytes),
		contract_management.PERMISSION_UPDATE, contract_management.UPDATE_AUTH)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

type ModifyChainConfigHandler struct {
}

func (handler *ModifyChainConfigHandler) LoginVerify() bool {
	return true
}

func (handler *ModifyChainConfigHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindModifyChainConfigHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorMarshalParameters)
		return
	}

	_, err = contract_management.SaveVote(params.ChainId, params.Reason, string(jsonBytes),
		contract_management.BLOCK_UPDATE, contract_management.UPDATE_CONFIG)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
