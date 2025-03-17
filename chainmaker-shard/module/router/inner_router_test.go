/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"testing"
	"time"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetInnerRouter(t *testing.T) {
	innerRouter := GetInnerRouter()
	require.NotNil(t, innerRouter)
	require.Equal(t, InnerRouterType, innerRouter.GetType())
}

func TestInnerRouter_Init(t *testing.T) {
	var chainIDs = []string{"chain1", "chain2"}
	innerRouter := GetInnerRouter()
	innerRouter.Init(chainIDs)
	newChainIDs := innerRouter.GetChainIDs()
	require.Equal(t, len(chainIDs), len(newChainIDs))
	require.Equal(t, chainIDs[0], newChainIDs[0])
	require.Equal(t, chainIDs[1], newChainIDs[1])
}

func TestInnerRouter_Invoke(t *testing.T) {
	var chainIDs = []string{"chain1", "chain2"}
	event.InitLog(getLogger())
	innerRouter := GetInnerRouter()
	innerRouter.Init(chainIDs)
	crossID := utils.NewUUID()
	transactionEvent := event.NewExecuteTransactionEvent(crossID, "chain1", []byte(""), nil)
	response, err := innerRouter.Invoke(transactionEvent, time.Second)
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, false, response.IsSuccess())
	require.Equal(t, crossID, response.GetCrossID())
	require.Equal(t, "chain1", response.GetChainID())
}

func getLogger() *zap.SugaredLogger {
	config := []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			LogLevel:     logger.INFO,
			FilePath:     "logs/default.log",
			MaxAge:       365,
			RotationTime: 1,
			LogInConsole: false,
			ShowColor:    true,
		},
	}
	logger.InitLogConfig(config)
	return logger.GetLogger("default")
}
