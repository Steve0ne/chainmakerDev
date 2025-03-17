/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

// Proofer interface of proof
type Proofer interface {

	// GetChainID return chain-id
	GetChainID() string

	// GetTxKey return the key of proof
	GetTxKey() string

	// GetBlockHeight return the height of proof
	GetBlockHeight() int64

	// GetIndex return index of tx
	GetIndex() int

	// GetContract return contract info
	GetContract() *eventproto.ContractInfo

	// GetExtra return extra data
	GetExtra() []byte
}

// NewProof create new proof
func NewProof(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) *eventproto.Proof {
	return &eventproto.Proof{
		ChainId:     chainID,
		TxKey:       txKey,
		BlockHeight: blockHeight,
		Index:       index,
		Contract:    contract,
		Extra:       extra,
	}
}
