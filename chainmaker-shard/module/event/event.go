/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	"go.uber.org/zap"
)

type MarshalType byte

const (
	BinaryMarshalType MarshalType = 0
	JsonMarshalType   MarshalType = 1
)

var log *zap.SugaredLogger

// InitLog init instance of log
func InitLog(zapLog *zap.SugaredLogger) {
	log = zapLog
}

// Event event interface of all instance event
type Event interface {
	// 返回消息实例类型
	GetType() eventproto.EventType
}

// NewCrossTx create new cross tx
func NewCrossTx(chainID string, index int32, executePayload, commitPayload, rollbackPayload []byte) *eventproto.CrossTx {
	return &eventproto.CrossTx{
		ChainId:         chainID,
		Index:           index,
		ExecutePayload:  executePayload,
		CommitPayload:   commitPayload,
		RollbackPayload: rollbackPayload,
	}
}

// NewCrossTxs create new CrossTxs instance by array of CrossTx
func NewCrossTxs(events []*eventproto.CrossTx) *eventproto.CrossTxs {
	return &eventproto.CrossTxs{
		Events: events,
	}
}
