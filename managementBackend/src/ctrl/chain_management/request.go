/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

type Nodes struct {
	OrgId    string
	NodeList string
}

type AddChainParams struct {
	ChainId           string
	ChainName         string
	BlockTxCapacity   uint32
	BlockInterval     uint32
	TxTimeout         uint32
	Consensus         int32
	Nodes             []Nodes
	Monitor           int
	ChainmakerImprove int
	Address           string
	Tls               int //0 开启 1 关闭   （默认开启）
	DockerVm          int //0 关闭 1 开始   （默认关闭）
}

func (params *AddChainParams) IsLegal() bool {
	if params.BlockTxCapacity < 0 || params.BlockInterval < 0 ||
		params.TxTimeout < 0 || params.ChainId == "" || params.ChainName == "" {
		return false
	}
	return true
}

type GetCertUserListParams struct {
	OrgId string
}

func (params *GetCertUserListParams) IsLegal() bool {
	return true
}

type GetCertOrgListParams struct {
	ChainId string
}

func (params *GetCertOrgListParams) IsLegal() bool {
	return true
}

type GetCertNodeListParams struct {
	ChainId string
	OrgId   string
}

func (params *GetCertNodeListParams) IsLegal() bool {
	return params.OrgId != ""
}

type SubscribeChainParams struct {
	ChainId  string
	OrgId    string
	NodeName string
	UserName string
	Tls      int //0 开启 1 关闭
}

func (params *SubscribeChainParams) IsLegal() bool {
	if params.OrgId == "" || params.ChainId == "" ||
		params.NodeName == "" || params.UserName == "" {
		return false
	}
	return true
}

type DeleteChainParams struct {
	ChainId string
}

func (params *DeleteChainParams) IsLegal() bool {
	return params.ChainId != ""
}

type DownloadChainConfigParams struct {
	ChainId string
}

func (params *DownloadChainConfigParams) IsLegal() bool {
	return params.ChainId != ""
}

func BindAddChainHandler(ctx *gin.Context) *AddChainParams {
	var body = &AddChainParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindDeleteChainHandler(ctx *gin.Context) *DeleteChainParams {
	var body = &DeleteChainParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetCertUserListHandler(ctx *gin.Context) *GetCertUserListParams {
	var body = &GetCertUserListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetCertOrgListHandler(ctx *gin.Context) *GetCertOrgListParams {
	var body = &GetCertOrgListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetCertNodeListHandler(ctx *gin.Context) *GetCertNodeListParams {
	var body = &GetCertNodeListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindSubscribeChainHandler(ctx *gin.Context) *SubscribeChainParams {
	var body = &SubscribeChainParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindDownloadChainConfigHandler(ctx *gin.Context) *DownloadChainConfigParams {
	var body = &DownloadChainConfigParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
