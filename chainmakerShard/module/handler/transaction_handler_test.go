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

func TestTransactionProcessHandler(t *testing.T) {
	// test get chain call
	TPH := GetTransactionProcessHandler()

	// test set state
	TPH.SetStateDB(nil)

	// test set logger
	log := logger.GetLogger(logger.ModuleHandler)
	TPH.SetLogger(log)

	// test get type
	ht := TPH.GetType()
	require.Equal(t, ht, TransactionProcess)
}
