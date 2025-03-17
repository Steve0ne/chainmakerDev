/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"errors"
	"fmt"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/event/coder"
	"chainmaker.org/chainmaker-cross/router"
	"chainmaker.org/chainmaker-cross/store"
	storetype "chainmaker.org/chainmaker-cross/store/types"
	"go.uber.org/zap"
)

var transactionProcessHandler *TransactionProcessHandler

func init() {
	transactionProcessHandler = &TransactionProcessHandler{
		dispatcher: router.GetDispatcher(),
		coders:     coder.GetEventCoderTools(),
	}
}

// TransactionProcessHandler the struct of transaction process handler
type TransactionProcessHandler struct {
	dispatcher *router.RouterDispatcher // 路由表，包含其他代理的长连接
	db         store.StateDB            // 存储
	log        *zap.SugaredLogger       // log
	coders     *coder.EventCoderTools   // 编解码器
}

// GetTransactionProcessHandler return the instance of TransactionProcessHandler
func GetTransactionProcessHandler() *TransactionProcessHandler {
	return transactionProcessHandler
}

// SetStateDB set state database
func (t *TransactionProcessHandler) SetStateDB(db store.StateDB) {
	t.db = db
}

// SetLogger set logger
func (t *TransactionProcessHandler) SetLogger(log *zap.SugaredLogger) {
	t.log = log
}

// GetType return the type of handler
func (t *TransactionProcessHandler) GetType() HandlerType {
	return TransactionProcess
}

// Handle handle event which will transfer event to other cross-chain proxy
func (t *TransactionProcessHandler) Handle(eve event.Event, _ bool) (interface{}, error) {
	// 接收到的事件是事务事件，该事件需要传递到innerRouter来处理
	eveTy := eve.GetType()
	if eveTy != eventproto.TransactionCtxEventType {
		t.log.Errorf("can not support this event for [%v]", eveTy)
		return nil, fmt.Errorf("can not support this event for [%v]", eveTy)
	}
	// 进行强制类型转换
	if txEventCtx, ok := eve.(*event.TransactionEventContext); ok {
		ctxKey := txEventCtx.GetKey()
		t.recordReceivedEvent(txEventCtx)
		txEvent := txEventCtx.GetEvent()
		opFuncType := txEvent.OpFunc
		crossID, chainID := txEvent.GetCrossID(), txEvent.GetChainID()
		proofResponse, err := t.dispatcher.Invoke(txEvent, conf.TxMsgResultMaxWaitTimeout)
		if err != nil {
			// 记录状态
			if err := t.db.FinishChainCrossState(crossID, chainID, []byte(err.Error()), storetype.StateFailed); err != nil {
				t.log.Errorf("cross[%v]->chain[%v] finish chain cross state failed, ", crossID, chainID, err)
			}
			pResp := &event.ProofResponse{
				ProofResponse: eventproto.ProofResponse{
					CrossId:    txEvent.GetCrossID(),
					Key:        ctxKey,
					Code:       event.FailureResp,
					Msg:        err.Error(),
					OpFunc:     txEvent.OpFunc,
					TxResponse: &eventproto.TxResponse{},
				},
			}
			pResp.SetChainID(txEvent.GetChainID())
			return pResp, err
		}
		// 设置上下文Key
		proofResponse.SetKey(ctxKey)
		// 记录到数据库
		if proofResponse.Code == event.SuccessResp {
			switch opFuncType {
			case event.ExecuteOpFunc:
				if err := t.db.WriteChainCrossState(crossID, chainID, storetype.StateExecuteSuccess, nil); err != nil {
					t.writeStateErrorLog(crossID, chainID, storetype.StateExecuteSuccess, err)
				}
			case event.CommitOpFunc:
				if err := t.db.WriteChainCrossState(crossID, chainID, storetype.StateCommitSuccess, nil); err != nil {
					t.writeStateErrorLog(crossID, chainID, storetype.StateCommitSuccess, err)
				}
			case event.RollbackOpFunc:
				if err := t.db.WriteChainCrossState(crossID, chainID, storetype.StateRollbackSuccess, nil); err != nil {
					t.writeStateErrorLog(crossID, chainID, storetype.StateRollbackSuccess, err)
				}
			}
		} else {
			switch opFuncType {
			case event.ExecuteOpFunc:
				if err := t.db.WriteChainCrossState(crossID, chainID, storetype.StateExecuteFailed, nil); err != nil {
					t.writeStateErrorLog(crossID, chainID, storetype.StateExecuteFailed, err)
				}
			case event.CommitOpFunc:
				if err := t.db.WriteChainCrossState(crossID, chainID, storetype.StateCommitFailed, nil); err != nil {
					t.writeStateErrorLog(crossID, chainID, storetype.StateCommitFailed, err)
				}
			case event.RollbackOpFunc:
				if err := t.db.WriteChainCrossState(crossID, chainID, storetype.StateRollbackFailed, nil); err != nil {
					t.writeStateErrorLog(crossID, chainID, storetype.StateRollbackFailed, err)
				}
			}
		}
		// 等待处理完成
		return proofResponse, err
	} else {
		return nil, errors.New("can not support this event")
	}
}

func (t *TransactionProcessHandler) recordReceivedEvent(eve *event.TransactionEventContext) {
	if err := t.db.WriteChainCrossState(eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), storetype.StateReceived, nil); err != nil {
		t.log.Errorf("cross[%v]->chain[%v] write chain cross state failed, ", eve.GetEvent().GetCrossID(), eve.GetEvent().GetChainID(), err)
	}
}

func (t *TransactionProcessHandler) writeStateErrorLog(crossID, chainID string, state storetype.State, err error) {
	t.log.Errorf("cross[%v]->chain[%v] write chain cross state[%v] failed, ", crossID, chainID, state, err)
}
