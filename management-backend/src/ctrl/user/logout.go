/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package user

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/entity"
	loggers "management_backend/src/logger"
	"management_backend/src/session"
)

type LogoutHandler struct{}

func (handler *LogoutHandler) LoginVerify() bool {
	return true
}

func (handler *LogoutHandler) Handle(user *entity.User, ctx *gin.Context) {
	// 将session清空
	err := session.UserStoreClear(ctx)
	if err != nil {
		log := loggers.GetLogger(loggers.ModuleSession)
		log.Debugf("clean captcha failed")
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
