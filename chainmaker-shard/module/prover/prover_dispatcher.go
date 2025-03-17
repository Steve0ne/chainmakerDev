/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package prover

import (
	"fmt"
	"sync"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
)

var dispatcher *ProverDispatcher

func init() {
	dispatcher = &ProverDispatcher{
		provers: make(map[string]Prover),
	}
}

// ProverDispatcher is dispatcher of prover
type ProverDispatcher struct {
	sync.RWMutex
	provers map[string]Prover // Prover 的map， key 为chainID
}

// GetProverDispatcher return the instance of ProverDispatcher
func GetProverDispatcher() *ProverDispatcher {
	return dispatcher
}

// Register put the prover to map
func (pd *ProverDispatcher) Register(prover Prover) {
	pd.Lock()
	defer pd.Unlock()
	chainIDs := prover.GetChainIDs()
	for _, chainID := range chainIDs {
		if _, exist := pd.provers[chainID]; !exist {
			pd.provers[chainID] = prover
		}
	}
}

// ToProof convert to Proof for the inputs
func (pd *ProverDispatcher) ToProof(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) (*eventproto.Proof, error) {
	if prover, exist := pd.GetProver(chainID); exist {
		return prover.ToProof(chainID, txKey, blockHeight, index, contract, extra)
	}
	// 不存在的话走正常的处理逻辑
	return event.NewProof(chainID, txKey, blockHeight, index, contract, extra), nil
}

// Prove load result of this proof
func (pd *ProverDispatcher) Prove(proof *eventproto.Proof) (bool, error) {
	chainID := proof.GetChainID()
	if prover, exist := pd.GetProver(chainID); exist {
		return prover.Prove(proof)
	}
	return false, fmt.Errorf("can not find prover for chainID [%v]", chainID)
}

// GetProver return prover by chain id
func (pd *ProverDispatcher) GetProver(chainID string) (Prover, bool) {
	pd.RLock()
	defer pd.RUnlock()
	prover, exist := pd.provers[chainID]
	return prover, exist
}
