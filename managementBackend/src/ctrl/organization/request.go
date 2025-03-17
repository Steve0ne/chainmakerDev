/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package organization

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

// 组织管理

type GetOrgListParams struct {
	ChainId  string
	OrgName  string
	PageNum  int
	PageSize int
}

func (params *GetOrgListParams) IsLegal() bool {
	if params.ChainId == "" {
		return false
	}
	if params.PageNum < 0 || params.PageSize == 0 {
		return false
	}
	return true
}

func BindGetOrgListHandler(ctx *gin.Context) *GetOrgListParams {
	var body = &GetOrgListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
