/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package overview

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

// 区块链概览

type GeneralDataParams struct {
	ChainId string
}

func (params *GeneralDataParams) IsLegal() bool {
	return params.ChainId != ""
}

type GetChainDetailParams struct {
	Id      int64
	ChainId string
}

func (params *GetChainDetailParams) IsLegal() bool {
	if params.Id == 0 && params.ChainId == "" {
		return false
	}
	return true
}

type GetAuthOrgListParams struct {
	ChainId string
	Type    int
}

func (params *GetAuthOrgListParams) IsLegal() bool {
	return params.ChainId != ""
}

type GetAuthListParams struct {
	ChainId string
}

func (params *GetAuthListParams) IsLegal() bool {
	return params.ChainId != ""
}

type GetAuthRoleListParams struct {
	ChainId string
	Type    int
}

func (params *GetAuthRoleListParams) IsLegal() bool {
	return params.ChainId != ""
}

type ModifyChainConfigParams struct {
	ChainId         string
	BlockTxCapacity uint32
	TxTimeout       uint32
	BlockInterval   uint32
	Reason          string
}

func (params *ModifyChainConfigParams) IsLegal() bool {
	return params.ChainId != ""
}

type OrgListParams struct {
	OrgId string
}

type RoleListParams struct {
	Role int
}

type ModifyChainAuthParams struct {
	ChainId    string
	Type       int
	Rule       int
	PercentNum string
	OrgList    []*OrgListParams
	RoleList   []*RoleListParams
}

func (params *ModifyChainAuthParams) IsLegal() bool {
	return params.ChainId != ""
}

func BindGeneralDataHandler(ctx *gin.Context) *GeneralDataParams {
	var body = &GeneralDataParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetChainDetailHandler(ctx *gin.Context) *GetChainDetailParams {
	var body = &GetChainDetailParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetAuthOrgListHandler(ctx *gin.Context) *GetAuthOrgListParams {
	var body = &GetAuthOrgListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetAuthRoleListHandler(ctx *gin.Context) *GetAuthRoleListParams {
	var body = &GetAuthRoleListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetAuthListHandler(ctx *gin.Context) *GetAuthListParams {
	var body = &GetAuthListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindModifyChainConfigHandler(ctx *gin.Context) *ModifyChainConfigParams {
	var body = &ModifyChainConfigParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindModifyChainAuthHandler(ctx *gin.Context) *ModifyChainAuthParams {
	var body = &ModifyChainAuthParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
