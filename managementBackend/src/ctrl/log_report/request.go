/*
   Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
     SPDX-License-Identifier: Apache-2.0
*/
package log_report

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

// GetLogListParams 获取 日志 列表
type GetLogListParams struct {
	ChainId  string
	PageNum  int
	PageSize int
}

func (params *GetLogListParams) IsLegal() bool {
	if params.ChainId == "" || params.PageNum < 0 || params.PageSize <= 0 {
		return false
	}
	return true
}

// PullErrorLogParams 拉取错误日志
type PullErrorLogParams struct {
	ChainId string
}

func (params *PullErrorLogParams) IsLegal() bool {
	if params.ChainId == "" {
		return false
	}
	return true
}

type AutoReportLogFileParams struct {
	ChainId   string
	Automatic int
}

func (params *AutoReportLogFileParams) IsLegal() bool {
	if params.ChainId == "" {
		return false
	}
	return true
}

type DownloadLogFileParams struct {
	Id int64
}

func (params *DownloadLogFileParams) IsLegal() bool {
	if params.Id == 0 {
		return false
	}
	return true
}

type ReportLogFileParams struct {
	Id int64
}

func (params *ReportLogFileParams) IsLegal() bool {
	if params.Id == 0 {
		return false
	}
	return true
}

// 链错误日志获取

func BindGetLogListHandler(ctx *gin.Context) *GetLogListParams {
	var body = &GetLogListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindDownloadLogFileHandler(ctx *gin.Context) *DownloadLogFileParams {
	var body = &DownloadLogFileParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindPullErrorLogHandler(ctx *gin.Context) *PullErrorLogParams {
	var body = &PullErrorLogParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindAutoReportLogFileHandler(ctx *gin.Context) *AutoReportLogFileParams {
	var body = &AutoReportLogFileParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindReportLogFileHandler(ctx *gin.Context) *ReportLogFileParams {
	var body = &ReportLogFileParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
