/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package impl

import (
	"chainmaker.org/chainmaker-cross/event"
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
)

// TrustProver is a trust prover
type TrustProver struct {
	chainIDs []string
}

// NewTrustProver create new trust prover
func NewTrustProver(chainIDs []string) *TrustProver {
	return &TrustProver{
		chainIDs: chainIDs,
	}
}

// ToProof convert to Proof for the inputs
func (t *TrustProver) ToProof(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) (*eventproto.Proof, error) {
	return event.NewProof(chainID, txKey, blockHeight, index, contract, extra), nil
}

// Prove return true
func (t *TrustProver) Prove(proof *eventproto.Proof) (bool, error) {
	return true, nil
}

// GetType return type of prover
func (t *TrustProver) GetType() ProverType {
	return TrustProverType
}

// GetChainIDs return chain-ids
func (t *TrustProver) GetChainIDs() []string {
	return t.chainIDs
}
