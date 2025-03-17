/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_invoke

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/ctrl/explorer"
	"management_backend/src/db/chain"
	"management_backend/src/db/contract"
	"management_backend/src/entity"
)

type GetInvokeRecordListHandler struct{}

func (getInvokeRecordListHandler *GetInvokeRecordListHandler) LoginVerify() bool {
	return true
}

func (getInvokeRecordListHandler *GetInvokeRecordListHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetInvokeRecordListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	invokeRecordsList, count, err :=
		contract.GetInvokeRecordsByChainId(params.PageNum, params.PageSize, params.ChainId, params.TxId, params.Status, params.TxStatus)
	if err != nil {
		certsView := arraylist.New()
		common.ConvergeListResponse(ctx, certsView.Values(), 0, nil)
		return
	}
	certsView := NewInvokeRecordListView(invokeRecordsList)
	common.ConvergeListResponse(ctx, certsView, count, nil)
}

type GetInvokeRecordDetailHandler struct{}

func (getInvokeRecordDetailHandler *GetInvokeRecordDetailHandler) LoginVerify() bool {
	return true
}

func (getInvokeRecordDetailHandler *GetInvokeRecordDetailHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetInvokeRecordDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	tx, err := chain.GetTxByTxId(params.ChainId, params.TxId)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorTXNotOnChain)
		return
	}
	txView := explorer.NewTransactionView(tx)
	common.ConvergeDataResponse(ctx, txView, nil)
}

type GetInvokeContractListHandler struct{}

func (getInvokeContractListHandler *GetInvokeContractListHandler) LoginVerify() bool {
	return true
}

func (getInvokeContractListHandler *GetInvokeContractListHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetInvokeContractListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}
	contractList, err := contract.GetContractList(params.ChainId)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	certsView := NewInvokeContractListView(contractList)
	common.ConvergeListResponse(ctx, certsView, int64(len(contractList)), nil)
}
