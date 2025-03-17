/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitEventHandler(t *testing.T) {
	InitEventHandlers(nil, nil)
}

func TestEventHandlerTools(t *testing.T) {
	EHT := GetEventHandlerTools()
	// 注册ChainCallHandler
	EHT.Register(getChainCallHandler(nil))
	// 注册TransactionProcessHandler
	EHT.Register(getTransactionProcessHandler(nil))
	// 注册CrossProcessHandler
	EHT.Register(getCrossProcessHandler(nil, nil))
	// 注册CrossSearchHandler
	EHT.Register(getCrossSearchHandler(nil))

	CCEH, ok := EHT.GetHandler(ChainCall)
	require.Equal(t, ok, true)
	require.NotNil(t, CCEH)
	TPEH, ok := EHT.GetHandler(TransactionProcess)
	require.NotNil(t, TPEH)
	require.Equal(t, ok, true)
	CPEH, ok := EHT.GetHandler(CrossProcess)
	require.Equal(t, ok, true)
	require.NotNil(t, CPEH)
	CSEH, ok := EHT.GetHandler(CrossSearch)
	require.Equal(t, ok, true)
	require.NotNil(t, CSEH)
}
