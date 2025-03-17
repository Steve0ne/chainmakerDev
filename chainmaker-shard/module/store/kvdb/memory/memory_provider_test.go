/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package memory

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMemProvider(t *testing.T) {
	// test new MemProvider
	mp := NewMemProvider()
	require.NotNil(t, mp)

	// prepare data
	keyValues := make(map[string]string)
	for i := 0; i < 10; i++ {
		keyValues["key"+strconv.Itoa(i)] = "value" + strconv.Itoa(i)
	}

	// TestPut
	for k, v := range keyValues {
		if err := mp.Put(k, []byte(v)); err != nil {
			t.Errorf("put %s -> %s failed", k, v)
		}
	}
	// TestGet
	for k, v := range keyValues {
		// TestPut
		bytes, exist := mp.Get(k)
		if !exist {
			t.Errorf("can not find value for %s", k)
		} else if string(bytes) != v {
			t.Errorf("%s 's value is not right", k)
		}
	}
	if _, exist := mp.Get("test"); exist {
		t.Error("db's data is error")
	}
	// TestHas
	for k := range keyValues {
		// TestPut
		exist, err := mp.Has(k)
		if err != nil {
			t.Errorf("check %s exist error:%s", k, err.Error())
		}
		if !exist {
			t.Errorf("can not find value for %s", k)
		}
	}
	if _, exist := mp.Get("test"); exist {
		t.Error("db's data is error")
	}
	// TestDelete
	for k := range keyValues {
		// TestPut
		err := mp.Delete(k)
		if err != nil {
			t.Errorf("delete %s from db error:%s", k, err.Error())
		}
	}
	// TestGetAgain
	for k := range keyValues {
		// TestPut
		_, exist := mp.Get(k)
		if exist {
			t.Errorf("find value for %s", k)
		}
	}
	if _, exist := mp.Get("test"); exist {
		t.Error("db's data is error")
	}

	// 完成后关闭
	mp.Close()
}
