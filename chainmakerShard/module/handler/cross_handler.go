/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"errors"
	"fmt"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/event/coder"
	"chainmaker.org/chainmaker-cross/store"
	"go.uber.org/zap"
)

var crossProcessHandler *CrossProcessHandler

func init() {
	crossEventCoder, exist := coder.GetEventCoderTools().GetDefaultCoder(eventproto.CrossEventType)
	if !exist {
		panic("can not find coder for cross event")
	}
	crossProcessHandler = &CrossProcessHandler{
		coder: crossEventCoder,
	}
}

// CrossProcessHandler the struct of handler which handle cross process
type CrossProcessHandler struct {
	eventChan chan event.Event   // CrossEvent 消息
	stateDB   store.StateDB      // 存储
	coder     event.EventCoder   // 编解码器
	log       *zap.SugaredLogger // log
}

// GetCrossProcessHandler return the instance of CrossProcessHandler
func GetCrossProcessHandler() *CrossProcessHandler {
	return crossProcessHandler
}

// SetEventChan set event channel
func (c *CrossProcessHandler) SetEventChan(eventChan chan event.Event) {
	c.eventChan = eventChan
}

// SetStateDB set state database
func (c *CrossProcessHandler) SetStateDB(stateDB store.StateDB) {
	c.stateDB = stateDB
}

// SetLogger set logger
func (c *CrossProcessHandler) SetLogger(logger *zap.SugaredLogger) {
	c.log = logger
}

// GetType return type of this handler
func (c *CrossProcessHandler) GetType() HandlerType {
	return CrossProcess
}

// Handle handle event which input it to event channel
func (c *CrossProcessHandler) Handle(eve event.Event, _ bool) (interface{}, error) {
	// 接收到的事件是事务事件，该事件需要传递到innerRouter来处理
	eveTy := eve.GetType()
	if eveTy != eventproto.CrossEventType {
		c.log.Errorf("can not support this event for [%v]", eveTy)
		return nil, fmt.Errorf("can not support this event for [%v]", eveTy)
	}
	// 进行强制类型转换
	if crossEvent, ok := eve.(*eventproto.CrossEvent); ok {
		c.log.Infof("receive cross event cross = %s", crossEvent.GetCrossID())
		// 放入channel即可
		c.eventChan <- crossEvent
		return nil, nil
	} else {
		return nil, errors.New("can not support this event")
	}
}
