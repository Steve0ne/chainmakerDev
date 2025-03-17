/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package kvdb

import (
	"bytes"
	"strconv"
	"testing"
	"time"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/store/kvdb/factory"
	storetypes "chainmaker.org/chainmaker-cross/store/types"
)

func TestKvStateDB_StartAndFinishCross(t *testing.T) {
	stateDB := newKvStateDB(t)
	defer stateDB.Close()
	crossID := strconv.Itoa(time.Now().Nanosecond())
	// 创建 对应content
	content := []byte("this is cross event")
	// 开始
	err := stateDB.StartCross(crossID, content)
	if err != nil {
		t.Errorf("start cross %s error: %s", crossID, err.Error())
	}
	// 查询状态
	crossState, _, exist := stateDB.ReadCrossState(crossID)
	if !exist {
		t.Errorf("can not find state for %s", crossID)
	}
	if crossState != storetypes.StateInit {
		t.Errorf("db state is error for %s", crossID)
	}
	// 将状态调整为完成
	var result = []byte("it is success")
	err = stateDB.FinishCross(crossID, result, storetypes.StateSuccess)
	if err != nil {
		t.Errorf("finish cross %s error: %s", crossID, err.Error())
	}
	// 重新查询该跨链状态
	crossState, dbResult, exist := stateDB.ReadCrossState(crossID)
	if !exist {
		t.Errorf("can not find state for %s", crossID)
	}
	if !bytes.Equal(dbResult, result) {
		t.Error("read state is not equal")
	}
	if crossState != storetypes.StateSuccess {
		t.Errorf("db state is error for %s", crossID)
	}
}

func TestKvStateDB_ChainIDs(t *testing.T) {
	stateDB := newKvStateDB(t)
	defer stateDB.Close()
	crossID := strconv.Itoa(time.Now().Nanosecond())
	// 创建 对应content
	content := []byte("this is cross event")
	// 开始
	err := stateDB.StartCross(crossID, content)
	if err != nil {
		t.Errorf("start cross %s error: %s", crossID, err.Error())
	}
	chainIDs := []string{
		"chain1",
		"chain2",
	}
	err = stateDB.WriteChainIDs(crossID, chainIDs)
	if err != nil {
		t.Errorf("write cross %s 's chain_ids error: %s", crossID, err.Error())
	}
	// 写入两个chain的状态
	for _, chainID := range chainIDs {
		dbChainID := chainID
		if err := stateDB.FinishChainCrossState(crossID, dbChainID, []byte(dbChainID), storetypes.StateSuccess); err != nil {
			t.Errorf("write chain cross %s 's chain_id %s error: %s", crossID, chainID, err.Error())
		}
	}
	// 读取两个chain的状态
	for _, chainID := range chainIDs {
		dbChainID := chainID
		chainCrossState, result, exist := stateDB.ReadChainCrossState(crossID, chainID)
		if !exist {
			t.Errorf("read corss[%s] chain[%s] 's state error", crossID, dbChainID)
		}
		if !bytes.Equal(result, []byte(dbChainID)) {
			t.Errorf("read corss[%s] chain[%s] 's result error", crossID, dbChainID)
		}
		if chainCrossState != storetypes.StateSuccess {
			t.Errorf("read corss[%s] chain[%s] 's state wrong", crossID, dbChainID)
		}
	}
}

func newKvStateDB(t *testing.T) *KvStateDB {
	levelDBConfig := newLevelDBConfig()
	dbProvider, err := factory.NewKvDBProvider(storetypes.LevelDB, levelDBConfig)
	if err != nil {
		t.Errorf("init leveldb failed, %v", err)
		t.FailNow()
	}
	stateDB := NewKvStateDB(dbProvider)
	return stateDB
}

func newLevelDBConfig() *conf.LevelDBConfig {
	return &conf.LevelDBConfig{
		StorePath:       "./testdata",
		WriteBufferSize: 4,
		BloomFilterBits: 8,
	}
}
