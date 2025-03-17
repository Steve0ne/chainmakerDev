/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package impl

import (
	"fmt"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker/spv/v2/pb/api"
	spvserver "chainmaker.org/chainmaker/spv/v2/server"
	"go.uber.org/zap"
)

const (
	ProveMaxWait = 5000 //ms
)

// SpvProver is prover of spv
type SpvProver struct {
	chainIDs  []string             // spv 支持的chainID
	spvServer *spvserver.SPVServer // spv 服务
	log       *zap.SugaredLogger   // log
}

// NewSpvProver create new instance of SpvProver
func NewSpvProver(ymlFile string, chainIDs []string) *SpvProver {
	log := logger.GetLogger(logger.ModuleProver)
	spvServer, err := spvserver.NewSPVServer(ymlFile, log)
	if err != nil {
		panic(fmt.Errorf("create spv-prover failed, %v", err))
	}
	err = spvServer.Start()
	if err != nil {
		panic(fmt.Errorf("spv-prover start failed, %v", err))
	}
	return &SpvProver{
		chainIDs:  chainIDs,
		spvServer: spvServer,
		log:       log,
	}
}

// GetType return type of prover
func (s *SpvProver) GetType() ProverType {
	return SpvProverType
}

// GetChainIDs return chain-ids
func (s *SpvProver) GetChainIDs() []string {
	return s.chainIDs
}

// ToProof convert to Proof for the inputs
func (s *SpvProver) ToProof(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) (*eventproto.Proof, error) {
	return event.NewProof(chainID, txKey, blockHeight, index, contract, extra), nil
}

// Prove return if is ok for the proof
func (s *SpvProver) Prove(proof *eventproto.Proof) (bool, error) {
	if proof.GetContract() == nil {
		return false, fmt.Errorf("get proof event contract error, is nil")
	}
	if proof.Contract == nil {
		return false, fmt.Errorf("contract is nil")
	}
	// 处理parameters
	parameters := make([]*api.KVPair, 0)
	contractParameters := proof.Contract.Parameters
	for i := 0; i < len(contractParameters); i++ {
		param := contractParameters[i]
		parameters = append(parameters, &api.KVPair{
			Key:   param.Key,
			Value: []byte(param.Value),
		})
	}
	contractData := &api.ContractData{
		Name:    proof.GetContract().Name,
		Version: proof.GetContract().Version,
		Method:  proof.GetContract().Method,
		Params:  parameters,
		Extra:   proof.Extra,
	}
	txVerifyInfo := &api.TxValidationRequest{
		ChainId:      proof.GetChainID(),
		BlockHeight:  uint64(proof.GetBlockHeight()),
		Index:        int64(proof.GetIndex()),
		TxKey:        proof.GetTxKey(),
		ContractData: contractData,
		Extra:        proof.GetExtra(),
		Timeout:      ProveMaxWait,
	}

	err := s.spvServer.ValidTransaction(txVerifyInfo)
	if err != nil {
		s.log.Errorf("[%s][%s][%s][%s] prove failed", proof.GetContract().Name, proof.GetContract().Version, proof.GetContract().Method, err)
		return false, err
	}
	return true, nil
}
