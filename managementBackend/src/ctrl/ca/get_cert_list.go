/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain_participant"
	"management_backend/src/entity"
)

type GetCertListHandler struct{}

func (getCertListHandler *GetCertListHandler) LoginVerify() bool {
	return false
}

func (getCertListHandler *GetCertListHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetCertListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}
	certs, count, err := chain_participant.GetCertList(params.PageNum, params.PageSize,
		params.Type, params.OrgName, params.NodeName, params.UserName)
	if err != nil {
		certsView := arraylist.New()
		common.ConvergeListResponse(ctx, certsView.Values(), 0, nil)
		return
	}
	certsView := NewCertListView(certs)
	common.ConvergeListResponse(ctx, certsView, count*2, nil)
}
