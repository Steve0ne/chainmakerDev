/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package vote

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

type VoteParams struct {
	VoteId     int64
	VoteResult int
}

func (params *VoteParams) IsLegal() bool {
	return params.VoteId != 0
}

type GetVoteManageListParams struct {
	ChainId    string
	OrgId      string
	VoteType   *int
	VoteStatus *int
	PageNum    int
	PageSize   int
}

func (params *GetVoteManageListParams) IsLegal() bool {
	if params.ChainId == "" || params.OrgId == "" {
		return false
	}
	if params.PageNum < 0 || params.PageSize == 0 {
		return false
	}
	return true
}

type GetVoteDetailParams struct {
	VoteId int64
}

func (params *GetVoteDetailParams) IsLegal() bool {
	return params.VoteId != 0
}

// 投票管理
func BindVoteHandler(ctx *gin.Context) *VoteParams {
	var body = &VoteParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetVoteManageListHandler(ctx *gin.Context) *GetVoteManageListParams {
	var body = &GetVoteManageListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetVoteDetailHandler(ctx *gin.Context) *GetVoteDetailParams {
	var body = &GetVoteDetailParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
