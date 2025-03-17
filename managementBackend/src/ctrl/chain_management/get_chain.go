/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain"
	"management_backend/src/entity"
)

type GetChainListHandler struct{}

func (getChainListHandler *GetChainListHandler) LoginVerify() bool {
	return true
}

func (getChainListHandler *GetChainListHandler) Handle(user *entity.User, ctx *gin.Context) {

	chainList, err := chain.GetChainList()
	if err != nil {
		orgsView := arraylist.New()
		common.ConvergeListResponse(ctx, orgsView.Values(), 0, nil)
		return
	}
	chainListView := NewChainListView(chainList)
	common.ConvergeListResponse(ctx, chainListView, int64(len(chainListView)), nil)
}
