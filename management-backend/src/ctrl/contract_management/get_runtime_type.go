/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_management

import (
	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/entity"
)

var runtimeTypeList = arraylist.New("WASMER", "WXVM", "GASM", "EVM", "DOCKER_GO")

type GetRuntimeTypeListHandler struct{}

func (getRuntimeTypeListHandler *GetRuntimeTypeListHandler) LoginVerify() bool {
	return false
}

func (getRuntimeTypeListHandler *GetRuntimeTypeListHandler) Handle(user *entity.User, ctx *gin.Context) {
	runtimeTypeView := NewRuntimeTypeListView(pbcommon.RuntimeType_value)
	common.ConvergeListResponse(ctx, runtimeTypeView, int64(len(runtimeTypeView)), nil)
}

type RuntimeTypeListView struct {
	RuntimeTypeName string
	RuntimeTypeType int32
}

func NewRuntimeTypeListView(runtimeTypeMap map[string]int32) []interface{} {
	runtimeTypeViews := arraylist.New()
	for runtimeTypeName, runtimeTypeType := range runtimeTypeMap {
		runtimeTypeView := RuntimeTypeListView{
			RuntimeTypeName: runtimeTypeName,
			RuntimeTypeType: runtimeTypeType,
		}
		if runtimeTypeList.Contains(runtimeTypeView.RuntimeTypeName) {
			runtimeTypeViews.Add(runtimeTypeView)
		}
	}
	return runtimeTypeViews.Values()
}
