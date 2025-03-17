/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package node

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

type GetNodeListParams struct {
	ChainId  string
	NodeName string
	PageNum  int
	PageSize int
}

func (params *GetNodeListParams) IsLegal() bool {
	if params.ChainId == "" {
		return false
	}
	if params.PageNum < 0 || params.PageSize == 0 {
		return false
	}
	return true
}

type GetNodeDetailParams struct {
	ChainId string
	NodeId  int
}

func (params *GetNodeDetailParams) IsLegal() bool {
	if params.ChainId == "" || params.NodeId == 0 {
		return false
	}
	return true
}

func BindGetNodeListHandler(ctx *gin.Context) *GetNodeListParams {
	var body = &GetNodeListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetNodeDetailHandler(ctx *gin.Context) *GetNodeDetailParams {
	var body = &GetNodeDetailParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
