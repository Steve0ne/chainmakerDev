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

var crossTxCoder *CrossTxCoder

func init() {
	crossTxCoder = &CrossTxCoder{}
}

// GetCrossTxCoder return instance of cross tx coder
func GetCrossTxCoder() *CrossTxCoder {
	return crossTxCoder
}

// CrossTxCoder cross tx coder struct
type CrossTxCoder struct {
}

// GetEventType return event type of event coder
func (c *CrossTxCoder) GetEventType() eventproto.EventType {
	return eventproto.CrossTxType
}

// MarshalToBinary marshal event to binary data
func (c *CrossTxCoder) MarshalToBinary(eve event.Event) ([]byte, error) {
	eveTy := eve.GetType()
	if eveTy != c.GetEventType() {
		return nil, fmt.Errorf("can not support event type [%v]", eveTy)
	}
	if eve, ok := eve.(*eventproto.CrossTx); ok {
		return c.marshalToBinary(eve)
	} else {
		return nil, errors.New("can not parse to [event.CrossTx]")
	}
}

// UnmarshalFromBinary unmarshal to event from binary data
func (c *CrossTxCoder) UnmarshalFromBinary(bytes []byte) (event.Event, error) {
	var eveObject = &eventproto.CrossTx{}
	if err := JsonBinaryUnmarshal(bytes, byte(c.GetEventType()), eveObject); err != nil {
		return nil, err
	}
	return eveObject, nil
}

func (c *CrossTxCoder) marshalToBinary(crossTx *eventproto.CrossTx) ([]byte, error) {
	return JsonBinaryMarshal(c.GetEventType(), crossTx)
}
