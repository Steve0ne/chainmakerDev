/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"testing"
	"time"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/utils"
	"github.com/stretchr/testify/require"
)

func TestGetDispatcher(t *testing.T) {
	routerDispatcher := GetDispatcher()
	require.NotNil(t, routerDispatcher)
	testLogger := getLogger()
	event.InitLog(testLogger)
	routerDispatcher.SetLogger(testLogger)
}

func TestRouterDispatcher_Invoke(t *testing.T) {
	routerDispatcher := GetDispatcher()
	require.NotNil(t, routerDispatcher)
	testLogger := getLogger()
	event.InitLog(testLogger)
	routerDispatcher.SetLogger(testLogger)
	// 注册innerrouter
	var chainIDs = []string{"chain1", "chain2"}
	innerRouter := GetInnerRouter()
	innerRouter.Init(chainIDs)
	err := routerDispatcher.Register(innerRouter)
	require.Nil(t, err)
	crossID := utils.NewUUID()
	transactionEvent := event.NewExecuteTransactionEvent(crossID, "chain1", []byte(""), nil)
	response, err := routerDispatcher.Invoke(transactionEvent, time.Second)
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, false, response.IsSuccess())
	require.Equal(t, crossID, response.GetCrossID())
	require.Equal(t, "chain1", response.GetChainID())

	crossID = utils.NewUUID()
	transactionEvent = event.NewExecuteTransactionEvent(crossID, "chain3", []byte(""), nil)
	response, err = routerDispatcher.Invoke(transactionEvent, time.Second)
	require.NotNil(t, err)
	require.Nil(t, response)
}
