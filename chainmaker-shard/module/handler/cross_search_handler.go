/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"errors"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/event/coder"
	"chainmaker.org/chainmaker-cross/store"
	storetype "chainmaker.org/chainmaker-cross/store/types"
	"go.uber.org/zap"
)

var crossSearchHandler *CrossSearchHandler

func init() {
	crossSearchHandler = &CrossSearchHandler{}
}

// CrossSearchHandler the struct of handler which handle request of research cross event state
type CrossSearchHandler struct {
	stateDB store.StateDB      // 存储
	coder   event.EventCoder   // 编解码器
	logger  *zap.SugaredLogger // log
}

// GetCrossSearchHandler return instance of CrossSearchHandler
func GetCrossSearchHandler() *CrossSearchHandler {
	return crossSearchHandler
}

// Init init base param
func (c *CrossSearchHandler) Init() {
	if eventCoder, exist := coder.GetEventCoderTools().GetDefaultCoder(eventproto.CrossRespEventType); exist {
		c.coder = eventCoder
	} else {
		panic("can not find coder for cross response event")
	}
}

// SetStateDB set state database
func (c *CrossSearchHandler) SetStateDB(stateDB store.StateDB) {
	c.stateDB = stateDB
}

// SetLogger set logger
func (c *CrossSearchHandler) SetLogger(logger *zap.SugaredLogger) {
	c.logger = logger
}

// GetType return type of this handler
func (c *CrossSearchHandler) GetType() HandlerType {
	return CrossSearch
}

// Handle handle event which just search state from db
func (c *CrossSearchHandler) Handle(eve event.Event, _ bool) (interface{}, error) {
	// 首先查询本地
	if crossSearchEvent, ok := eve.(*eventproto.CrossSearchEvent); ok {
		return c.LoadCrossEventResp(crossSearchEvent.GetCrossID()), nil
	} else {
		c.logger.Error("event is not type of CrossSearchEvent")
		return nil, errors.New("event is not type of CrossSearchEvent")
	}
}

// LoadCrossEventResp load cross event response
func (c *CrossSearchHandler) LoadCrossEventResp(crossID string) event.Event {
	crossState, valBytes, exist := c.stateDB.ReadCrossState(crossID)
	if exist {
		if crossState == storetype.StateSuccess {
			eve, err := c.coder.UnmarshalFromBinary(valBytes)
			if err != nil {
				return event.NewCrossResponse(crossID, event.FailureResp, err.Error())
			}
			return eve
		}
		if crossState == storetype.StateFailed {
			return event.NewCrossResponse(crossID, event.FailureResp, string(valBytes))
		}
		return event.NewCrossResponse(crossID, event.UnknownResp, "you should research again")
	} else {
		return event.NewCrossResponse(crossID, event.ErrorResp, "can not find cross state from db")
	}
}
