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

func TestCrossSearchHandler(t *testing.T) {
	// test get chain call
	CSH := GetCrossSearchHandler()

	// test set state
	CSH.SetStateDB(nil)

	// test set logger
	log := logger.GetLogger(logger.ModuleHandler)
	CSH.SetLogger(log)

	// test get type
	ht := CSH.GetType()
	require.Equal(t, ht, CrossSearch)
}
