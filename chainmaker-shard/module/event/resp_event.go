/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"sync"
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
)

const (
	SuccessResp = iota
	FailureResp
	ErrorResp
	UnknownResp
)

// NewCrossResponse create new cross response object
func NewCrossResponse(crossID string, code int32, msg string) *eventproto.CrossResponse {
	return &eventproto.CrossResponse{
		CrossId:     crossID,
		Code:        code,
		Msg:         msg,
		TxResponses: make([]*eventproto.CrossTxResponse, 0),
	}
}

// DefaultCrossResponse return new cross response which include empty cross tx response
func DefaultCrossResponse() *eventproto.CrossResponse {
	return &eventproto.CrossResponse{
		TxResponses: make([]*eventproto.CrossTxResponse, 0),
	}
}

// NewCrossTxResponse create new tx response
func NewCrossTxResponse(chainID, txKey string, blockHeight int64, index int32, extra []byte) *eventproto.CrossTxResponse {
	return &eventproto.CrossTxResponse{
		ChainId:     chainID,
		TxKey:       txKey,
		BlockHeight: blockHeight,
		Index:       int32(index),
		Extra:       extra,
	}
}

// NewTxResponse create new tx response
func NewTxResponse(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) *eventproto.TxResponse {
	return &eventproto.TxResponse{
		ChainId:     chainID,
		TxKey:       txKey,
		BlockHeight: blockHeight,
		Index:       index,
		Contract:    contract,
		Extra:       extra,
	}
}

// CommonTxResponse the struct of common tx response
type CommonTxResponse struct {
	eventproto.TxResponse
	Code int32
	Msg  string
}

// NewCommonTxResponse create new common tx response
func NewCommonTxResponse(txResponse *eventproto.TxResponse, code int32, msg string) *CommonTxResponse {
	resp := &CommonTxResponse{
		Code: code,
		Msg:  msg,
	}
	resp.TxResponse = *txResponse
	return resp
}

// IsSuccess return whether is success
func (c *CommonTxResponse) IsSuccess() bool {
	return c.Code == SuccessResp
}

// NewContract create new contract
func NewContract(name, version, method string, extraData []byte) *eventproto.ContractInfo {
	return &eventproto.ContractInfo{
		Name:       name,
		Version:    version,
		Method:     method,
		Parameters: make([]*eventproto.ContractParameter, 0),
		ExtraData:  extraData,
	}
}

// NewContractParameter create new contract parameter
func NewContractParameter(key, value string) *eventproto.ContractParameter {
	return &eventproto.ContractParameter{
		Key:   key,
		Value: value,
	}
}

// NewContractValue create new ContractParameter
func NewContractValue(value string) *eventproto.ContractParameter {
	return NewContractParameter("", value)
}

// ProofResponse the struct of proof response
type ProofResponse struct {
	sync.Mutex // lock
	eventproto.ProofResponse
	ch          chan bool // 同步/异步通道
	isCompleted bool      // 是否完成的标签
}

// NewProofResponse create new proof response
func NewProofResponse(crossID, chainID string, opFunc eventproto.OpFuncType) *ProofResponse {
	pr := &ProofResponse{
		ProofResponse: eventproto.ProofResponse{
			CrossId: crossID,
			OpFunc:  opFunc,
			Code:    UnknownResp,
			TxResponse: &eventproto.TxResponse{
				ChainId: chainID,
			},
		},
		ch: make(chan bool, 1),
	}
	//pr.ChainId = chainID
	return pr
}

// NewProofResponseByProof create new proof response by proof
func NewProofResponseByProof(crossID, chainID, msg string, code int32, opFunc eventproto.OpFuncType, proof *eventproto.Proof) *ProofResponse {
	proofResp := NewProofResponse(crossID, chainID, opFunc)
	proofResp.TxResponse = NewTxResponse(proof.ChainId, proof.TxKey, proof.BlockHeight, proof.Index, proof.Contract, proof.Extra)
	proofResp.Msg = msg
	proofResp.Code = code
	return proofResp
}

// SetKey set key of proof response
func (p *ProofResponse) SetKey(key string) {
	p.Key = key
}

// GetKey return key of proof response
func (p *ProofResponse) GetKey() string {
	return p.Key
}

// GetType return type of this event
func (p *ProofResponse) GetType() eventproto.EventType {
	return eventproto.ProofRespEventType
}
func (p *ProofResponse) GetCrossID() string {
	return p.GetCrossId()
}
func (p *ProofResponse) SetChainID(chainID string) {
	if p.TxResponse == nil {
		p.TxResponse = &eventproto.TxResponse{}
	}
	p.TxResponse.ChainId = chainID
}
func (p *ProofResponse) GetChainID() string {
	return p.TxResponse.GetChainId()
}
func (p *ProofResponse) GetTxKey() string {
	return p.TxResponse.GetTxKey()
}

func (p *ProofResponse) GetBlockHeight() int64 {
	return p.TxResponse.GetBlockHeight()
}

func (p *ProofResponse) GetIndex() int32 {
	return p.TxResponse.GetIndex()
}

func (p *ProofResponse) GetContract() *eventproto.ContractInfo {
	return p.TxResponse.GetContract()
}

func (p *ProofResponse) GetExtra() []byte {
	return p.TxResponse.GetExtra()
}

// IsSuccess return whether is success
func (p *ProofResponse) IsSuccess() bool {
	return p.Code == SuccessResp
}

// DoneError update state to failed
func (p *ProofResponse) DoneError(msg string) {
	p.Lock()
	defer p.Unlock()
	if !p.isCompleted {
		p.Code, p.Msg = FailureResp, msg
		p.isCompleted = true
		p.ch <- true
	}
}

// Done update state and set tx-info
func (p *ProofResponse) Done(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) {
	p.Lock()
	defer p.Unlock()
	if !p.isCompleted {
		p.Code = SuccessResp
		p.TxResponse.ChainId, p.TxResponse.TxKey = chainID, txKey
		p.TxResponse.BlockHeight, p.TxResponse.Index = blockHeight, index
		p.TxResponse.Contract, p.TxResponse.Extra = contract, extra
		p.isCompleted = true
		p.ch <- true
	}
}

// Wait wait until {waitTime} or completed
func (p *ProofResponse) Wait(waitTime time.Duration) bool {
	// 等待ch处理完成
	select {
	case <-time.After(waitTime):
		log.Warnf("waiting response timeout for [%s]", p.CrossId)
	case <-p.ch:
		log.Infof("waiting response success for [%s]", p.CrossId)
	}
	return p.isCompleted
}
