/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package explorer

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/entity"
)

type GetTxListHandler struct {
}

func (handler *GetTxListHandler) LoginVerify() bool {
	return true
}

func (handler *GetTxListHandler) Handle(user *entity.User, ctx *gin.Context) {
	var (
		txList     []*dbcommon.Transaction
		totalCount int64
		offset     int
		limit      int
	)
	params := BindGetTxListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	offset = params.PageNum * params.PageSize
	limit = params.PageSize
	totalCount, txList, err := chain.GetTxList(params.ChainId, offset, limit, params.BlockHeight, params.ContractName)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	txInfos := convertToTxViews(txList)
	common.ConvergeListResponse(ctx, txInfos, totalCount, nil)
}

func convertToTxViews(txList []*dbcommon.Transaction) []interface{} {
	txViews := arraylist.New()
	for _, tx := range txList {
		txView := NewTransactionView(tx)
		txViews.Add(txView)
	}
	return txViews.Values()
}

type GetTxDetailHandler struct {
}

func (handler *GetTxDetailHandler) LoginVerify() bool {
	return true
}

func (handler *GetTxDetailHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetTxDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	var (
		tx  *dbcommon.Transaction
		err error
	)

	if params.Id != 0 {
		tx, err = chain.GetTxById(params.Id)
		if err != nil {
			common.ConvergeHandleFailureResponse(ctx, err)
			return
		}
	} else if params.TxId != "" {
		tx, err = chain.GetTxByTxId(params.ChainId, params.TxId)
		if err != nil {
			common.ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}
	txView := NewTransactionView(tx)
	common.ConvergeDataResponse(ctx, txView, nil)
}
