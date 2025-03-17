/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_management

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

type InstallContractParams struct {
	ChainId         string
	ContractName    string
	ContractVersion string
	CompileSaveKey  string
	EvmAbiSaveKey   string
	RuntimeType     int
	Parameters      []*ParameterParams
	Methods         []*Method
	Reason          string
}

func (params *InstallContractParams) IsLegal() bool {
	if params.ChainId == "" || params.ContractName == "" || params.ContractVersion == "" ||
		params.CompileSaveKey == "" {
		return false
	}
	return true
}

type UpgradeContractParams struct {
	ChainId         string
	ContractName    string
	CompileSaveKey  string
	EvmAbiSaveKey   string
	ContractVersion string
	RuntimeType     int
	Parameters      []*ParameterParams
	Methods         []*Method
	Reason          string
}

func (params *UpgradeContractParams) IsLegal() bool {
	if params.ChainId == "" || params.ContractName == "" ||
		params.ContractVersion == "" {
		return false
	}
	return true
}

type Method struct {
	MethodName string
	MethodKey  string
}

type ParameterParams struct {
	Key   string
	Value string
}

func (params *ParameterParams) IsLegal() bool {
	return true
}

type FreezeContractParams struct {
	ChainId      string
	ContractName string
	Reason       string
}

func (params *FreezeContractParams) IsLegal() bool {
	if params.ChainId == "" || params.ContractName == "" {
		return false
	}
	return true
}

type GetContractManageListParams struct {
	ChainId      string
	ContractName string
	common.RangeBody
}

func (params *GetContractManageListParams) IsLegal() bool {
	return params.ChainId != ""
}

type ContractDetailParams struct {
	Id uint64
}

func (params *ContractDetailParams) IsLegal() bool {
	return params.Id > 0
}

type ModifyContractParams struct {
	Id            int64
	Methods       []*Method
	EvmAbiSaveKey string
}

func (params *ModifyContractParams) IsLegal() bool {
	return params.Id > 0
}

func BindInstallContractHandler(ctx *gin.Context) *InstallContractParams {
	var body = &InstallContractParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindFreezeContractHandler(ctx *gin.Context) *FreezeContractParams {
	var body = &FreezeContractParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindUnFreezeContractHandler(ctx *gin.Context) *FreezeContractParams {
	var body = &FreezeContractParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindRevokeContractHandler(ctx *gin.Context) *FreezeContractParams {
	var body = &FreezeContractParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindUpgradeContractHandler(ctx *gin.Context) *UpgradeContractParams {
	var body = &UpgradeContractParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetContractManageListHandler(ctx *gin.Context) *GetContractManageListParams {
	var body = &GetContractManageListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindContractDetailParamsHandler(ctx *gin.Context) *ContractDetailParams {
	var body = &ContractDetailParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindModifyContractParamsHandler(ctx *gin.Context) *ModifyContractParams {
	var body = &ModifyContractParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
