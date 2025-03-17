/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"sync"
)

var handlerTools *EventHandlerTools

func init() {
	handlerTools = &EventHandlerTools{
		handlers: make(map[HandlerType]EventHandler),
	}
}

// EventHandlerTools the struct of event handler tools
type EventHandlerTools struct {
	sync.RWMutex
	handlers map[HandlerType]EventHandler
}

// GetEventHandlerTools return the instance of event handler tools
func GetEventHandlerTools() *EventHandlerTools {
	return handlerTools
}

// Register add event handler to event handler tools
func (d *EventHandlerTools) Register(handler EventHandler) {
	d.Lock()
	defer d.Unlock()
	ty := handler.GetType()
	d.handlers[ty] = handler
}

// GetHandler return the instance of event handler by type
func (d *EventHandlerTools) GetHandler(handlerType HandlerType) (EventHandler, bool) {
	d.RLock()
	defer d.RUnlock()
	handler, exist := d.handlers[handlerType]
	return handler, exist
}
