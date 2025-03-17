/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"fmt"
	"management_backend/src/utils"
	"net/http"

	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain_participant"
	"management_backend/src/entity"
)

type DownLoadCertHandler struct{}

func (downLoadCertHandler *DownLoadCertHandler) LoginVerify() bool {
	return true
}

func (downLoadCertHandler *DownLoadCertHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindDownloadCertHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	certInfo, err := chain_participant.GetCertById(params.CertId)
	if err != nil {
		log.Error("ErrorGetCert err : " + err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorGetCert)
		return
	}

	var detail string
	var suffix string
	if params.CertUse == KEY_FOR_SIGN || params.CertUse == KEY_FOR_TLS {
		detail = certInfo.PrivateKey
		suffix = ".key"
	} else {
		detail = certInfo.Cert
		suffix = ".crt"
	}

	var certName string
	content := []byte(detail)
	if certInfo.CertType == chain_participant.ORG_CA {
		certName = certInfo.OrgName + suffix
	}
	if certInfo.CertType == chain_participant.ADMIN || certInfo.CertType == chain_participant.CLIENT {
		certName = certInfo.CertUserName + suffix
	}
	if certInfo.CertType == chain_participant.CONSENSUS || certInfo.CertType == chain_participant.COMMON {
		certName = certInfo.NodeName + suffix
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", "attachment; filename="+utils.Base64Encode([]byte(certName)))
	ctx.Header("Content-Type", "application/text/plain")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
	ctx.Header("Access-Control-Expose-Headers", "Content-Disposition")
	_, err = ctx.Writer.Write(content)
	if err != nil {
		log.Error("ctx write content err : " + err.Error())
	}
}
