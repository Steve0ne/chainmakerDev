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

var proofRespEventCoder *ProofRespEventCoder

func init() {
	proofRespEventCoder = &ProofRespEventCoder{}
}

// GetProofRespEventCoder return instance of proof response event coder
func GetProofRespEventCoder() *ProofRespEventCoder {
	return proofRespEventCoder
}

// ProofRespEventCoder proof response event coder struct
type ProofRespEventCoder struct {
}

// GetEventType return event type of event coder
func (c *ProofRespEventCoder) GetEventType() eventproto.EventType {
	return eventproto.ProofRespEventType
}

// MarshalToBinary marshal event to binary data
func (c *ProofRespEventCoder) MarshalToBinary(eve event.Event) ([]byte, error) {
	eveTy := eve.GetType()
	if eveTy != c.GetEventType() {
		return nil, fmt.Errorf("can not support event type [%v]", eveTy)
	}
	if eve, ok := eve.(*event.ProofResponse); ok {
		return c.marshalToBinary(eve)
	} else {
		return nil, errors.New("can not parse to [event.ProofResponse]")
	}
}

// UnmarshalFromBinary unmarshal to event from binary data
func (c *ProofRespEventCoder) UnmarshalFromBinary(bytes []byte) (event.Event, error) {
	var eveObject = &event.ProofResponse{}
	if err := JsonBinaryUnmarshal(bytes, byte(c.GetEventType()), eveObject); err != nil {
		return nil, err
	}
	return eveObject, nil
}

func (c *ProofRespEventCoder) marshalToBinary(crossEvent *event.ProofResponse) ([]byte, error) {
	return JsonBinaryMarshal(c.GetEventType(), crossEvent)
}
