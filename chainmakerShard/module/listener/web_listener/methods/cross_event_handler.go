/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package methods

import (
	"net/http"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/handler"
	"chainmaker.org/chainmaker-cross/net/net_http"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _ ContextHandler = (*CrossEventContextHandler)(nil)

// CrossEventContextHandler is handler which handle cross event
type CrossEventContextHandler struct {
	eventHandler handler.EventHandler
}

// NewCrossEventContextHandler create new cross event context handler
func NewCrossEventContextHandler(log *zap.SugaredLogger) *CrossEventContextHandler {
	if eveHandler, exist := handler.GetEventHandlerTools().GetHandler(handler.CrossProcess); exist {
		return &CrossEventContextHandler{
			eventHandler: eveHandler,
		}
	} else {
		log.Error("can not find handler for cross event")
		return nil
	}
}

// Handle receive cross event and start handle
func (c *CrossEventContextHandler) Handle(ctx *gin.Context) {
	// 获取cross-event
	crossEvent := &eventproto.CrossEvent{}
	if err := ctx.ShouldBindJSON(crossEvent); err != nil {
		log.Error("resolve param error:", err)
		return
	}
	go func(crossEvent *eventproto.CrossEvent) {
		// 独立goroutine处理该状态
		if _, err := c.eventHandler.Handle(crossEvent, false); err != nil {
			// 打印日志信息
			log.Errorf("handle cross event[%s] error", crossEvent.GetCrossID(), err)
		}
	}(crossEvent)
	// 返回crossID
	crossID := crossEvent.GetCrossID()
	jsonResponse(ctx, http.StatusOK, NewDefaultCrossEventResp(crossID))
}

// DefaultCrossEventResp is default cross event response
type DefaultCrossEventResp struct {
	CrossID string
}

// NewDefaultCrossEventResp create new default cross event response
func NewDefaultCrossEventResp(crossID string) *DefaultCrossEventResp {
	return &DefaultCrossEventResp{
		CrossID: crossID,
	}
}

// CrossEventSearchContextHandler is handler which will handle cross event search context
type CrossEventSearchContextHandler struct {
	CrossEventContextHandler
}

// NewCrossEventSearchContextHandler create new instance of CrossEventSearchContextHandler
func NewCrossEventSearchContextHandler(log *zap.SugaredLogger) *CrossEventSearchContextHandler {
	if eveHandler, exist := handler.GetEventHandlerTools().GetHandler(handler.CrossSearch); exist {
		hd := &CrossEventSearchContextHandler{}
		hd.eventHandler = eveHandler
		return hd
	} else {
		log.Error("can not find handler for cross search event")
		return nil
	}
}

// Handle handle the request of cross search event
func (l *CrossEventSearchContextHandler) Handle(ctx *gin.Context) {
	// 获取cross-event
	crossSearchEvent := &eventproto.CrossSearchEvent{}
	if err := ctx.ShouldBindJSON(crossSearchEvent); err != nil {
		log.Error("resolve param error:", err)
		return
	}
	crossResult, err := l.eventHandler.Handle(crossSearchEvent, true)
	if err != nil {
		// 打印日志信息
		log.Errorf("handle cross event[%s] error", crossSearchEvent.GetCrossID(), err)
	}
	jsonResponse(ctx, http.StatusOK, crossResult)
}

type TransactionEventContextHandler struct {
	eventHandler handler.EventHandler
}

func NewTransactionEventContextHandler(log *zap.SugaredLogger) *TransactionEventContextHandler {
	if eveHandler, exist := handler.GetEventHandlerTools().GetHandler(handler.TransactionProcess); exist {
		return &TransactionEventContextHandler{
			eventHandler: eveHandler,
		}
	} else {
		log.Error("can not find handler for transaction event")
		return nil
	}
}

func (t *TransactionEventContextHandler) Handle(ctx *gin.Context) {
	req := &Request{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		log.Error("resolve param error:", err)
		jsonOkResponse(ctx, Response{
			Code:    100,
			Message: "resolve param error",
		})
		return
	}
	tec := &event.TransactionEventContext{}
	if err := req.Unmarshal(tec); err != nil {
		log.Error("request unmarshal error:", err)
		jsonOkResponse(ctx, Response{
			Code:    100,
			Message: err.Error(),
		})
		return
	}
	result, err := t.eventHandler.Handle(tec, true)
	if err != nil {
		// 打印错误信息
		log.Error("handle transaction event failed: ", err)
		jsonOkResponse(ctx, Response{
			Code:    100,
			Message: err.Error(),
		})
	} else {
		// 需要结果是*event.ProofResponse
		if resp, ok := result.(*event.ProofResponse); ok {
			data, err := net_http.NewMessage(resp, event.BinaryMarshalType)
			if err != nil {
				log.Error("ProofResponse convert to NewMessage fail", err)
				return
			}
			log.Infof("resp result is [%+v]", *resp)
			jsonOkResponse(ctx, Response{
				Data: data,
			})
		} else {
			// 打印信息
			log.Error("http resp result is not the type of ProofResponse")
			jsonOkResponse(ctx, Response{
				Code:    100,
				Message: "resp result can not convert to ProofResponse",
			})
		}
	}
}
