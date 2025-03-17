/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package user

import (
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
)

// 用户管理
// RegisterParams 注册
type RegisterParams struct {
	UserName string
	Name     string
	Password string
}

func (registerParams *RegisterParams) IsLegal() bool {
	if registerParams.UserName == "" || registerParams.Password == "" {
		return false
	}
	return true
}

// LoginParams 登录
type LoginParams struct {
	UserName string
	Password string
	Captcha  string
}

func (loginParams *LoginParams) IsLegal() bool {
	if loginParams.UserName == "" || loginParams.Password == "" || loginParams.Captcha == "" {
		return false
	}
	return true
}

// GetUserListParams 查询账户列表
type GetUserListParams struct {
	PageNum  int
	PageSize int
}

func (getUserListParams *GetUserListParams) IsLegal() bool {
	if getUserListParams.PageNum < 0 || getUserListParams.PageSize == 0 {
		return false
	}
	return true
}

// ModifyPasswordParams 修改用户密码
type ModifyPasswordParams struct {
	Password    string
	OldPassword string
}

func (modifyPasswordParams *ModifyPasswordParams) IsLegal() bool {
	return len(modifyPasswordParams.Password) > 0 && len(modifyPasswordParams.OldPassword) > 0
}

// UserIdParams 参数为账户id
// 用于启用，禁用账户，账户登出等。
type UserIdParams struct {
	UserId int64
}

func (userIdParams *UserIdParams) IsLegal() bool {
	return userIdParams.UserId > 0
}

// User Login & management

func BindLoginHandler(ctx *gin.Context) *LoginParams {
	var body = &LoginParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindRegisterHandler(ctx *gin.Context) *RegisterParams {
	var body = &RegisterParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindGetUserListHandler(ctx *gin.Context) *GetUserListParams {
	var body = &GetUserListParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindPasswordHandler(ctx *gin.Context) *ModifyPasswordParams {
	var body = &ModifyPasswordParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindUserIdParams(ctx *gin.Context) *UserIdParams {
	var body = &UserIdParams{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
