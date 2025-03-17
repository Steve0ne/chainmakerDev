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

type GetAuthRoleListHandler struct {
}

func (handler *GetAuthRoleListHandler) LoginVerify() bool {
	return true
}

func (handler *GetAuthRoleListHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetAuthRoleListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	var (
		err error
	)

	roleList, err := policy.GetRoleList(params.ChainId, params.Type)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}

	roleListInterface := make([]interface{}, len(roleList))
	for i, v := range roleList {
		roleListInterface[i] = v
	}
	common.ConvergeListResponse(ctx, roleListInterface, int64(len(roleList)), nil)
}
