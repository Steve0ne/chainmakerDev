/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/relation"
	"management_backend/src/entity"
)

type GetCertNodeHandler struct{}

func (getCertNodeHandler *GetCertNodeHandler) LoginVerify() bool {
	return true
}

func (getCertNodeHandler *GetCertNodeHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetCertNodeListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	orgNodeList, err := relation.GetOrgNode(params.OrgId)
	if err != nil {
		orgsView := arraylist.New()
		common.ConvergeListResponse(ctx, orgsView.Values(), 0, nil)
		return
	}
	nodesView := NewCertNodeListView(orgNodeList)

	common.ConvergeListResponse(ctx, nodesView, int64(len(nodesView)), nil)
}
