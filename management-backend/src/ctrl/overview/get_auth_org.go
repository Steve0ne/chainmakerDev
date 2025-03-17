/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package overview

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/policy"
	"management_backend/src/entity"
)

type GetAuthOrgListHandler struct {
}

func (handler *GetAuthOrgListHandler) LoginVerify() bool {
	return true
}

func (handler *GetAuthOrgListHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetAuthOrgListHandler(ctx)

	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	var (
		orgList []*dbcommon.ChainPolicyOrg
	)

	orgList, err := policy.GetOrgListByPolicyType(params.ChainId, params.Type)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	orgInfos := convertToPolicyOrgView(orgList)
	common.ConvergeListResponse(ctx, orgInfos, int64(len(orgInfos)), nil)
}

func convertToPolicyOrgView(orgList []*dbcommon.ChainPolicyOrg) []interface{} {
	views := arraylist.New()
	for _, org := range orgList {
		view := NewPolicyOrgView(org)
		views.Add(view)
	}
	return views.Values()
}
