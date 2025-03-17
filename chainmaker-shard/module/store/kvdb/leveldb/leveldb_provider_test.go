/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package leveldb

import (
	"strconv"
	"testing"

	"chainmaker.org/chainmaker-cross/conf"
)

func TestNewLevelDBProvider(t *testing.T) {
	levelDBConfig := newLevelDBConfig()
	dbProvider := NewLevelDBProvider(levelDBConfig)
	if dbProvider == nil {
		t.Error("leveldb is nil")
		t.FailNow()
	}
	// 完成后关闭
	dbProvider.Close()
}

func TestLevelDBProvider_PutAndGetAndHashAndDelete(t *testing.T) {
	levelDBConfig := newLevelDBConfig()
	dbProvider := NewLevelDBProvider(levelDBConfig)
	keyValues := make(map[string]string)
	for i := 0; i < 10; i++ {
		keyValues["key"+strconv.Itoa(i)] = "value" + strconv.Itoa(i)
	}
	// TestPut
	for k, v := range keyValues {
		if err := dbProvider.Put(k, []byte(v)); err != nil {
			t.Errorf("put %s -> %s failed", k, v)
		}
	}
	// TestGet
	for k, v := range keyValues {
		// TestPut
		bytes, exist := dbProvider.Get(k)
		if !exist {
			t.Errorf("can not find value for %s", k)
		} else if string(bytes) != v {
			t.Errorf("%s 's value is not right", k)
		}
	}
	if _, exist := dbProvider.Get("test"); exist {
		t.Error("db's data is error")
	}
	// TestHas
	for k := range keyValues {
		// TestPut
		exist, err := dbProvider.Has(k)
		if err != nil {
			t.Errorf("check %s exist error:%s", k, err.Error())
		}
		if !exist {
			t.Errorf("can not find value for %s", k)
		}
	}
	if _, exist := dbProvider.Get("test"); exist {
		t.Error("db's data is error")
	}
	// TestDelete
	for k := range keyValues {
		// TestPut
		err := dbProvider.Delete(k)
		if err != nil {
			t.Errorf("delete %s from db error:%s", k, err.Error())
		}
	}
	// TestGetAgain
	for k := range keyValues {
		// TestPut
		_, exist := dbProvider.Get(k)
		if exist {
			t.Errorf("find value for %s", k)
		}
	}
	if _, exist := dbProvider.Get("test"); exist {
		t.Error("db's data is error")
	}
	// 完成后关闭
	dbProvider.Close()
}

func TestLevelDBProvider_WriteBatch(t *testing.T) {
	//levelDBConfig := newLevelDBConfig()
	//dbProvider := NewLevelDBProvider(levelDBConfig)
	//if dbProvider == nil {
	//	t.Error("leveldb is nil")
	//	t.FailNow()
	//}
	//kvDBBatcher := kvdb.NewKvDBBatcher()
	//for i := 0; i < 10; i++ {
	//	kvDBBatcher.Add("batch" + strconv.Itoa(i), []byte("value" + strconv.Itoa(i)))
	//}
	//err := dbProvider.WriteBatch(kvDBBatcher)
	//if err != nil {
	//	t.Errorf("write batch error %s", err.Error())
	//}
	//// TestGet
	//for _, kv := range kvDBBatcher.GetKvs() {
	//	// TestPut
	//	dbVal, exist := dbProvider.Get(kv.GetKey())
	//	if !exist {
	//		t.Errorf("can not find value for %s", kv.GetKey())
	//	}
	//
	//	if !bytes.Equal(dbVal, kv.GetValue()) {
	//		t.Errorf("%s 's value is not right", kv.GetKey())
	//	}
	//}
	//// 完成后关闭
	//dbProvider.Close()
}

func TestLevelDBProvider_Close(t *testing.T) {
	levelDBConfig := newLevelDBConfig()
	dbProvider := NewLevelDBProvider(levelDBConfig)
	keyValues := make(map[string]string)
	for i := 0; i < 10; i++ {
		keyValues["key"+strconv.Itoa(i)] = "value" + strconv.Itoa(i)
	}
	// TestPut
	for k, v := range keyValues {
		if err := dbProvider.Put(k, []byte(v)); err != nil {
			t.Errorf("put %s -> %s failed", k, v)
		}
	}
	// TestClose
	dbProvider.Close()
	// TestGet
	err := dbProvider.Put("test", []byte("test-value"))
	if err == nil {
		t.Error("db close's error")
	}
}

func newLevelDBConfig() *conf.LevelDBConfig {
	return &conf.LevelDBConfig{
		StorePath:       "./testdata",
		WriteBufferSize: 4,
		BloomFilterBits: 8,
	}
}
