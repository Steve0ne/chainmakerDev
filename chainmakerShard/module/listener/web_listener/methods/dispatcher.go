/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package methods

import (
	"chainmaker.org/chainmaker-cross/conf"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	nonHandlerError = errors.New("can not find context handler to handle request")
	handlerMap      = make(map[string]ContextHandler)
	log             *zap.SugaredLogger
)

// InitHandlers init all handlers
func InitHandlers(logg *zap.SugaredLogger) {
	log = logg
	handlerMap[InvokeCrossEventMethod] = NewCrossEventContextHandler(log)
	handlerMap[GetCrossEventMethod] = NewCrossEventSearchContextHandler(log)
	if conf.Config.ListenerConfig.WebConfig.OpenTxRoute {
		handlerMap[TransactionEventMethod] = NewTransactionEventContextHandler(log)
	}
}

// Dispatch is dispatcher which will dispatch the request
func Dispatch(ctx *gin.Context) {
	contextHandler := ParseUrl(ctx)
	if contextHandler == nil {
		log.Error("can not find context handler to handle request")
		// 返回错误信息
		jsonResponse(ctx, http.StatusNotImplemented, nonHandlerError)
		return
	}
	contextHandler.Handle(ctx)
}

// jsonResponse wrapper the response
func jsonResponse(ctx *gin.Context, httpStatus int, data interface{}) {
	ctx.JSON(httpStatus, data)
}

// ParseUrl load method of this request
func ParseUrl(ctx *gin.Context) ContextHandler {
	log.Infof("Receive http request[%s]", ctx.Request.URL.String())
	param, ok := ctx.GetQuery(MethodType)
	if !ok {
		return nil
	}
	if handler, exist := handlerMap[param]; exist {
		return handler
	}
	return nil
}
