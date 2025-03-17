/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package coder

import (
	"errors"
	"fmt"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
)

var crossEventCoder *CrossEventCoder

func init() {
	crossEventCoder = &CrossEventCoder{}
}

// GetCrossEventCoder return instance of cross event coder
func GetCrossEventCoder() *CrossEventCoder {
	return crossEventCoder
}

// CrossEventCoder cross event coder struct
type CrossEventCoder struct {
}

// GetEventType return event type of event coder
func (c *CrossEventCoder) GetEventType() eventproto.EventType {
	return eventproto.CrossEventType
}

// MarshalToBinary marshal event to binary data
func (c *CrossEventCoder) MarshalToBinary(eve event.Event) ([]byte, error) {
	eveTy := eve.GetType()
	if eveTy != c.GetEventType() {
		return nil, fmt.Errorf("can not support event type [%v]", eveTy)
	}
	if eve, ok := eve.(*eventproto.CrossEvent); ok {
		return c.marshalToBinary(eve)
	} else {
		return nil, errors.New("can not parse to [event.CrossEvent]")
	}
}

// UnmarshalFromBinary unmarshal to event from binary data
func (c *CrossEventCoder) UnmarshalFromBinary(bytes []byte) (event.Event, error) {
	var crossEvent = &eventproto.CrossEvent{}
	if err := JsonBinaryUnmarshal(bytes, byte(c.GetEventType()), crossEvent); err != nil {
		return nil, err
	}
	return crossEvent, nil
}

func (c *CrossEventCoder) marshalToBinary(crossEvent *eventproto.CrossEvent) ([]byte, error) {
	return JsonBinaryMarshal(c.GetEventType(), crossEvent)
}
