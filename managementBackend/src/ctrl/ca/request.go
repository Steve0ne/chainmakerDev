/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

type GenerateCertParams struct {
	OrgId    string
	OrgName  string
	NodeName string
	NodeRole int
	CertType int
	NodeIp   string
	UserName string
	UserRole int
}

func (params *GenerateCertParams) IsLegal() bool {
	if params.OrgName == "" || params.OrgId == "" || params.CertType > 2 {
		return false
	}
	return true
}

type GetCertParams struct {
	CertId  int64
	CertUse int
}

func (params *GetCertParams) IsLegal() bool {
	return params.CertId >= 0
}

type GetCertListParams struct {
	Type     int
	OrgName  string
	NodeName string
	UserName string
	common.RangeBody
}

func (params *GetCertListParams) IsLegal() bool {
	if params.Type < 0 || params.Type > 3 {
		return false
	}
	return true
}

type ImportCertParams struct {
	Type     int
	Role     int
	OrgId    string
	OrgName  string
	NodeName string
	UserName string
	CaCert   string
	CaKey    string
	SignCert string
	SignKey  string
	TlsCert  string
	TlsKey   string
	NodeIp   string
}

func (params *ImportCertParams) IsLegal() bool {
	if params.Type < 0 || params.Type > 4 || params.OrgId == "" || params.OrgName == "" {
		return false
	}
	return true
}

type DownloadCertParams struct {
	CertId  int64
	CertUse int
}

func (params *DownloadCertParams) IsLegal() bool {
	return params.CertId >= 0
}

func BindGenerateCertHandler(ctx *gin.Context) *GenerateCertParams {
	var body = &GenerateCertParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetCertHandler(ctx *gin.Context) *GetCertParams {
	var body = &GetCertParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetCertListHandler(ctx *gin.Context) *GetCertListParams {
	var body = &GetCertListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindImportCertHandler(ctx *gin.Context) *ImportCertParams {
	var body = &ImportCertParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindDownloadCertHandler(ctx *gin.Context) *DownloadCertParams {
	var body = &DownloadCertParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
