/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package overview

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/policy"
	"management_backend/src/entity"
)

type GetAuthListHandler struct {
}

func (handler *GetAuthListHandler) LoginVerify() bool {
	return true
}

func (handler *GetAuthListHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetAuthListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}
	chainPolicy, err := policy.GetChainPolicyByChainId(params.ChainId)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	authListView := NewAuthListView(chainPolicy)
	common.ConvergeListResponse(ctx, authListView, int64(len(authListView)), nil)
}
