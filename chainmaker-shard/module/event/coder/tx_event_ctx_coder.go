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

var transactionEventCtxCoder *TransactionEventCtxCoder

func init() {
	transactionEventCtxCoder = &TransactionEventCtxCoder{}
}

// GetTransactionEventCtxCoder return instance of transaction event ctx coder
func GetTransactionEventCtxCoder() *TransactionEventCtxCoder {
	return transactionEventCtxCoder
}

// TransactionEventCtxCoder transaction event ctx coder struct
type TransactionEventCtxCoder struct {
}

// GetEventType return event type of event coder
func (c *TransactionEventCtxCoder) GetEventType() eventproto.EventType {
	return eventproto.TransactionCtxEventType
}

// MarshalToBinary marshal event to binary data
func (c *TransactionEventCtxCoder) MarshalToBinary(eve event.Event) ([]byte, error) {
	eveTy := eve.GetType()
	if eveTy != c.GetEventType() {
		return nil, fmt.Errorf("can not support event type [%v]", eveTy)
	}
	if eve, ok := eve.(*event.TransactionEventContext); ok {
		return c.marshalToBinary(eve)
	} else {
		return nil, errors.New("can not parse to [event.TransactionEventContext]")
	}
}

// UnmarshalFromBinary unmarshal to event from binary data
func (c *TransactionEventCtxCoder) UnmarshalFromBinary(bytes []byte) (event.Event, error) {
	var eveObject = &event.TransactionEventContext{}
	if err := JsonBinaryUnmarshal(bytes, byte(c.GetEventType()), eveObject); err != nil {
		return nil, err
	}
	return eveObject, nil
}

func (c *TransactionEventCtxCoder) marshalToBinary(crossEvent *event.TransactionEventContext) ([]byte, error) {
	return JsonBinaryMarshal(c.GetEventType(), crossEvent)
}
