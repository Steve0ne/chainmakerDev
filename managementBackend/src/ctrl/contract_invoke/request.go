/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_invoke

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

type ParameterParams struct {
	Key   string
	Value string
}

func (params *ParameterParams) IsLegal() bool {
	return true
}

type GetInvokeContractListParams struct {
	ChainId string
}

func (params *GetInvokeContractListParams) IsLegal() bool {
	return params.ChainId != ""
}

type InvokeContractListParams struct {
	ChainId      string
	ContractName string
	MethodName   string
	Parameters   []*ParameterParams
}

func (params *InvokeContractListParams) IsLegal() bool {
	if params.ChainId == "" || params.ContractName == "" || params.MethodName == "" {
		return false
	}
	return true
}

type ReInvokeContractParams struct {
	InvokeRecordId int64
}

func (params *ReInvokeContractParams) IsLegal() bool {
	return params.InvokeRecordId >= 1
}

type GetInvokeRecordListParams struct {
	ChainId  string
	TxId     string
	Status   int
	TxStatus int
	common.RangeBody
}

func (params *GetInvokeRecordListParams) IsLegal() bool {
	return params.ChainId != ""
}

type GetInvokeRecordDetailParams struct {
	ChainId string
	TxId    string
}

func (params *GetInvokeRecordDetailParams) IsLegal() bool {
	if params.ChainId == "" || params.TxId == "" {
		return false
	}
	return true
}

func BindGetInvokeContractListHandler(ctx *gin.Context) *GetInvokeContractListParams {
	var body = &GetInvokeContractListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindInvokeContractHandler(ctx *gin.Context) *InvokeContractListParams {
	var body = &InvokeContractListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindReInvokeContractHandler(ctx *gin.Context) *ReInvokeContractParams {
	var body = &ReInvokeContractParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetInvokeRecordListHandler(ctx *gin.Context) *GetInvokeRecordListParams {
	var body = &GetInvokeRecordListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetInvokeRecordDetailHandler(ctx *gin.Context) *GetInvokeRecordDetailParams {
	var body = &GetInvokeRecordDetailParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
