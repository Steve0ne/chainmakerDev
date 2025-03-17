/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package inner_listener

import (
	"testing"

	"chainmaker.org/chainmaker-cross/handler"
	"github.com/stretchr/testify/require"
)

func TestInnerListenStart(t *testing.T) {
	// Test new inner listener
	l := NewInnerListener()
	handler.InitEventHandlers(nil, nil)

	// Test start
	err := l.ListenStart()
	require.NoError(t, err)
}

func TestInnerListenStop(t *testing.T) {
	// Test new inner listener
	l := NewInnerListener()
	handler.InitEventHandlers(nil, nil)

	// Test start
	err := l.ListenStart()
	require.NoError(t, err)

	// Test stop
	err = l.Stop()
	require.NoError(t, err)
}
