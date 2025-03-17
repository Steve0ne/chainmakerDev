/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"testing"

	"chainmaker.org/chainmaker-cross/event"
	"github.com/stretchr/testify/require"
)

func TestRouterFactory(t *testing.T) {
	testLogger := getLogger()
	event.InitLog(testLogger)
	var chainIDs = []string{"chain1", "chain2"}
	routerDispatcher := InitRouters(chainIDs)
	require.NotNil(t, routerDispatcher)
	router, b := routerDispatcher.getInnerRouter("chain1")
	require.NotNil(t, router)
	require.True(t, b)
	router, b = routerDispatcher.getInnerRouter("chain3")
	require.Nil(t, router)
	require.False(t, b)
}
