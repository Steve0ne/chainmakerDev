/*
   Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
     SPDX-License-Identifier: Apache-2.0
*/
package log_report

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/entity"
)

type GetLogListHandler struct{}

func (handler *GetLogListHandler) LoginVerify() bool {
	return true
}

func (handler *GetLogListHandler) Handle(user *entity.User, ctx *gin.Context) {
	var (
		logList    []*dbcommon.ChainErrorLog
		totalCount int64
		offset     int
		limit      int
	)

	params := BindGetLogListHandler(ctx)

	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	offset = params.PageNum * params.PageSize
	limit = params.PageSize
	totalCount, logList, err := chain.GetLogList(params.ChainId, offset, limit)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	logInfos := convertToLogInfoListViews(logList)
	common.ConvergeListResponse(ctx, logInfos, totalCount, nil)
}
