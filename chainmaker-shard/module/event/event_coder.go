/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

// EventCoder retrieve the coder for each Event
type EventCoder interface {

	// GetEventType get event implements type
	GetEventType() eventproto.EventType

	// MarshalToBinary encode to binary
	MarshalToBinary(Event) ([]byte, error)

	// UnmarshalFromBinary decode from binary
	UnmarshalFromBinary([]byte) (Event, error)
}
