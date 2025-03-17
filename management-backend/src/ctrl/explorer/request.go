/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package explorer

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

// 浏览器 接口的请求
type GetTxListParams struct {
	ChainId      string
	TxId         string
	BlockHeight  *int64
	ContractName string
	PageNum      int
	PageSize     int
}

func (params *GetTxListParams) IsLegal() bool {
	if params.ChainId == "" || params.PageNum < 0 || params.PageSize <= 0 {
		return false
	}
	return true
}

type GetTxDetailParams struct {
	ChainId string
	Id      uint64
	TxId    string
}

func (params *GetTxDetailParams) IsLegal() bool {
	if params.ChainId == "" {
		return false
	}
	if params.Id == 0 && params.TxId == "" {
		return false
	}
	return true
}

// 区块查询

type GetBlockListParams struct {
	ChainId     string
	BlockHash   string
	BlockHeight string
	PageNum     int
	PageSize    int
}

func (params *GetBlockListParams) IsLegal() bool {
	if params.ChainId == "" || params.PageNum < 0 || params.PageSize <= 0 {
		return false
	}
	return true
}

type GetBlockDetailParams struct {
	ChainId     string
	Id          uint64
	BlockHeight uint64
	BlockHash   string
}

func (params *GetBlockDetailParams) IsLegal() bool {
	if params.ChainId == "" {
		return false
	}
	if params.Id == 0 && params.BlockHeight == 0 {
		return false
	}
	return true
}

// 合约查询

type GetContractListParams struct {
	ChainId      string
	ContractName string
	PageNum      int
	PageSize     int
}

func (params *GetContractListParams) IsLegal() bool {
	return params.ChainId != ""
}

type GetContractDetailParams struct {
	ChainId      string
	Id           uint64
	ContractName string
}

func (params *GetContractDetailParams) IsLegal() bool {
	if params.Id > 0 {
		return true
	}

	if params.ChainId == "" && params.ContractName == "" {
		return false
	}
	return true
}

type HomePageSearchParams struct {
	KeyWord string
	ChainId string
}

func (params *HomePageSearchParams) IsLegal() bool {
	if params.KeyWord == "" || params.ChainId == "" {
		return false
	}
	return true
}

// Explorer Handler

func BindGetTxListHandler(ctx *gin.Context) *GetTxListParams {
	var body = &GetTxListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetTxDetailHandler(ctx *gin.Context) *GetTxDetailParams {
	var body = &GetTxDetailParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetBlockListHandler(ctx *gin.Context) *GetBlockListParams {
	var body = &GetBlockListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetBlockDetailHandler(ctx *gin.Context) *GetBlockDetailParams {
	var body = &GetBlockDetailParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetContractListHandler(ctx *gin.Context) *GetContractListParams {
	var body = &GetContractListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetContractDetailHandler(ctx *gin.Context) *GetContractDetailParams {
	var body = &GetContractDetailParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindHomePageSearchHandler(ctx *gin.Context) *HomePageSearchParams {
	var body = &HomePageSearchParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
