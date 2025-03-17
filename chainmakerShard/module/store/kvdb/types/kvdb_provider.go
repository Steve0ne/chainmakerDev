/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package types

// KvDBProvider is the provider for kvdb
type KvDBProvider interface {

	// Get get value from database for key
	Get(key string) ([]byte, bool)

	// Put saves the key-values
	Put(key string, value []byte) error

	// Has return true if the given key exist, or return false if none exists
	Has(key string) (bool, error)

	// Delete deletes the given key
	Delete(key string) error

	// WriteBatch writes a batch in an atomic operation
	WriteBatch(batch *KvDBBatcher) error

	// Close close the database
	Close()
}

// KvDBBatcher the struct for batch key-values
type KvDBBatcher struct {
	// 不加锁，由调用方处理
	kvs []*Kv // 键值对数组
}

// NewKvDBBatcher create new KvDBBatcher
func NewKvDBBatcher() *KvDBBatcher {
	return &KvDBBatcher{
		kvs: make([]*Kv, 0),
	}
}

// Add add key-value to KvDBBatcher
func (b *KvDBBatcher) Add(key string, value []byte) {
	b.kvs = append(b.kvs, NewKv(key, value))
}

// GetKvs return all the key-values
func (b *KvDBBatcher) GetKvs() []*Kv {
	return b.kvs
}

// Len return length of kvs
func (b *KvDBBatcher) Len() int {
	return len(b.kvs)
}

// Kv the key-value
type Kv struct {
	key   string
	value []byte
}

// NewKv create new Kv instance
func NewKv(key string, value []byte) *Kv {
	return &Kv{
		key:   key,
		value: value,
	}
}

// GetKey return key of KV
func (kv *Kv) GetKey() string {
	return kv.key
}

// GetValue return value of KV
func (kv *Kv) GetValue() []byte {
	return kv.value
}
