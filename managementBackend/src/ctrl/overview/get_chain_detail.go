/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package overview

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	dbchain "management_backend/src/db/chain"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/entity"
)

type GetChainDetailHandler struct {
}

func (handler *GetChainDetailHandler) LoginVerify() bool {
	return true
}

func (handler *GetChainDetailHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetChainDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	var (
		chain *dbcommon.Chain
		err   error
	)

	if params.Id != 0 {
		chain, err = dbchain.GetChainById(params.Id)
		if err != nil {
			common.ConvergeHandleFailureResponse(ctx, err)
			return
		}
	} else {
		chain, err = dbchain.GetChainByChainId(params.ChainId)
		if err != nil {
			common.ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}

	chainView := NewChainView(chain)
	common.ConvergeDataResponse(ctx, chainView, nil)
}
