/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package kvdb

import (
	"fmt"
	"strings"
	"sync"

	"chainmaker.org/chainmaker-cross/logger"
	kvdbtypes "chainmaker.org/chainmaker-cross/store/kvdb/types"
	storetypes "chainmaker.org/chainmaker-cross/store/types"
	"go.uber.org/zap"
)

const (
	IDSep                  string = ","
	StateBytesIndex        int    = 0
	CrossKeyFormat         string = "C/%s"     // k:C/{CrossID}			v:[]byte		跨链消息
	CrossResultKeyFormat   string = "CR/%s"    // k:CR/{CrossID}			v:[]byte		跨链消息的结果
	CrossChainIDsFormat    string = "L/%s"     // k:L/{CrossID}			v:[]{ChainIDs}	跨链消息涉及的ChainID
	CrossStateFormat       string = "S/%s"     // k:S/{CrossID}			v:int			跨链消息状态
	ChainCrossStateFormat  string = "S/%s/%s"  // k:S/{CrossID}/{ChainID}	v:int			跨链消息某条链的状态
	ChainCrossResultFormat string = "C/%s/%s"  // k:C/{CrossID}/{ChainID}	v:int			跨链消息某条链的执行结果
	UnfinishedCrossSetKey  string = "UF/CROSS" // k:S/{CrossID}			v:int			未完成的跨链交易
)

// KvStateDB is the struct which will be call by other module
type KvStateDB struct {
	sync.Mutex                        // lock
	provider   kvdbtypes.KvDBProvider // DataBase提供者
	logger     *zap.SugaredLogger     // log
}

// NewKvStateDB create new instance of KvStateDB
func NewKvStateDB(provider kvdbtypes.KvDBProvider) *KvStateDB {
	return &KvStateDB{
		provider: provider,
		logger:   logger.GetLogger(logger.ModuleStorage),
	}
}

// StartCross start the cross transaction for crossID
func (k *KvStateDB) StartCross(crossID string, content []byte) error {
	k.Lock()
	defer k.Unlock()
	// 必须加锁
	batch := kvdbtypes.NewKvDBBatcher()
	unfinishedSetKey, unfinishedSetValue := k.redefineAddToUnfinishedSet(crossID)
	batch.Add(crossKey(crossID), content)                                 // 内容
	batch.Add(crossStateKey(crossID), []byte{byte(storetypes.StateInit)}) // 初始化
	batch.Add(unfinishedSetKey, unfinishedSetValue)                       // 添加到未完成集合
	return k.provider.WriteBatch(batch)
}

// FinishCross finish the cross transaction for crossID
func (k *KvStateDB) FinishCross(crossID string, result []byte, state storetypes.State) error {
	k.Lock()
	defer k.Unlock()
	batch := kvdbtypes.NewKvDBBatcher()
	unfinishedSetKey, unfinishedSetValue, exist := k.redefineRemoveFromUnfinishedSet(crossID)
	batch.Add(crossStateKey(crossID), []byte{byte(state)})
	if exist {
		batch.Add(unfinishedSetKey, unfinishedSetValue) // 重置未完成集合
	}
	batch.Add(crossResultKey(crossID), result)
	return k.provider.WriteBatch(batch)
}

// ReadCross read cross content for crossID
func (k *KvStateDB) ReadCross(crossID string) ([]byte, error) {
	// 生成key
	crossKey := crossKey(crossID)
	bz, ok := k.provider.Get(crossKey)
	if !ok {
		return nil, fmt.Errorf("can not find value for key[%s]", crossKey)
	}
	return bz, nil
}

// WriteCross write cross content into database
func (k *KvStateDB) WriteCross(crossID string, content []byte) error {
	// 生成key
	crossKey := crossKey(crossID)
	return k.provider.Put(crossKey, content)
}

// WriteChainIDs write the relationship of crossID and chainIDs into database
func (k *KvStateDB) WriteChainIDs(crossID string, chainIDs []string) error {
	// 字符串拼接
	value := strings.Join(chainIDs, IDSep)
	crossChainsKey := crossChainsKey(crossID)
	return k.provider.Put(crossChainsKey, []byte(value))
}

// ReadChainIDs read chain ids for crossID
func (k *KvStateDB) ReadChainIDs(crossID string) ([]string, bool) {
	crossChainsKey := crossChainsKey(crossID)
	valueBytes, exist := k.provider.Get(crossChainsKey)
	if !exist {
		return nil, exist
	}
	return strings.Split(string(valueBytes), IDSep), true
}

// WriteCrossState write cross state for crossID
func (k *KvStateDB) WriteCrossState(crossID string, state storetypes.State) error {
	stateKey := crossStateKey(crossID)
	return k.provider.Put(stateKey, []byte{byte(state)})
}

// ReadCrossState read cross state for crossID
func (k *KvStateDB) ReadCrossState(crossID string) (storetypes.State, []byte, bool) {
	stateKey := crossStateKey(crossID)
	stateValBytes, exist := k.provider.Get(stateKey)
	if !exist {
		k.logger.Errorf("can not find value for key[%s]", stateKey)
		return storetypes.StateUnknown, nil, false
	}
	stateVal := stateValBytes[StateBytesIndex]
	resultKey := crossResultKey(crossID)
	resultValBytes, exist := k.provider.Get(resultKey)
	if !exist {
		// 表示未找到结果
		k.logger.Warnf("can not find value for key[%s]", resultKey)
		return storetypes.State(stateVal), nil, true
	}
	return storetypes.State(stateVal), resultValBytes, exist
}

// FinishChainCrossState finish cross transaction for crossID and chainID
func (k *KvStateDB) FinishChainCrossState(crossID, chainID string, result []byte, state storetypes.State) error {
	batch := kvdbtypes.NewKvDBBatcher()
	batch.Add(chainCrossStateKey(crossID, chainID), []byte{byte(state)})
	if state == storetypes.StateSuccess {
		// 表明当前链处理成功
		batch.Add(chainCrossResultKey(crossID, chainID), result)
	}
	return k.provider.WriteBatch(batch)
}

// WriteChainCrossState write the state for crossID and chainID
func (k *KvStateDB) WriteChainCrossState(crossID, chainID string, state storetypes.State, content []byte) error {
	batch := kvdbtypes.NewKvDBBatcher()
	batch.Add(chainCrossStateKey(crossID, chainID), []byte{byte(state)})
	if content != nil {
		batch.Add(chainCrossResultKey(crossID, chainID), content)
	}
	return k.provider.WriteBatch(batch)
}

// ReadChainCrossState load state for crossID and chainID
func (k *KvStateDB) ReadChainCrossState(crossID, chainID string) (storetypes.State, []byte, bool) {
	chainCrossStateKey := chainCrossStateKey(crossID, chainID)
	stateValBytes, exist := k.provider.Get(chainCrossStateKey)
	if !exist {
		k.logger.Errorf("can not find value for key[%s]", chainCrossStateKey)
		return storetypes.StateUnknown, nil, false
	}
	chainState := storetypes.State(stateValBytes[StateBytesIndex])
	chainCrossResultKey := chainCrossResultKey(crossID, chainID)
	result, exist := k.provider.Get(chainCrossResultKey)
	return chainState, result, exist
}

// ReadUnfinishedCrossIDs load all the unfinished crossID array
func (k *KvStateDB) ReadUnfinishedCrossIDs() []string {
	k.Lock()
	defer k.Unlock()
	crossIDsBytes, exist := k.provider.Get(UnfinishedCrossSetKey)
	if !exist || crossIDsBytes == nil {
		k.logger.Info("the value is not storage for ", UnfinishedCrossSetKey)
		return nil
	}
	crossIDsString := string(crossIDsBytes)
	return strings.Split(crossIDsString, IDSep)
}

// DeleteCrossIDFromUnfinished delete crossID from UnfinishedCrossIDs
func (k *KvStateDB) DeleteCrossIDFromUnfinished(crossID string) error {
	k.Lock()
	defer k.Unlock()
	batch := kvdbtypes.NewKvDBBatcher()
	unfinishedSetKey, unfinishedSetValue, exist := k.redefineRemoveFromUnfinishedSet(crossID)
	if exist {
		if len(unfinishedSetValue) == 0 {
			unfinishedSetValue = nil
		}
		batch.Add(unfinishedSetKey, unfinishedSetValue) // 重置未完成集合
	}
	return k.provider.WriteBatch(batch)
}

// Close close the database
func (k *KvStateDB) Close() {
	k.provider.Close()
}

// GetLogger return the logger
func (k *KvStateDB) GetLogger() *zap.SugaredLogger {
	return k.logger
}

func (k *KvStateDB) redefineAddToUnfinishedSet(crossID string) (string, []byte) {
	crossIDsBytes, exist := k.provider.Get(UnfinishedCrossSetKey)
	if !exist || crossIDsBytes == nil {
		crossIDArray := []string{crossID}
		value := strings.Join(crossIDArray, IDSep)
		return UnfinishedCrossSetKey, []byte(value)
	}
	crossIDsString := string(crossIDsBytes)
	value := crossIDsString + IDSep + crossID
	return UnfinishedCrossSetKey, []byte(value)
}

func (k *KvStateDB) redefineRemoveFromUnfinishedSet(crossID string) (string, []byte, bool) {
	crossIDsBytes, exist := k.provider.Get(UnfinishedCrossSetKey)
	if !exist {
		return "", nil, false
	}
	crossIDsString := string(crossIDsBytes)
	crossIDs := strings.Split(crossIDsString, IDSep)
	newCrossIDs := make([]string, 0)
	var isExist = false
	// 判断是否在其中
	for _, id := range crossIDs {
		if id != crossID {
			newCrossIDs = append(newCrossIDs, id)
		} else {
			// 表明相等
			isExist = true
		}
	}
	value := strings.Join(newCrossIDs, IDSep)
	return UnfinishedCrossSetKey, []byte(value), isExist
}

func crossKey(crossID string) string {
	return fmt.Sprintf(CrossKeyFormat, crossID)
}

func crossResultKey(crossID string) string {
	return fmt.Sprintf(CrossResultKeyFormat, crossID)
}

func crossChainsKey(crossID string) string {
	return fmt.Sprintf(CrossChainIDsFormat, crossID)
}

func crossStateKey(crossID string) string {
	return fmt.Sprintf(CrossStateFormat, crossID)
}

func chainCrossStateKey(crossID, chainID string) string {
	return fmt.Sprintf(ChainCrossStateFormat, crossID, chainID)
}

func chainCrossResultKey(crossID, chainID string) string {
	return fmt.Sprintf(ChainCrossResultFormat, crossID, chainID)
}
