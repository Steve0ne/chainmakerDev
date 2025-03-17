/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/utils"
)

const (
	DefaultVersion = "v0.9.0"
)

// NewCrossEvent create cross event by array of CrossTx
func NewCrossEvent(txEvents []*eventproto.CrossTx) *eventproto.CrossEvent {
	crossID := utils.NewUUID()
	return &eventproto.CrossEvent{
		CrossId:   crossID,
		Version:   DefaultVersion,
		Timestamp: time.Now().Unix(),
		TxEvents:  NewCrossTxs(txEvents),
	}
}

// NewEmptyCrossEvent create empty cross event which have not any CrossTx
func NewEmptyCrossEvent() *eventproto.CrossEvent {
	crossID := utils.NewUUID()
	return &eventproto.CrossEvent{
		CrossId:   crossID,
		Version:   DefaultVersion,
		Timestamp: time.Now().Unix(),
		TxEvents:  &eventproto.CrossTxs{},
	}
}

// NewCrossSearchEvent create cross search event
func NewCrossSearchEvent(crossID string) *eventproto.CrossSearchEvent {
	return &eventproto.CrossSearchEvent{
		CrossId: crossID,
	}
}
