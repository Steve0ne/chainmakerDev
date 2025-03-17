/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain_participant"
	"management_backend/src/entity"
)

type GetCertHandler struct{}

func (getCertHandler *GetCertHandler) LoginVerify() bool {
	return true
}

func (getCertHandler *GetCertHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetCertHandler(ctx)
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
	if params.CertUse == KEY_FOR_SIGN || params.CertUse == KEY_FOR_TLS {
		detail = certInfo.PrivateKey
	} else {
		detail = certInfo.Cert
	}
	certView := NewCertDetailView(detail)
	common.ConvergeDataResponse(ctx, certView, nil)
}
