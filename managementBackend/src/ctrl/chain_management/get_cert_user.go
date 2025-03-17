/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain_participant"
	"management_backend/src/entity"
)

type GetCertUserListHandler struct{}

func (getCertUserListHandler *GetCertUserListHandler) LoginVerify() bool {
	return true
}

func (getCertUserListHandler *GetCertUserListHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetCertUserListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	userCertList, count, err := chain_participant.GetSignUserCertList(params.OrgId)
	if err != nil {
		userCertsView := arraylist.New()
		common.ConvergeListResponse(ctx, userCertsView.Values(), 0, nil)
		return
	}

	userCertsView := NewCertUserListView(userCertList)
	common.ConvergeListResponse(ctx, userCertsView, count, nil)
}
