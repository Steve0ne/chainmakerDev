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

var transactionProofCoder *TransactionProofCoder

func init() {
	transactionProofCoder = &TransactionProofCoder{}
}

// GetTransactionProofCoder return instance of transaction proof coder
func GetTransactionProofCoder() *TransactionProofCoder {
	return transactionProofCoder
}

// TransactionEventCtxCoder transaction proof coder struct
type TransactionProofCoder struct {
}

// GetEventType return event type of event coder
func (c *TransactionProofCoder) GetEventType() eventproto.EventType {
	return eventproto.TxProofType
}

// MarshalToBinary marshal event to binary data
func (c *TransactionProofCoder) MarshalToBinary(eve event.Event) ([]byte, error) {
	eveTy := eve.GetType()
	if eveTy != c.GetEventType() {
		return nil, fmt.Errorf("can not support event type [%v]", eveTy)
	}
	if eve, ok := eve.(*eventproto.Proof); ok {
		return c.marshalToBinary(eve)
	} else {
		return nil, errors.New("can not parse to [event.TransactionEvent]")
	}
}

// UnmarshalFromBinary unmarshal to event from binary data
func (c *TransactionProofCoder) UnmarshalFromBinary(bytes []byte) (event.Event, error) {
	var eveObject = &eventproto.Proof{}
	if err := JsonBinaryUnmarshal(bytes, byte(c.GetEventType()), eveObject); err != nil {
		return nil, err
	}
	return eveObject, nil
}

func (c *TransactionProofCoder) marshalToBinary(crossEvent *eventproto.Proof) ([]byte, error) {
	return JsonBinaryMarshal(c.GetEventType(), crossEvent)
}
