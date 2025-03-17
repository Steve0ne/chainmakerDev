/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package coder

import (
	"sync"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
)

var tools *EventCoderTools

func init() {
	tools = &EventCoderTools{
		coders: make(map[eventproto.EventType]*EventCoders),
	}
	tools.InitEventCoder(eventproto.CrossEventType, GetCrossEventCoder())
	tools.InitEventCoder(eventproto.TransactionEventType, GetTransactionEventCoder())
	tools.InitEventCoder(eventproto.CrossTxType, GetCrossTxCoder())
	tools.InitEventCoder(eventproto.CrossRespEventType, GetCrossRespEventCoder())
	tools.InitEventCoder(eventproto.ProofRespEventType, GetProofRespEventCoder())
	tools.InitEventCoder(eventproto.TransactionCtxEventType, GetTransactionEventCtxCoder())
	tools.InitEventCoder(eventproto.TxProofType, GetTransactionProofCoder())
}

// GetEventCoderTools return instance of event coder tools
func GetEventCoderTools() *EventCoderTools {
	return tools
}

// EventCoderTools event coder tools struct
type EventCoderTools struct {
	sync.RWMutex                                       // 读写锁
	coders       map[eventproto.EventType]*EventCoders // 编解码器的Map
}

// InitEventCoder register link between event-type and event coder
func (tools *EventCoderTools) InitEventCoder(eventType eventproto.EventType, coder event.EventCoder) {
	tools.coders[eventType] = NewEventCoders(coder)
}

// GetCoders return coder by event-type
func (tools *EventCoderTools) GetCoders(eventType eventproto.EventType) (*EventCoders, bool) {
	tools.RLock()
	defer tools.RUnlock()
	return tools.getCoders(eventType)
}

// getCoders non-use lock
func (tools *EventCoderTools) getCoders(eventType eventproto.EventType) (*EventCoders, bool) {
	eventCoders, exist := tools.coders[eventType]
	return eventCoders, exist
}

// GetDefaultCoder return default coder by event-type
func (tools *EventCoderTools) GetDefaultCoder(eventType eventproto.EventType) (event.EventCoder, bool) {
	tools.RLock()
	defer tools.RUnlock()
	return tools.getCoder(eventType, "")
}

// GetCoder return coder by event-type and chain-id
func (tools *EventCoderTools) GetCoder(eventType eventproto.EventType, chainID string) (event.EventCoder, bool) {
	tools.RLock()
	defer tools.RUnlock()
	return tools.getCoder(eventType, chainID)
}

func (tools *EventCoderTools) getCoder(eventType eventproto.EventType, chainID string) (event.EventCoder, bool) {
	if eventCoders, exist := tools.getCoders(eventType); exist {
		if chainID == "" {
			return eventCoders.GetDefaultCoder()
		}
		return eventCoders.GetCoder(chainID)
	}
	return nil, false
}

//EventCoders event coder struct
type EventCoders struct {
	sync.RWMutex                                  // lock
	eventTy           eventproto.EventType        // 事件类型
	defaultEventCoder event.EventCoder            // 默认编解码器
	cs                map[string]event.EventCoder // 编解码器实例
}

// NewEventCoders create new event coders
func NewEventCoders(defaultCoder event.EventCoder) *EventCoders {
	coders := &EventCoders{
		eventTy:           defaultCoder.GetEventType(),
		defaultEventCoder: defaultCoder,
		cs:                make(map[string]event.EventCoder),
	}
	return coders
}

// RegisterEventCoder register event coder by chain-id and event coder
func (coders *EventCoders) RegisterEventCoder(chainID string, coder event.EventCoder) {
	coders.Lock()
	defer coders.Unlock()
	coders.cs[chainID] = coder
}

// GetCoder return coder by chain-id
func (coders *EventCoders) GetCoder(chainID string) (event.EventCoder, bool) {
	coders.RLock()
	defer coders.RUnlock()
	eventCoder, exist := coders.cs[chainID]
	return eventCoder, exist
}

// GetDefaultCoder return default coder
func (coders *EventCoders) GetDefaultCoder() (event.EventCoder, bool) {
	coders.RLock()
	defer coders.RUnlock()
	return coders.defaultEventCoder, coders.defaultEventCoder != nil
}
