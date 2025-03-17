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

var transactionEventCoder *TransactionEventCoder

func init() {
	transactionEventCoder = &TransactionEventCoder{}
}

// GetTransactionEventCoder return instance of transaction event coder
func GetTransactionEventCoder() *TransactionEventCoder {
	return transactionEventCoder
}

// TransactionEventCoder transaction event coder struct
type TransactionEventCoder struct {
}

// GetEventType return event type of event coder
func (c *TransactionEventCoder) GetEventType() eventproto.EventType {
	return eventproto.TransactionEventType
}

// MarshalToBinary marshal event to binary data
func (c *TransactionEventCoder) MarshalToBinary(eve event.Event) ([]byte, error) {
	eveTy := eve.GetType()
	if eveTy != c.GetEventType() {
		return nil, fmt.Errorf("can not support event type [%v]", eveTy)
	}
	if eve, ok := eve.(*eventproto.TransactionEvent); ok {
		return c.marshalToBinary(eve)
	} else {
		return nil, errors.New("can not parse to [event.TransactionEvent]")
	}
}

// UnmarshalFromBinary unmarshal to event from binary data
func (c *TransactionEventCoder) UnmarshalFromBinary(bytes []byte) (event.Event, error) {
	var eveObject = &eventproto.TransactionEvent{}
	if err := JsonBinaryUnmarshal(bytes, byte(c.GetEventType()), eveObject); err != nil {
		return nil, err
	}
	return eveObject, nil
}

func (c *TransactionEventCoder) marshalToBinary(crossEvent *eventproto.TransactionEvent) ([]byte, error) {
	return JsonBinaryMarshal(c.GetEventType(), crossEvent)
}
