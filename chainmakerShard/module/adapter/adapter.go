/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package adapter

import (
	"chainmaker.org/chainmaker-cross/event"
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
)

// ChainAdapter adapter of chain interface
type ChainAdapter interface {

	// GetChainID return chain id
	GetChainID() string

	// Invoke invoke transaction event and return response
	Invoke(txEvent *eventproto.TransactionEvent) (*eventproto.TxResponse, error)

	// SaveProof save the proof and verify in the chain
	SaveProof(crossID, proofKey string, txProof *eventproto.Proof, verifyResult bool) (*eventproto.TxResponse, error)

	// Prove prove the proof
	Prove(txProof *eventproto.Proof) bool

	// QueryByTxKey query tx response by tx-key
	QueryByTxKey(txKey string) (*event.CommonTxResponse, error)

	// QueryTx query tx and return tx response
	QueryTx(payload []byte) (*event.CommonTxResponse, error)
}
