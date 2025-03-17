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

func TestCrossProcessHandler(t *testing.T) {
	// test get chain call
	CPH := GetCrossProcessHandler()

	// test set state
	CPH.SetStateDB(nil)

	// test set logger
	log := logger.GetLogger(logger.ModuleHandler)
	CPH.SetLogger(log)

	// test get type
	ht := CPH.GetType()
	require.Equal(t, ht, CrossProcess)
}
