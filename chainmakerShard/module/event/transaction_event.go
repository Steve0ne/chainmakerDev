/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
)

//type OpFuncType int32

const (
	ExecuteOpFunc  = eventproto.OpFuncType_ExecuteOpFunc
	CommitOpFunc   = eventproto.OpFuncType_CommitOpFunc
	RollbackOpFunc = eventproto.OpFuncType_RollbackOpFunc
)

// NewExecuteTransactionEvent create new execute transaction event
func NewExecuteTransactionEvent(crossID, chainID string, payload []byte, proofKey string, txProof *eventproto.Proof) *eventproto.TransactionEvent {
	return &eventproto.TransactionEvent{
		CrossId: crossID,
		OpFunc:  ExecuteOpFunc,
		ChainId: chainID,
		Payload: payload,
		ProofKey: proofKey,
		TxProof: txProof,
	}
}

// NewRollbackTransactionEvent create new rollback transaction event
func NewRollbackTransactionEvent(crossID string, chainID string, payload []byte) *eventproto.TransactionEvent {
	return &eventproto.TransactionEvent{
		CrossId: crossID,
		OpFunc:  RollbackOpFunc,
		ChainId: chainID,
		Payload: payload,
	}
}

// NewCommitTransactionEvent create new commit transaction event
func NewCommitTransactionEvent(crossID string, chainID string, payload []byte) *eventproto.TransactionEvent {
	return &eventproto.TransactionEvent{
		CrossId: crossID,
		OpFunc:  CommitOpFunc,
		ChainId: chainID,
		Payload: payload,
	}
}
