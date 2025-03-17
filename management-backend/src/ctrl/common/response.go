/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package common

import (
	"fmt"

	"management_backend/src/utils"
)

// SuccessDataResponse 成功的单一数据应答
type SuccessDataResponse struct {
	Response DataResponse
}

// SuccessListResponse 成功的列表数据应答
type SuccessListResponse struct {
	Response ListResponse
}

type SuccessListStatusResponse struct {
	Response ListStatusResponse
}

// FailureResponse 失败的应答
type FailureResponse struct {
	Response ErrorResponse
}

// DataResponse 单一对象
type DataResponse struct {
	Data interface{}
	//RequestId string
}

// ListResponse 集合对象
type ListResponse struct {
	GroupList  []interface{}
	TotalCount int64
	//RequestId  string
}

type ListStatusResponse struct {
	GroupList  []interface{}
	TotalCount int64
	Status     int
	RequestId  string
}

// ErrorResponse 异常应答
type ErrorResponse struct {
	Error Error
	//RequestId string
}

// Error 错误
type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s - %s", e.Code, e.Message)
}

type StatusIntegerResponse struct {
	Status int
}

type StatusResponse struct {
	Status string
}

func NewStatusResponse() *StatusResponse {
	return &StatusResponse{
		Status: "OK",
	}
}

type TokenResponse struct {
	Token string
}

type DownloadResponse struct {
	Content string
}

type UploadResponse struct {
	FileKey string
}

func NewDownloadResponse(content []byte) *DownloadResponse {
	base64Encode := utils.Base64Encode(content)
	return &DownloadResponse{
		Content: base64Encode,
	}
}

func NewUploadResponse(key string) *UploadResponse {
	return &UploadResponse{
		FileKey: key,
	}
}

func NewSuccessDataResponse(data interface{}) *SuccessDataResponse {
	dataResponse := DataResponse{
		//RequestId: NewRandomRequestId(),
		Data: data,
	}
	return &SuccessDataResponse{
		Response: dataResponse,
	}
}

func NewSuccessListResponse(datas []interface{}, count int64) *SuccessListResponse {
	listResp := ListResponse{
		GroupList:  datas,
		TotalCount: count,
		//RequestId:  NewRandomRequestId(),
	}
	return &SuccessListResponse{
		Response: listResp,
	}
}

func NewSuccessListStatusResponse(datas []interface{}, status int, count int64) *SuccessListStatusResponse {
	listResp := ListStatusResponse{
		GroupList:  datas,
		TotalCount: count,
		Status:     status,
		//RequestId:  NewRandomRequestId(),
	}
	return &SuccessListStatusResponse{
		Response: listResp,
	}
}

func NewFailureResponse(err *Error) *FailureResponse {
	errResponse := ErrorResponse{
		Error: *err,
	}
	return &FailureResponse{
		Response: errResponse,
	}
}

// NewError 创建错误
func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
