/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package memory

import (
	"sync"

	kvdbtypes "chainmaker.org/chainmaker-cross/store/kvdb/types"
)

// MemProvider storage of memory
type MemProvider struct {
	sync.RWMutex
	cache map[string][]byte
}

// NewMemProvider create new memory provider
func NewMemProvider() *MemProvider {
	return &MemProvider{
		cache: make(map[string][]byte),
	}
}

// Get return value for this key
func (m *MemProvider) Get(key string) ([]byte, bool) {
	m.RLock()
	defer m.RUnlock()
	value, exist := m.cache[key]
	return value, exist
}

// Put put the key and value into memory
func (m *MemProvider) Put(key string, value []byte) error {
	m.Lock()
	defer m.Unlock()
	m.cache[key] = value
	return nil
}

// Has return true if this key been existed
func (m *MemProvider) Has(key string) (bool, error) {
	m.RLock()
	defer m.RUnlock()
	_, exist := m.cache[key]
	return exist, nil
}

// Delete delete the key from memory
func (m *MemProvider) Delete(key string) error {
	m.Lock()
	defer m.Unlock()
	delete(m.cache, key)
	return nil
}

// WriteBatch write batch into memory
func (m *MemProvider) WriteBatch(batch *kvdbtypes.KvDBBatcher) error {
	kvs := batch.GetKvs()
	m.Lock()
	defer m.Unlock()
	for i := 0; i < len(kvs); i++ {
		kv := kvs[i]
		m.cache[kv.GetKey()] = kv.GetValue()
	}
	return nil
}

// Close clear the memory
func (m *MemProvider) Close() {
	m.cache = nil
}
