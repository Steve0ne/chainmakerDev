/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package prover

import (
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	"chainmaker.org/chainmaker-cross/prover/impl"
)

// Prover is the interface to prove proof
type Prover interface {

	// GetType return the type of prover
	GetType() impl.ProverType

	// GetChainIDs return all the chain ids
	GetChainIDs() []string

	// ToProof convert to Proof for the inputs
	ToProof(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) (*eventproto.Proof, error)

	// Prove load result of this proof
	Prove(proof *eventproto.Proof) (bool, error)
}
