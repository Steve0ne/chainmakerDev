/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"testing"

	"chainmaker.org/chainmaker-cross/logger"
	"github.com/stretchr/testify/require"
)

func TestChainCallHandler(t *testing.T) {
	// test get chain call
	CCH := GetChainCallHandler()

	// test set state
	CCH.SetStateDB(nil)

	// test set logger
	log := logger.GetLogger(logger.ModuleHandler)
	CCH.SetLogger(log)

	// test get type
	ht := CCH.GetType()
	require.Equal(t, ht, ChainCall)
}
