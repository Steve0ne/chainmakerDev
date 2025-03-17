/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUUID(t *testing.T) {
	var length = 10
	uuidMap := make(map[string]string)
	for i := 0; i < length; i++ {
		uuidMap[NewUUID()] = "test"
	}
	if len(uuidMap) != length {
		t.Error("uuid is repeat")
	}
}

func TestNewRandomKey(t *testing.T) {
	key := NewRandomKey()
	require.NotNil(t, key)
}
