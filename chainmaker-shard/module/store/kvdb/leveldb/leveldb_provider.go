/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package leveldb

import (
	"fmt"
	"sync"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/logger"
	types "chainmaker.org/chainmaker-cross/store/kvdb/types"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"go.uber.org/zap"
)

const defaultBloomFilterBits = 10

// LevelDBProvider Provider provides handle to db instances
type LevelDBProvider struct {
	sync.Mutex                    // lock
	db         *leveldb.DB        // leveldb
	wo         *opt.WriteOptions  // leveldb options
	logger     *zap.SugaredLogger // log
}

// NewLevelDBProvider construct a new Provider for state operation with given chainId
func NewLevelDBProvider(levelDBConf *conf.LevelDBConfig) *LevelDBProvider {
	if levelDBConf == nil {
		panic("can not create leveldb because it's config is nil")
	}
	dbOpts := &opt.Options{}
	writeBufferSize := levelDBConf.WriteBufferSize
	if writeBufferSize <= 0 {
		//default value 4MB
		dbOpts.WriteBuffer = 4 * opt.MiB
	} else {
		dbOpts.WriteBuffer = writeBufferSize * opt.MiB
	}
	bloomFilterBits := levelDBConf.BloomFilterBits
	if bloomFilterBits <= 0 {
		bloomFilterBits = defaultBloomFilterBits
	}
	dbOpts.Filter = filter.NewBloomFilter(bloomFilterBits)
	dbPath := levelDBConf.StorePath
	db, err := leveldb.OpenFile(conf.FinalCfgPath(dbPath), dbOpts)
	if err != nil {
		panic(fmt.Sprintf("Error opening leveldbprovider: %s", err))
	}
	return &LevelDBProvider{
		db:     db,
		wo:     &opt.WriteOptions{Sync: true},
		logger: logger.GetLogger(logger.ModuleStorage),
	}
}

// Get return value by key
func (l *LevelDBProvider) Get(key string) ([]byte, bool) {
	value, err := l.db.Get([]byte(key), nil)
	if err != nil {
		return nil, false
	}
	return value, true
}

// Put put key and value
func (l *LevelDBProvider) Put(key string, value []byte) error {
	if key == "" {
		return errors.New("error writing leveldb with nil key")
	}
	if value == nil {
		return errors.New("error writing leveldb with nil value")
	}
	err := l.db.Put([]byte(key), value, l.wo)
	if err != nil {
		return err
	}
	return nil
}

// Has return true if this key been existed
func (l *LevelDBProvider) Has(key string) (bool, error) {
	exist, err := l.db.Has([]byte(key), nil)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// Delete delete key from data base
func (l *LevelDBProvider) Delete(key string) error {
	err := l.db.Delete([]byte(key), l.wo)
	if err != nil {
		return err
	}
	return nil
}

// WriteBatch write batch into database
func (l *LevelDBProvider) WriteBatch(batch *types.KvDBBatcher) error {
	if batch.Len() == 0 {
		return errors.New("error writing with nil batch")
	}
	levelBatch := &leveldb.Batch{}
	for _, v := range batch.GetKvs() {
		batchVal := v
		dbKey, dbValue := batchVal.GetKey(), batchVal.GetValue()
		if dbValue == nil {
			// 表示删除
			levelBatch.Delete([]byte(dbKey))
		} else {
			levelBatch.Put([]byte(dbKey), dbValue)
		}
	}
	err := l.db.Write(levelBatch, l.wo)
	if err != nil {
		return err
	}
	return nil
}

// Close close the leveldb
func (l *LevelDBProvider) Close() {
	if err := l.db.Close(); err != nil {
		l.logger.Error("close leveldb failed", err)
	}
	l.logger.Info("Module storage stopped")
}
