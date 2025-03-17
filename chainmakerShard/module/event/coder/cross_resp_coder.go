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

var crossRespEventCoder *CrossRespEventCoder

func init() {
	crossRespEventCoder = &CrossRespEventCoder{}
}

// GetCrossRespEventCoder return instance of cross response event coder
func GetCrossRespEventCoder() *CrossRespEventCoder {
	return crossRespEventCoder
}

// CrossRespEventCoder cross response event coder struct
type CrossRespEventCoder struct {
}

// GetEventType return event type of event coder
func (c *CrossRespEventCoder) GetEventType() eventproto.EventType {
	return eventproto.CrossRespEventType
}

// MarshalToBinary marshal event to binary data
func (c *CrossRespEventCoder) MarshalToBinary(eve event.Event) ([]byte, error) {
	eveTy := eve.GetType()
	if eveTy != c.GetEventType() {
		return nil, fmt.Errorf("can not support event type [%v]", eveTy)
	}
	if eve, ok := eve.(*eventproto.CrossResponse); ok {
		return c.marshalToBinary(eve)
	} else {
		return nil, errors.New("can not parse to [event.CrossResponse]")
	}
}

// UnmarshalFromBinary unmarshal to event from binary data
func (c *CrossRespEventCoder) UnmarshalFromBinary(bytes []byte) (event.Event, error) {
	var eveObject = &eventproto.CrossResponse{}
	if err := JsonBinaryUnmarshal(bytes, byte(c.GetEventType()), eveObject); err != nil {
		return nil, err
	}
	return eveObject, nil
}

func (c *CrossRespEventCoder) marshalToBinary(crossEvent *eventproto.CrossResponse) ([]byte, error) {
	return JsonBinaryMarshal(c.GetEventType(), crossEvent)
}
