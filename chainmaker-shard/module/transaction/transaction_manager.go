/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package transaction

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/adapter"
	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/event/coder"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/prover"
	"chainmaker.org/chainmaker-cross/router"
	"chainmaker.org/chainmaker-cross/store"
	storetype "chainmaker.org/chainmaker-cross/store/types"
	"go.uber.org/zap"
)

const (
	RetryCount                     = 5000
	EventChannelLength             = 1024 * 64
	RetryPeriod                    = time.Second * 15
	SupportedChainCount            = 2
	ChainFirstIdx                  = 0
	ChainSecondIdx                 = 1
	CrossChainProofSaveErrorFormat = "save chain[%s] cross[%s] proof to chain error"
	DBCrossStateErrorFormat        = "save cross[%s] state[%v] error"
	DBCrossChainStateErrorFormat   = "save chain[%s] cross[%s] state[%v] error"
)

var (
	dbInitError              = errors.New("transaction manager's db is not init")
	crossEventCoderInitError = errors.New("can not find event coder to handle cross-event")
	crossRespCoderInitError  = errors.New("can not find event coder to handle cross-resp")
	txProofCoderInitError    = errors.New("can not find event coder to handle tx-proof")
	crossChainStateSuccess   = "cross chain success"
)

var transactionManager *Manager

func init() {
	transactionManager = &Manager{
		eventCh:           make(chan event.Event, EventChannelLength),
		routerDispatcher:  router.GetDispatcher(),
		proverDispatcher:  prover.GetProverDispatcher(),
		adapterDispatcher: adapter.GetChainAdapterDispatcher(),
		crossEventCoder:   coder.GetCrossEventCoder(),
		crossRespCoder:    coder.GetCrossRespEventCoder(),
		txProofCoder:      coder.GetTransactionProofCoder(),
	}
}

// Manager is module for transaction
type Manager struct {
	db                store.StateDB                   // 存储
	eventCh           chan event.Event                // 消息事件通道
	routerDispatcher  *router.RouterDispatcher        // 路由管理模块
	proverDispatcher  *prover.ProverDispatcher        // 证明管理模块
	adapterDispatcher *adapter.ChainAdapterDispatcher // 连接器管理模块
	crossEventCoder   event.EventCoder                // 跨链事件编解码器
	crossRespCoder    event.EventCoder                // 跨链返回编解码器
	txProofCoder      event.EventCoder                // 交易证明编解码器
	logger            *zap.SugaredLogger              // log
	cancel            context.CancelFunc              // 退出函数
}

// GetTransactionManager return the instance of transaction manager
func GetTransactionManager() *Manager {
	return transactionManager
}

// GetEventChan return event channel
func (tm *Manager) GetEventChan() chan event.Event {
	return tm.eventCh
}

// SetStateDB set state database
func (tm *Manager) SetStateDB(db store.StateDB) {
	tm.db = db
}

// SetLogger set logger
func (tm *Manager) SetLogger(logger *zap.SugaredLogger) {
	tm.logger = logger
}

// Start start the transaction manager
func (tm *Manager) Start() error {
	// check configs
	if tm.db == nil {
		return dbInitError
	}
	if tm.logger == nil {
		tm.logger = logger.GetLogger(logger.ModuleTransactionMgr)
	}
	if tm.crossEventCoder == nil {
		if eventCoder, exist := coder.GetEventCoderTools().GetDefaultCoder(eventproto.CrossEventType); exist {
			tm.crossEventCoder = eventCoder
		} else {
			return crossEventCoderInitError
		}
	}
	if tm.crossRespCoder == nil {
		if eventCoder, exist := coder.GetEventCoderTools().GetDefaultCoder(eventproto.CrossRespEventType); exist {
			tm.crossRespCoder = eventCoder
		} else {
			return crossRespCoderInitError
		}
	}
	if tm.txProofCoder == nil {
		if eventCoder, exist := coder.GetEventCoderTools().GetDefaultCoder(eventproto.TxProofType); exist {
			tm.txProofCoder = eventCoder
		} else {
			return txProofCoderInitError
		}
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	tm.cancel = cancelFunc
	tm.handleDBCrossEventsStart()
	tm.handleChanCrossEventsStart(ctx)
	return nil
}

// Stop stop the transaction manager
func (tm *Manager) Stop() {
	tm.cancel()
}

// handleDBCrossEventsStart
//- 处理服务重启后，事物模块中缓存的所有跨链事物
func (tm *Manager) handleDBCrossEventsStart() {
	go func() {
		crossEventBytesArray, err := tm.recoverUnfinishedCrossEvents()
		if err != nil {
			tm.logger.Error("Module transaction-manager recover unfinished cross error")
			return
		}
		if len(crossEventBytesArray) > 0 {
			tm.logger.Infof("Module transaction-manager recover unfinished cross: %s", crossEventBytesArray)
		}
		for _, crossEventBytes := range crossEventBytesArray {
			eve, err := tm.crossEventCoder.UnmarshalFromBinary(crossEventBytes)
			if err != nil {
				tm.logger.Error("Module transaction-manager recover unfinished cross error")
				continue
			}
			eventType := eve.GetType()
			if eventType == eventproto.CrossEventType {
				// 跨链事件
				tm.logger.Info("I recover cross event, will handle it!")
				if crossEvent, ok := eve.(*eventproto.CrossEvent); ok {
					tm.handleRecovery(crossEvent)
				} else {
					tm.logger.Warn("This recovery event can not convert to event.CrossEvent")
				}
			} else {
				tm.logger.Errorf("Recovery event type is [%v], so I will not handle it!", eventType)
			}
		}
	}()
}

// handleDBCrossEventsStart
//- 处理跨链事物
func (tm Manager) handleChanCrossEventsStart(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				// 上下文结束，打印日志
				tm.logger.Info("Module transaction-manager stopped")
				return
			case <-time.After(conf.LogWritePeriod):
				// 只是打印日志
				tm.logger.Info("transaction manager is running periodically!")
			case e := <-tm.eventCh:
				eventType := e.GetType()
				if eventType == eventproto.CrossEventType {
					// 跨链事件
					tm.logger.Info("transaction manager receive cross event, start handle it!")
					if crossEvent, ok := e.(*eventproto.CrossEvent); ok {
						if !crossEvent.IsValid() {
							tm.logger.Errorf("This event's crossTxs not match to it's index, so I will not handle it!")
							continue
						}
						tm.handle(crossEvent, false)
					} else {
						tm.logger.Warn("This event can not convert to CrossEvent")
					}
				} else {
					tm.logger.Errorf("transaction manager receive event is [%v], will not handle it!", eventType)
				}
			}
		}
	}(ctx)
}

// handle the cross event
func (tm *Manager) handle(eve *eventproto.CrossEvent, isSync bool) {
	if isSync {
		tm.innerHandle(eve)
	} else {
		go func(eve *eventproto.CrossEvent) {
			tm.innerHandle(eve)
		}(eve)
	}
}

// innerHandle which is inner handle for sync
func (tm *Manager) innerHandle(eve *eventproto.CrossEvent) {
	// 开始该事务处理
	crossID := eve.GetCrossID()
	// 目前只支持两条链的跨链操作
	chainIDs := eve.GetChainIDs()
	if len(chainIDs) != SupportedChainCount {
		tm.logger.Errorf("just support %v chains for cross-event %s", SupportedChainCount, crossID)
		return
	}
	content, err := tm.crossEventCoder.MarshalToBinary(eve)
	if err != nil {
		// 记录错误，无需处理其他
		_ = tm.db.FinishCross(crossID, []byte(err.Error()), storetype.StateFailed)
		return
	}
	if err = tm.db.StartCross(crossID, content); err != nil {
		tm.logger.Error("save event cross start error", err)
		return
	}
	txEvents := eve.GetPkgTxEvents()
	// sort by index
	sort.Sort(txEvents)
	tm.handleTxEvents(crossID, txEvents.GetCrossTxs())
}

// Handle handle the cross event which is unfinished last time
func (tm *Manager) handleRecovery(eve *eventproto.CrossEvent) {
	go func(eve *eventproto.CrossEvent) {
		// 开始该跨链事务处理
		crossID := eve.GetCrossID()
		chainIDs := eve.GetChainIDs()
		if len(chainIDs) != SupportedChainCount {
			tm.logger.Errorf("just support %v chains for cross-event %s", SupportedChainCount, crossID)
			return
		}
		txEvents := eve.GetPkgTxEvents()
		sort.Sort(txEvents)
		crossTxs := txEvents.Events //GetCrossTxs()
		firstEventTx, secondEventTx := crossTxs[ChainFirstIdx], crossTxs[ChainSecondIdx]
		// 检查第一笔交易的状态
		firstTxState, result, exist := tm.db.ReadChainCrossState(crossID, firstEventTx.GetChainID())
		if !exist {
			// 第一笔交易不存在，存在两种情况：
			// 1、该交易尚未发送就宕机，则需要两笔都重新发一下
			// 2、该交易已经发送，则需要判断其状态，然后再判断第二笔是否需要发送
			txResponse, err := tm.adapterDispatcher.Query(firstEventTx.GetChainID(), firstEventTx.ExecutePayload)
			if err != nil || txResponse == nil {
				// 表示出现错误、或没有应答，则重新执行
				tm.handleTxEvents(crossID, crossTxs)
			} else {
				// 判断当前交易的状态
				if txResponse.IsSuccess() {
					// 执行下一笔交易
					// 生成第一笔交易的证明
					proof := event.NewProof(txResponse.GetChainID(), txResponse.TxKey, txResponse.BlockHeight,
						txResponse.Index, txResponse.Contract, txResponse.Extra)
					if proofBytes, err := tm.txProofCoder.MarshalToBinary(proof); err != nil {
						tm.logger.Errorf("cross[%v]->chain[%v]'s tx-proof marshal error", crossID, proof.GetChainID(), err)
					} else {
						if err := tm.db.WriteChainCrossState(crossID, proof.GetChainID(), storetype.StateExecuteSuccess, proofBytes); err != nil {
							tm.logger.Errorf(DBCrossChainStateErrorFormat, proof.GetChainID(), crossID, storetype.StateExecuteSuccess)
						}
					}
					// 重新执行第二笔交易
					tm.retrySecondTx(crossID, firstEventTx, secondEventTx, proof)
				}
			}
			return
		} else {
			// 获取第二笔交易的状态
			secondTxState, secondResult, secondStateExist := tm.db.ReadChainCrossState(crossID, secondEventTx.GetChainID())
			if firstTxState == storetype.StateExecuteSuccess {
				if !secondStateExist {
					// 第二笔交易没有状态，那么需要重新进行提交
					// 对证明进行转换
					proofEvent, err := tm.txProofCoder.UnmarshalFromBinary(result)
					if err != nil {
						tm.logger.Errorf("unmarshal cross[%s]->chain[%s] failed", crossID, secondEventTx.GetChainID())
						// 处理本地回退
						if ok := tm.rollbackCrossTx(crossID, firstEventTx); ok {
							// 删除 unfinished
							tm.finishCrossEvent(crossID)
						}
						return
					}
					if proof, ok := proofEvent.(*eventproto.Proof); ok {
						// 重新执行第二笔交易
						tm.retrySecondTx(crossID, firstEventTx, secondEventTx, proof)
						return
					} else {
						// 打印日志，本地节点回滚
						tm.logger.Errorf("cross[%v]->chain[%v]'s proof can not be convert", crossID, firstEventTx.GetChainID())
						// 记录该链处理错误
						if err := tm.db.WriteChainCrossState(crossID, firstEventTx.GetChainID(), storetype.StateFailed, nil); err != nil {
							tm.logger.Errorf(DBCrossChainStateErrorFormat, firstEventTx.GetChainID(), crossID, storetype.StateFailed)
						}
						// 此时有错误，回滚第一笔交易
						tm.rollbackCrossTx(crossID, firstEventTx)
						// 记录整体状态，结束该事务
						tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("execute chain[%v] error", firstEventTx.GetChainID())))
						return
					}
				} else {
					// 判断其状态，若提交成功，则只需要提交第一笔交易
					if secondTxState == storetype.StateCommitSuccess {
						// 另外一笔交易已经提交成功，当前交易提交
						tm.commitDesignativeTx(crossID, firstEventTx)
						return
					}
					// 若操作失败，则需要重新提交两笔
					if secondTxState == storetype.StateCommitFailed {
						tm.commitTwoTxs(crossID, firstEventTx, secondEventTx)
						return
					}
					// 若操作失败，则需要回退两笔交易
					if secondTxState == storetype.StateProofConvertFailed || secondTxState == storetype.StateProofFailed {
						// 记录证明到本地
						secondProofEvent, err := tm.txProofCoder.UnmarshalFromBinary(secondResult)
						if err != nil {
							tm.logger.Errorf("cross[%v]->chain[%v]'s proof can not be convert", crossID, secondEventTx.GetChainID())
						} else {
							// 保存证明信息，写到链上
							if proof, ok := secondProofEvent.(*eventproto.Proof); ok {
								err = tm.saveProof(firstEventTx.GetChainID(), crossID, firstEventTx.ProofKey, proof, false)
								if err != nil {
									tm.logger.Errorf("cross[%v]->chain[%v]'s save chain[%s]'s proof error, ", crossID, firstEventTx.GetChainID(), secondEventTx.GetChainID(), err)
								}
							}
						}
						// 两个都进行回退操作
						tm.rollbackTwoTxs(crossID, firstEventTx, secondEventTx)
						return
					}
					if secondTxState == storetype.StateProofSuccess {
						// 记录证明到本地
						secondProofEvent, err := tm.txProofCoder.UnmarshalFromBinary(secondResult)
						if err != nil {
							tm.logger.Errorf("cross[%v]->chain[%v]'s proof can not be convert", crossID, secondEventTx.GetChainID())
						} else {
							// 保存证明信息，写到链上
							if proof, ok := secondProofEvent.(*eventproto.Proof); ok {
								err = tm.saveProof(firstEventTx.GetChainID(), crossID, firstEventTx.ProofKey, proof, true)
								if err != nil {
									tm.logger.Errorf("cross[%v]->chain[%v]'s save chain[%s]'s proof error, ", crossID, firstEventTx.GetChainID(), secondEventTx.GetChainID(), err)
								}
							}
						}
						// 另外一笔已经成功执行，需要提交两笔交易
						tm.commitTwoTxs(crossID, firstEventTx, secondEventTx)
						return
					}
				}
			} else if firstTxState == storetype.StateCommitSuccess {
				// 提交成功，则判断第二笔是否提交成功
				if secondTxState == storetype.StateCommitSuccess {
					// 两笔都提交成功的话，只需要删除即可
					tm.finishCrossEvent(crossID)
					return
				} else {
					// 提交第二笔
					tm.commitDesignativeTx(crossID, secondEventTx)
					return
				}
			} else if firstTxState == storetype.StateCommitFailed {
				// 提交失败的话，判断第二笔是否成功
				if secondTxState == storetype.StateCommitSuccess {
					// 只需要提交第一笔即可
					tm.commitDesignativeTx(crossID, firstEventTx)
					return
				} else {
					tm.commitTwoTxs(crossID, firstEventTx, secondEventTx)
					return
				}
			} else if firstTxState == storetype.StateRollbackSuccess {
				// 判断第二笔是否回滚成功
				if secondTxState == storetype.StateRollbackSuccess {
					tm.finishCrossEvent(crossID)
					return
				} else {
					tm.rollbackDesignativeTx(crossID, secondEventTx)
					return
				}
			} else if firstTxState == storetype.StateRollbackFailed {
				if secondTxState == storetype.StateRollbackSuccess {
					tm.rollbackDesignativeTx(crossID, firstEventTx)
					return
				} else {
					tm.rollbackTwoTxs(crossID, firstEventTx, secondEventTx)
					return
				}
			}
		}
	}(eve)
}

// handleTxEvents
func (tm *Manager) handleTxEvents(crossID string, crossTxs []*eventproto.CrossTx) {
	handledPkgTxEvents := make([]*eventproto.CrossTx, 0)
	allResponse := make([]*event.ProofResponse, 0)
	var majorProof *eventproto.Proof
	// 仅支持两条链操作，A-B，A的结果由B来验证，B的结果由A来验证
	for _, pkgTxEve := range crossTxs {
		pkgTxEvent := pkgTxEve
		chainID := pkgTxEvent.GetChainID()
		resp, err := tm.execute(crossID, pkgTxEvent, majorProof)
		if err != nil {
			tm.logger.Errorf("cross[%v]->chain[%v]'s execute payload error, ", crossID, chainID, err)
			// 记录该链处理错误
			if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateFailed, nil); err != nil {
				tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateFailed)
			}
			// 此时有错误，需要回滚之前已经完成的提交
			tm.rollbackHandledEvents(crossID, handledPkgTxEvents, err)
			// 记录整体状态，结束该事务
			tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("execute chain[%v] error", chainID)))
			return
		} else {
			if resp.IsSuccess() {
				tm.logger.Infof("cross[%v]->chain[%v]'s execute payload success", crossID, chainID)
				allResponse = append(allResponse, resp)
				// 成功，则进行下一个处理
				handledPkgTxEvents = append(handledPkgTxEvents, pkgTxEvent)
				if majorProof == nil {
					// 表示为第一条链的操作
					majorProof, err = tm.toProof(chainID, resp)
					tm.logger.Infof("cross[%v]->chain[%v]'s response convert to proof success", crossID, chainID)
					if err != nil {
						tm.logger.Errorf("cross[%v]->chain[%v]'s response convert to proof error", crossID, chainID, err)
						if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateProofConvertFailed, nil); err != nil {
							tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateProofConvertFailed)
						}
						// 表示无法转换，需要进行回滚
						tm.rollbackHandledEvents(crossID, handledPkgTxEvents, err)
						tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("convert chain[%v]'s response to proof error", chainID)))
						return
					}
					// 表示交易执行成功，并且证明OK
					if proofBytes, err := tm.txProofCoder.MarshalToBinary(majorProof); err != nil {
						tm.logger.Errorf("cross[%v]->chain[%v]'s tx-proof marshal error", crossID, chainID, err)
					} else {
						if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateExecuteSuccess, proofBytes); err != nil {
							tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateExecuteSuccess)
						}
					}
				} else {
					// 表示非第一个操作，那么需要对其进行证明
					// 表示为第二条链上的操作
					var (
						remoteProof      *eventproto.Proof
						remoteProofBytes []byte
					)
					remoteProof, err = tm.toProof(chainID, resp)
					tm.logger.Infof("cross[%v]->chain[%v]'s response convert to proof success", crossID, chainID)
					if err != nil {
						tm.logger.Errorf("cross[%v]->chain[%v]'s response convert to proof error", crossID, chainID, err)
						if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateProofConvertFailed, nil); err != nil {
							tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateProofConvertFailed)
						}
						// 表示无法转换，需要进行回滚
						tm.rollbackHandledEvents(crossID, handledPkgTxEvents, err)
						tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("convert chain[%v]'s response to proof error", chainID)))
						return
					}
					remoteProofBytes, err = tm.txProofCoder.MarshalToBinary(remoteProof)
					// TODO 后续采用更优雅的方式
					firstCrossTx := crossTxs[0]
					if ok, _ := tm.proverDispatcher.Prove(remoteProof); !ok {
						// 证明失败
						tm.logger.Errorf("cross[%v]->chain[%v]'s proof check error", crossID, chainID)
						if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateProofFailed, remoteProofBytes); err != nil {
							tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateProofFailed)
						}
						// 证明完成后需要将该证据写入到链上
						if err := tm.saveProof(firstCrossTx.GetChainID(), crossID, firstCrossTx.ProofKey, remoteProof, ok); err != nil {
							tm.logger.Errorf(CrossChainProofSaveErrorFormat, chainID, crossID)
						}
						// 证明失败，则需要回滚
						tm.rollbackHandledEvents(crossID, handledPkgTxEvents, fmt.Errorf("can not prove chain[%v]'s proof", chainID))
						tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("can not prove chain[%v]'s proof", chainID)))
						return
					} else {
						// 证明成功
						tm.logger.Infof("cross[%v]->chain[%v]'s proof check success", crossID, chainID)
						if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateProofSuccess, remoteProofBytes); err != nil {
							tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateProofSuccess)
						}
						// 证明完成后需要将该证据写入到链上
						if err := tm.saveProof(firstCrossTx.GetChainID(), crossID, firstCrossTx.ProofKey, remoteProof, ok); err != nil {
							tm.logger.Errorf(CrossChainProofSaveErrorFormat, chainID, crossID)
						}
					}
				}
			} else {
				tm.logger.Errorf("cross[%v]->chain[%v]'s execute failed, %s", crossID, chainID, resp.Msg)
				if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateFailed, nil); err != nil {
					tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateFailed)
				}
				// 重新生成需要回滚的交易对象列表
				rollbackEvents := make([]*eventproto.CrossTx, 0)
				rollbackEvents = append(rollbackEvents, handledPkgTxEvents...)
				// 也需要回滚当前的交易，当前交易是否回滚由事务合约控制
				rollbackEvents = append(rollbackEvents, pkgTxEvent)
				// 失败的情况下需要回滚之前已完成的提交
				tm.rollbackHandledEvents(crossID, rollbackEvents, fmt.Errorf("chain[%v]'s execute failed for %s", chainID, resp.Msg))
				// 记录整体状态，结束该事务
				tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("chain[%v]'s execute failed for %s", chainID, resp.Msg)))
				return
			}
		}
	}
	// 执行到此处表示prepare阶段已全部处理完成，下面需要并发进行commit操作
	tm.commitAll(crossID, handledPkgTxEvents, allResponse)
}

// retrySecondTx
func (tm *Manager) retrySecondTx(crossID string, firstEventTx, secondEventTx *eventproto.CrossTx, proof *eventproto.Proof) {
	// 重新执行第二笔交易, sync execute
	proofResponse, err := tm.execute(crossID, secondEventTx, proof)
	if err != nil {
		tm.logger.Errorf("cross[%v]->chain[%v]'s execute payload error, ", crossID, secondEventTx.GetChainID(), err)
		// 记录该链处理错误
		if err := tm.db.WriteChainCrossState(crossID, secondEventTx.GetChainID(), storetype.StateFailed, nil); err != nil {
			tm.logger.Errorf(DBCrossChainStateErrorFormat, secondEventTx.GetChainID(), crossID, storetype.StateFailed)
		}
		// 此时有错误，回滚第一笔交易
		tm.rollbackCrossTx(crossID, firstEventTx)
		// 记录整体状态，结束该事务
		tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("execute chain[%v] error", secondEventTx.GetChainID())))
		return
	} else {
		// 判断结果
		chainID := secondEventTx.GetChainID()
		if proofResponse.IsSuccess() {
			var (
				remoteProof      *eventproto.Proof
				remoteProofBytes []byte
			)
			remoteProof, err = tm.toProof(chainID, proofResponse)
			tm.logger.Infof("cross[%v]->chain[%v]'s response convert to proof success", crossID, chainID)
			if err != nil {
				tm.logger.Errorf("cross[%v]->chain[%v]'s response convert to proof error", crossID, chainID, err)
				if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateProofConvertFailed, nil); err != nil {
					tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateProofConvertFailed)
				}
				// 表示无法转换，需要进行回滚
				tm.rollbackHandledEvents(crossID, []*eventproto.CrossTx{firstEventTx, secondEventTx}, err)
				tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("convert chain[%v]'s response to proof error", chainID)))
				return
			}
			remoteProofBytes, err = tm.txProofCoder.MarshalToBinary(remoteProof)
			if ok, _ := tm.proverDispatcher.Prove(remoteProof); !ok {
				tm.logger.Errorf("cross[%v]->chain[%v]'s proof check error", crossID, chainID)
				if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateProofFailed, remoteProofBytes); err != nil {
					tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateProofFailed)
				}
				// 证明完成后需要将该证据写入到链上，此时是写到当前链（即链1）
				if err := tm.saveProof(firstEventTx.GetChainID(), crossID, firstEventTx.ProofKey, remoteProof, ok); err != nil {
					tm.logger.Errorf(CrossChainProofSaveErrorFormat, chainID, crossID)
				}
				// 证明异常，则需要回滚
				tm.rollbackHandledEvents(crossID, []*eventproto.CrossTx{firstEventTx, secondEventTx}, fmt.Errorf("can not prove chain[%v]'s proof", chainID))
				tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("can not prove chain[%v]'s proof", chainID)))
				return
			} else {
				// 证明成功
				tm.logger.Infof("cross[%v]->chain[%v]'s proof check success", crossID, chainID)
				if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateProofSuccess, remoteProofBytes); err != nil {
					tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateProofSuccess)
				}
				// 证明完成后需要将该证据写入到链上，此时是写到当前链（即链1）
				if err := tm.saveProof(firstEventTx.GetChainID(), crossID, firstEventTx.ProofKey, remoteProof, ok); err != nil {
					tm.logger.Errorf(CrossChainProofSaveErrorFormat, chainID, crossID)
				}
				// 提交两笔交易
				firstProofResponse := event.NewProofResponseByProof(crossID, chainID, crossChainStateSuccess, event.SuccessResp, event.ExecuteOpFunc, proof)
				tm.commitAll(crossID, []*eventproto.CrossTx{firstEventTx, secondEventTx}, []*event.ProofResponse{firstProofResponse, proofResponse})
				return
			}
		} else {
			// 应答失败，则回滚本地结果
			tm.logger.Errorf("cross[%v]->chain[%v]'s execute failed, %s", crossID, chainID, proofResponse.Msg)
			if err := tm.db.WriteChainCrossState(crossID, chainID, storetype.StateFailed, nil); err != nil {
				tm.logger.Errorf(DBCrossChainStateErrorFormat, chainID, crossID, storetype.StateFailed)
			}
			// 失败的情况下需要回滚之前已完成的提交
			tm.rollbackCrossTx(crossID, firstEventTx)
			// 记录整体状态，结束该事务
			tm.recordInterruptedState(crossID, []byte(fmt.Sprintf("chain[%v]'s execute failed for %s", chainID, proofResponse.Msg)))
			return
		}
	}
}

// commitTwoTxs
func (tm *Manager) commitTwoTxs(crossID string, firstEventTx, secondEventTx *eventproto.CrossTx) {
	crossEventTxs := []*eventproto.CrossTx{firstEventTx, secondEventTx}
	var (
		wg           sync.WaitGroup
		successCount int32 = 0
	)
	wg.Add(len(crossEventTxs))
	for _, eventTx := range crossEventTxs {
		commitSuccess := tm.commitCrossTx(crossID, eventTx)
		wg.Done()
		chainID := eventTx.GetChainID()
		if commitSuccess {
			tm.logger.Infof("cross[%v]->chain[%v] commit success", crossID, chainID)
			tm.recordChainState(crossID, chainID, storetype.StateCommitSuccess)
			atomic.AddInt32(&successCount, 1)
		} else {
			tm.logger.Infof("cross[%v]->chain[%v] commit failed", crossID, chainID)
			tm.recordChainState(crossID, chainID, storetype.StateCommitFailed)
		}
	}
	wg.Wait()
	if int(successCount) >= len(crossEventTxs) {
		tm.finishCrossEvent(crossID)
	}
}

// commitDesignativeTx
func (tm *Manager) commitDesignativeTx(crossID string, eventTx *eventproto.CrossTx) {
	commitSuccess := tm.commitCrossTx(crossID, eventTx)
	chainID := eventTx.GetChainID()
	if commitSuccess {
		tm.logger.Infof("cross[%v]->chain[%v] commit success", crossID, chainID)
		tm.recordChainState(crossID, chainID, storetype.StateCommitSuccess)
		tm.finishCrossEvent(crossID)
	} else {
		tm.logger.Infof("cross[%v]->chain[%v] commit failed", crossID, chainID)
		tm.recordChainState(crossID, chainID, storetype.StateCommitFailed)
	}
}

// rollbackTwoTxs
func (tm *Manager) rollbackTwoTxs(crossID string, firstEventTx, secondEventTx *eventproto.CrossTx) {
	crossEventTxs := []*eventproto.CrossTx{firstEventTx, secondEventTx}
	var (
		wg           sync.WaitGroup
		successCount int32 = 0
	)
	wg.Add(len(crossEventTxs))
	for _, crossTx := range crossEventTxs {
		go func(crossTx *eventproto.CrossTx) {
			rollbackSuccess := tm.rollbackCrossTx(crossID, crossTx)
			wg.Done()
			chainID := crossTx.GetChainID()
			if rollbackSuccess {
				tm.logger.Infof("cross[%v]->chain[%v] rollback success", crossID, chainID)
				tm.recordChainState(crossID, chainID, storetype.StateRollbackSuccess)
				atomic.AddInt32(&successCount, 1)
			} else {
				tm.logger.Infof("cross[%v]->chain[%v] rollback failed", crossID, chainID)
				tm.recordChainState(crossID, chainID, storetype.StateRollbackFailed)
			}
		}(crossTx)
	}
	wg.Wait()
	if int(successCount) >= len(crossEventTxs) {
		tm.finishCrossEvent(crossID)
	}
}

func (tm *Manager) rollbackDesignativeTx(crossID string, eventTx *eventproto.CrossTx) {
	rollbackSuccess := tm.rollbackCrossTx(crossID, eventTx)
	chainID := eventTx.GetChainID()
	if rollbackSuccess {
		tm.logger.Infof("cross[%v]->chain[%v] rollback success", crossID, chainID)
		tm.recordChainState(crossID, chainID, storetype.StateRollbackSuccess)
		tm.finishCrossEvent(crossID)
	} else {
		tm.logger.Infof("cross[%v]->chain[%v] rollback failed", crossID, chainID)
		tm.recordChainState(crossID, chainID, storetype.StateRollbackFailed)
	}
}

func (tm *Manager) execute(crossID string, crossTx *eventproto.CrossTx, proof *eventproto.Proof) (*event.ProofResponse, error) {
	chainID := crossTx.GetChainID()
	tm.logger.Infof("cross[%v]->chain[%v]'s execute start", crossID, chainID)
	// 创建交易
	eve := event.NewExecuteTransactionEvent(crossID, chainID, crossTx.GetExecutePayload(), crossTx.ProofKey, proof)
	return tm.routerDispatcher.Invoke(eve, conf.TxMsgResultMaxWaitTimeout)
}

// commit
// 提交，暂时不需要进行证明
func (tm *Manager) commit(crossID string, pkgTxEvent *eventproto.CrossTx) (*event.ProofResponse, error) {
	return tm.secondPhaseHandle(crossID, pkgTxEvent, event.CommitOpFunc)
}

// rollback
// 回滚，暂时不需要进行结果证明
func (tm *Manager) rollback(crossID string, pkgTxEvent *eventproto.CrossTx) (*event.ProofResponse, error) {
	return tm.secondPhaseHandle(crossID, pkgTxEvent, event.RollbackOpFunc)
}

func (tm *Manager) secondPhaseHandle(crossID string, crossTx *eventproto.CrossTx, opFunc eventproto.OpFuncType) (*event.ProofResponse, error) {
	chainID := crossTx.GetChainID()
	var eve *eventproto.TransactionEvent
	// 创建交易
	if opFunc == event.CommitOpFunc {
		eve = event.NewCommitTransactionEvent(crossID, chainID, crossTx.GetCommitPayload())
	} else if opFunc == event.RollbackOpFunc {
		eve = event.NewRollbackTransactionEvent(crossID, chainID, crossTx.GetRollbackPayload())
	} else {
		return nil, fmt.Errorf("can not support operate func [%v]", opFunc)
	}
	return tm.routerDispatcher.Invoke(eve, conf.TxMsgResultMaxWaitTimeout)
}

// saveProof save proof to chain
func (tm *Manager) saveProof(chainID, crossID, proofTxKey string, proof *eventproto.Proof, verifiedResult bool) error {
	_, err := tm.adapterDispatcher.SaveProof(chainID, crossID, proofTxKey, proof, verifiedResult)
	return err
}

func (tm *Manager) rollbackHandledEvents(crossID string, handledPkgTxEvents []*eventproto.CrossTx, err error) {
	handledEventSize := len(handledPkgTxEvents)
	if handledEventSize > 0 {
		tm.logger.Infof("cross[%v] there are %v event need rollback", crossID, handledEventSize)
		var wg sync.WaitGroup
		wg.Add(handledEventSize)
		var rollbackSuccessSize int32 = 0
		// 并发回滚即可
		for _, handledPkgTxEvent := range handledPkgTxEvents {
			go func(txEve *eventproto.CrossTx) {
				chainID := txEve.GetChainID()
				tm.logger.Infof("cross[%v]->chain[%v] will rollback", crossID, chainID)
				rollbackSuccess := tm.rollbackCrossTx(crossID, txEve)
				wg.Done()
				// 进行状态记录
				if rollbackSuccess {
					tm.logger.Infof("cross[%v]->chain[%v] rollback success", crossID, chainID)
					tm.recordChainState(crossID, chainID, storetype.StateRollbackSuccess)
					atomic.AddInt32(&rollbackSuccessSize, 1)
				} else {
					tm.logger.Warnf("cross[%v]->chain[%v] rollback failed", crossID, chainID)
					tm.recordChainState(crossID, chainID, storetype.StateRollbackFailed)
				}
			}(handledPkgTxEvent)
		}
		wg.Wait()
		// 判断是否全部回滚成功
		if int(rollbackSuccessSize) >= handledEventSize {
			tm.logger.Infof("cross[%v] rollback completed", crossID)
			// 记录失败
			tm.recordFailedFinishedState(crossID, err)
		}
	} else {
		tm.logger.Warn("there are non events will be rollback")
	}
}

func (tm *Manager) rollbackCrossTx(crossID string, txEve *eventproto.CrossTx) bool {
	var rollbackSuccess = false
	chainID := txEve.GetChainID()
	re, err := tm.rollback(crossID, txEve)
	// 异常或操作失败均需要重试
	if err != nil || !re.IsSuccess() {
		if err != nil {
			tm.logger.Warnf("cross[%v]->chain[%v] rollback failed, ", crossID, chainID, err)
		}
		if !re.IsSuccess() {
			tm.logger.Warnf("cross[%v]->chain[%v] rollback failed -> %s", crossID, chainID, re.Msg)
		}
		// 进行重试操作
		for i := 0; i < RetryCount; i++ {
			time.Sleep(RetryPeriod) // 先进行休眠
			re, err = tm.rollback(crossID, txEve)
			if err == nil && re.IsSuccess() {
				// 操作成功，状态记录
				rollbackSuccess = true
				break
			} else if err != nil {
				tm.logger.Warnf("cross[%v]->chain[%v]->[%v] rollback failed, ", crossID, chainID, i+1, err)
			} else if !re.IsSuccess() {
				tm.logger.Warnf("cross[%v]->chain[%v]->[%v] rollback failed -> %s", crossID, chainID, i+1, re.Msg)
			}
		}
	} else {
		// 操作成功，状态更新
		rollbackSuccess = true
	}
	return rollbackSuccess
}

func (tm *Manager) commitAll(crossID string, handledPkgTxEvents []*eventproto.CrossTx, allResponse []*event.ProofResponse) {
	handledEventSize := len(handledPkgTxEvents)
	if handledEventSize > 0 {
		tm.logger.Infof("cross[%v] there are %v event need commit", crossID, handledEventSize)
		var wg sync.WaitGroup
		wg.Add(handledEventSize)
		// 并发commit
		// 并发回滚即可
		for _, handledPkgTxEvent := range handledPkgTxEvents {
			go func(txEve *eventproto.CrossTx) {
				chainID := txEve.GetChainID()
				tm.logger.Infof("cross[%v]->chain[%v] will commit", crossID, chainID)
				commitSuccess := tm.commitCrossTx(crossID, txEve)
				wg.Done()
				if commitSuccess {
					tm.logger.Infof("cross[%v]->chain[%v] commit success", crossID, chainID)
					tm.recordChainState(crossID, chainID, storetype.StateCommitSuccess)
				} else {
					tm.logger.Infof("cross[%v]->chain[%v] commit failed", crossID, chainID)
					tm.recordChainState(crossID, chainID, storetype.StateCommitFailed)
				}
			}(handledPkgTxEvent)
		}
		wg.Wait()
		tm.recordSuccessFinishedState(crossID, allResponse)
	} else {
		tm.logger.Warn("there are non events will be commit")
	}
}

func (tm *Manager) commitCrossTx(crossID string, txEve *eventproto.CrossTx) bool {
	re, err := tm.commit(crossID, txEve)
	var commitSuccess = false
	// 异常或操作失败均需要重试
	if err != nil || !re.IsSuccess() {
		// 进行重试操作
		for i := 0; i < RetryCount; i++ {
			time.Sleep(RetryPeriod) // 先进行休眠
			re, err = tm.commit(crossID, txEve)
			if err == nil && re.IsSuccess() {
				// 操作成功，状态记录
				commitSuccess = true
				break
			}
		}
	} else {
		// 操作成功，状态更新
		commitSuccess = true
	}
	return commitSuccess
}

func (tm *Manager) toProof(chainID string, response *event.ProofResponse) (*eventproto.Proof, error) {
	return tm.proverDispatcher.ToProof(chainID, response.GetTxKey(), response.GetBlockHeight(), response.GetIndex(),
		response.GetContract(), response.GetExtra())
}

func (tm *Manager) recordInterruptedState(crossID string, result []byte) {
	if err := tm.db.FinishCross(crossID, result, storetype.StateFailed); err != nil {
		tm.logger.Errorf(DBCrossStateErrorFormat, crossID, storetype.StateFailed)
	}
}

func (tm *Manager) recordSuccessFinishedState(crossID string, allResponse []*event.ProofResponse) {
	var (
		code  int32 = event.SuccessResp
		state       = storetype.StateSuccess
	)
	// 生成CrossResponse
	crossResponse := event.NewCrossResponse(crossID, code, crossChainStateSuccess)
	for _, resp := range allResponse {
		r := resp
		tm.logger.Infof("cross[%s]'s result = [%v]", crossID, r)
		crossResponse.AddTxResponse(event.NewCrossTxResponse(r.GetChainID(), r.GetTxKey(), r.GetBlockHeight(), r.GetIndex(), r.GetExtra()))
	}
	// 对crossResponse进行序列化操作
	binary, err := tm.crossRespCoder.MarshalToBinary(crossResponse)
	if err != nil {
		tm.logger.Info("marshal cross response failed,", err)
	} else {
		if err := tm.db.FinishCross(crossID, binary, state); err != nil {
			tm.logger.Errorf(DBCrossStateErrorFormat, crossID, state)
		}
	}
}

func (tm *Manager) recordFailedFinishedState(crossID string, err error) {
	var (
		code  int32 = event.FailureResp
		state       = storetype.StateFailed
	)
	// 生成CrossResponse
	crossResponse := event.NewCrossResponse(crossID, code, err.Error())
	// 对crossResponse进行序列化操作
	binary, err := tm.crossRespCoder.MarshalToBinary(crossResponse)
	if err != nil {
		tm.logger.Info("marshal cross response failed,", err)
	} else {
		if err := tm.db.FinishCross(crossID, binary, state); err != nil {
			tm.logger.Errorf(DBCrossStateErrorFormat, crossID, state)
		}
	}
}

func (tm *Manager) recordChainState(crossID, chainID string, state storetype.State) {
	if err := tm.db.WriteChainCrossState(crossID, chainID, state, nil); err != nil {
		tm.logger.Errorf(DBCrossChainStateErrorFormat, crossID, state)
	}
}

func (tm *Manager) finishCrossEvent(crossID string) {
	// 只是从数据库中删除该值
	if err := tm.db.DeleteCrossIDFromUnfinished(crossID); err != nil {
		tm.logger.Errorf("delete crossID[%s] from db failed, ", crossID, err)
	}
}

// RecoverUnfinishedCross - 从数据库中还原未完成的跨链交易
func (tm *Manager) recoverUnfinishedCrossEvents() ([][]byte, error) {
	var (
		contents [][]byte
		content  []byte
		err      error
	)
	crossIDs := tm.db.ReadUnfinishedCrossIDs()
	if len(crossIDs) == 0 {
		tm.logger.Info("No cross event recovery")
		return nil, nil
	}

	for _, id := range crossIDs {
		if id == "" {
			continue
		}
		content, err = tm.db.ReadCross(id)
		if err != nil {
			tm.logger.Errorf("Read CrossID: [%s] error: ", id, err)
			continue
		}
		contents = append(contents, content)
	}
	if len(contents) == 0 {
		tm.logger.Info("No cross event recovery")
		return nil, nil
	}
	return contents, nil
}
