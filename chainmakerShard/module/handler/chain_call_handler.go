/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"errors"
	"fmt"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/adapter"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/store"
	"go.uber.org/zap"
)

var chainCallHandler *ChainCallHandler

func init() {
	chainCallHandler = &ChainCallHandler{
		dispatcher:   adapter.GetChainAdapterDispatcher(),
		respContexts: event.GetProofResponseContexts(),
	}
}

// GetChainCallHandler return instance of ChainCallHandler
func GetChainCallHandler() *ChainCallHandler {
	return chainCallHandler
}

// ChainCallHandler the struct of ChainCallHandler
type ChainCallHandler struct {
	dispatcher   *adapter.ChainAdapterDispatcher // 链交互的转接器的Map
	respContexts *event.ProofResponseContexts    // 根据转接返回结果，提供验证信息的Map
	db           store.StateDB                   // 存储
	log          *zap.SugaredLogger              // log
}

// SetStateDB set state database
func (c *ChainCallHandler) SetStateDB(db store.StateDB) {
	c.db = db
}

// SetLogger set logger of chain call handler
func (c *ChainCallHandler) SetLogger(log *zap.SugaredLogger) {
	c.log = log
}

// GetType return type of this handler
func (c *ChainCallHandler) GetType() HandlerType {
	return ChainCall
}

// Handle handle event which send this message to real service chain
func (c *ChainCallHandler) Handle(eve event.Event, syncWait bool) (interface{}, error) {
	eveTy := eve.GetType()
	if eveTy != eventproto.TransactionCtxEventType {
		c.log.Errorf("can not support this event for [%v]", eveTy)
		return nil, fmt.Errorf("can not support this event for [%v]", eveTy)
	}
	// 进行强制类型转换
	if txCtxEvent, ok := eve.(*event.TransactionEventContext); ok {
		// 类型转换正确，开始处理
		txEvent := txCtxEvent.GetEvent()
		chainID := txEvent.GetChainID()
		key := txCtxEvent.GetKey()
		c.log.Infof("chain[%s]->key[%s] start handle event", chainID, key)
		resp, err := c.dispatcher.Invoke(chainID, txEvent)
		if err != nil {
			// 出错，填充到contexts中
			c.respContexts.DoneError(key, err.Error())
		} else if resp == nil {
			c.respContexts.DoneError(key, "get empty invoke response")
		} else {
			// 成功，填充结果
			c.respContexts.Done(key, resp.GetChainID(), resp.TxKey, resp.BlockHeight, resp.Index, resp.Contract, resp.Extra)
		}
		return nil, nil
	}
	return nil, errors.New("can not support this event")
}
