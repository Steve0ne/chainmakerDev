/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"sync"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/utils"
)

var proofResponseContexts *ProofResponseContexts

func init() {
	proofResponseContexts = &ProofResponseContexts{
		contexts: make(map[string]*ProofResponseContext),
	}
}

// ProofResponseContexts set of proof response context
type ProofResponseContexts struct {
	sync.RWMutex
	contexts map[string]*ProofResponseContext
}

// GetProofResponseContexts return instance of proof response contexts
func GetProofResponseContexts() *ProofResponseContexts {
	return proofResponseContexts
}

// Register add context to contexts
func (ctxs *ProofResponseContexts) Register(ctx *ProofResponseContext) bool {
	ctxs.Lock()
	defer ctxs.Unlock()
	key := ctx.GetKey()
	if _, exist := ctxs.contexts[key]; !exist {
		ctxs.contexts[key] = ctx
		return true
	}
	return false
}

// DoneError set state to error
func (ctxs *ProofResponseContexts) DoneError(key, msg string) bool {
	ctx, exist := ctxs.getContext(key)
	if exist {
		return ctx.DoneError(msg)
	}
	return false
}

// DoneByProofResp update state by proof response
func (ctxs *ProofResponseContexts) DoneByProofResp(proofResp *ProofResponse) {
	ctxs.Done(proofResp.Key, proofResp.TxResponse.ChainId, proofResp.TxResponse.TxKey, proofResp.TxResponse.BlockHeight, proofResp.TxResponse.Index, proofResp.TxResponse.Contract, proofResp.TxResponse.Extra)
}

// Done update state to done
func (ctxs *ProofResponseContexts) Done(key, chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) bool {
	ctx, exist := ctxs.getContext(key)
	if exist {
		return ctx.Done(chainID, txKey, blockHeight, index, contract, extra)
	}
	return false
}

// Remove remove context which key = {key}
func (ctxs *ProofResponseContexts) Remove(key string) {
	ctxs.Lock()
	defer ctxs.Unlock()
	delete(ctxs.contexts, key)
}

func (ctxs *ProofResponseContexts) getContext(key string) (*ProofResponseContext, bool) {
	ctxs.RLock()
	defer ctxs.RUnlock()
	if ctx, exist := ctxs.contexts[key]; exist {
		return ctx, true
	}
	return nil, false
}

// NewProofResponseContext create new proof response context
func NewProofResponseContext(resp *ProofResponse) *ProofResponseContext {
	key := utils.NewRandomKey()
	resp.SetKey(key) // 设置ProofResponse的Key
	return &ProofResponseContext{
		key:  key,
		resp: resp,
	}
}

// ProofResponseContext the struct of proof response context
type ProofResponseContext struct {
	sync.Mutex                // lock
	key        string         // 随机的key
	resp       *ProofResponse // 验证消息
	completed  bool           // 验证完成的标签
}

// GetKey return key of context
func (ctx *ProofResponseContext) GetKey() string {
	return ctx.key
}

// DoneError set context state to error
func (ctx *ProofResponseContext) DoneError(msg string) bool {
	ctx.Lock()
	defer ctx.Unlock()
	if ctx.completed {
		return false
	}
	ctx.resp.DoneError(msg)
	ctx.completed = true
	return true
}

// Done set state to success and set tx's info in chain
func (ctx *ProofResponseContext) Done(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) bool {
	ctx.Lock()
	defer ctx.Unlock()
	if ctx.completed {
		return false
	}
	ctx.resp.Done(chainID, txKey, blockHeight, index, contract, extra)
	ctx.completed = true
	return true
}

//TransactionEventContext struct of transaction event context
type TransactionEventContext struct {
	Key   string
	Event *eventproto.TransactionEvent // 跨链消息
}

// NewTransactionEventContext create new transaction event context
func NewTransactionEventContext(key string, eve *eventproto.TransactionEvent) *TransactionEventContext {
	return &TransactionEventContext{
		Key:   key,
		Event: eve,
	}
}

// GetType return type of this event
func (ctx *TransactionEventContext) GetType() eventproto.EventType {
	return eventproto.TransactionCtxEventType
}

// GetKey return the key of context
func (ctx *TransactionEventContext) GetKey() string {
	return ctx.Key
}

//GetEvent return event in context
func (ctx *TransactionEventContext) GetEvent() *eventproto.TransactionEvent {
	return ctx.Event
}
