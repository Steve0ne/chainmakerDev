/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package methods

import (
	"net/http"

	"chainmaker.org/chainmaker-cross/net/net_http"
	"github.com/gin-gonic/gin"
)

type Response = net_http.Response
type Request = net_http.Request

func jsonOkResponse(ctx *gin.Context, data Response) {
	ctx.JSON(http.StatusOK, data)
}
