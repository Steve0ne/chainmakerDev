/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package adapter

import (
	"fmt"
	"sync"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
	"go.uber.org/zap"
)

var dispatcher *ChainAdapterDispatcher

func init() {
	dispatcher = &ChainAdapterDispatcher{
		adapters: make(map[string]ChainAdapter),
	}
}

// GetChainAdapterDispatcher return instance of ChainAdapterDispatcher which from config file
func GetChainAdapterDispatcher() *ChainAdapterDispatcher {
	return dispatcher
}

// ChainAdapterDispatcher dispatcher of chain adapter which can communication with real chain
type ChainAdapterDispatcher struct {
	sync.RWMutex                         // 读写锁
	adapters     map[string]ChainAdapter // 转接器的Map，按链ID区分
	log          *zap.SugaredLogger      // 日志模块
}

// SetLog set module of logger
func (d *ChainAdapterDispatcher) SetLog(log *zap.SugaredLogger) {
	d.log = log
}

// Register register adapter to adapter dispatcher
func (d *ChainAdapterDispatcher) Register(chainAdapter ChainAdapter) {
	chainID := chainAdapter.GetChainID()
	d.log.Infof("register adapter for chain[%s]", chainID)
	// 无需处理并发
	d.adapters[chainID] = chainAdapter
}

// Invoke transfer transaction event to real adapter
func (d *ChainAdapterDispatcher) Invoke(chainID string, tx *eventproto.TransactionEvent) (*eventproto.TxResponse, error) {
	d.RLock()
	defer d.RUnlock()
	if adapter, exist := d.adapters[chainID]; exist {
		d.log.Infof("find chain[%s]'s adapter", chainID)
		// 判断是否需要进行证明
		if tx.NeedProve() {
			var verifyResult = false
			if verifyResult = adapter.Prove(tx.TxProof); !verifyResult {
				// 证明失败，打印信息，然后返回error
				d.log.Errorf("cross[%s]->chain[%s]'s tx-proof prove failed", tx.GetCrossID(), tx.GetChainID())
				return nil, fmt.Errorf("cross[%s]->chain[%s]'s tx-proof prove failed", tx.GetCrossID(), tx.GetChainID())
			}
			d.log.Infof("cross[%s]->chain[%s]'s tx-proof prove success", tx.GetCrossID(), tx.GetChainID())
			// 将证明及其内容上链
			proofResponse, err := adapter.SaveProof(tx.GetCrossID(), tx.ProofKey, tx.TxProof, verifyResult)
			if err != nil {
				d.log.Errorf("cross[%s]->chain[%s] save proof error, ", tx.GetCrossID(), tx.GetChainID(), err)
				// 保存数据失败，不影响交易主流程，错误不返回
				//return nil, err
			} else {
				// 打印内容，后续写入数据库
				d.log.Infof("save proof success, cross[%s]->chain[%s] txKey[%s] block[%v] index[%v]",
					tx.GetCrossID(), tx.GetChainID(), proofResponse.TxKey, proofResponse.BlockHeight, proofResponse.Index)
			}
		}
		return adapter.Invoke(tx)
	}
	d.log.Errorf("can not find adapter for chain[%s]", chainID)
	return nil, fmt.Errorf("can not find adapter for chain[%v]", chainID)
}

func (d *ChainAdapterDispatcher) SaveProof(chainID, crossID, proofTxKey string, txProof *eventproto.Proof, verifyResult bool) (*eventproto.TxResponse, error) {
	d.RLock()
	defer d.RUnlock()
	if adapter, exist := d.adapters[chainID]; exist {
		d.log.Infof("find chain[%s]'s adapter", chainID)
		proofResponse, err := adapter.SaveProof(crossID, proofTxKey, txProof, verifyResult)
		if err != nil {
			d.log.Errorf("cross[%s]->chain[%s] save proof error, ", crossID, chainID, err)
			// 保存数据失败，此时直接返回错误
			return nil, err
		}
		// 打印内容，后续写入数据库
		d.log.Infof("save proof success, cross[%s]->chain[%s] txKey[%s] block[%v] index[%v]",
			crossID, chainID, proofResponse.TxKey, proofResponse.BlockHeight, proofResponse.Index)
		return proofResponse, err
	}
	d.log.Errorf("can not find adapter for chain[%s]", chainID)
	return nil, fmt.Errorf("can not find adapter for chain[%v]", chainID)
}

// QueryByTxKey query transaction by chain-id and tx-key
func (d *ChainAdapterDispatcher) QueryByTxKey(chainID string, txKey string) (*event.CommonTxResponse, error) {
	d.RLock()
	defer d.RUnlock()
	if adapter, exist := d.adapters[chainID]; exist {
		d.log.Infof("find chain[%s]'s adapter", chainID)
		return adapter.QueryByTxKey(txKey)
	}
	d.log.Errorf("can not find adapter for chain[%s]", chainID)
	return nil, fmt.Errorf("can not find adapter for chain[%v]", chainID)
}

// Query query transaction by chain-id and payload
func (d *ChainAdapterDispatcher) Query(chainID string, payload []byte) (*event.CommonTxResponse, error) {
	d.RLock()
	defer d.RUnlock()
	if adapter, exist := d.adapters[chainID]; exist {
		d.log.Infof("find chain[%s]'s adapter", chainID)
		return adapter.QueryTx(payload)
	}
	d.log.Errorf("can not find adapter for chain[%s]", chainID)
	return nil, fmt.Errorf("can not find adapter for chain[%v]", chainID)
}

// GetChainIDs return all chain-ids which support by all adapters
func (d *ChainAdapterDispatcher) GetChainIDs() []string {
	chainIDs := make([]string, 0)
	for adapterKey := range d.adapters {
		chainID := adapterKey
		chainIDs = append(chainIDs, chainID)
	}
	return chainIDs
}
