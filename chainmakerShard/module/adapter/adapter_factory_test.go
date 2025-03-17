/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package adapter

import (
	"fmt"
	"testing"

	"chainmaker.org/chainmaker-cross/adapter/chainmaker"
	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/logger"
	"github.com/stretchr/testify/require"
)

func TestInitAdapters(t *testing.T) {
	// mock config
	conf.Config.AdapterConfigs = []*conf.AdapterConfig{
		{
			Provider:   "chainmaker",
			ChainID:    "chain1",
			ConfigPath: "", // TODO
		},
	}
	// init
	adapter := InitAdapters()
	require.NotNil(t, adapter)
}

func TestGetChainAdapterDispatcher(t *testing.T) {
	ad := GetChainAdapterDispatcher()
	require.NotNil(t, ad)
	chainAdapter, err := chainmaker.NewChainMakerAdapter("chainID", "", nil, nil)
	require.Equal(t, err, fmt.Errorf("connect chainmaker node address is empty"))
	require.Nil(t, chainAdapter)
}

func TestChainAdapterDispatcher_Register(t *testing.T) {
	ad := GetChainAdapterDispatcher()
	mockAdapter := chainmaker.MockNilChainMakerAdapter()
	log := logger.GetLogger(logger.ModuleAdapter)
	ad.SetLog(log)
	ad.Register(mockAdapter)
}

func TestChainAdapterDispatcher_Invoke(t *testing.T) {
	// test nil
	ad := GetChainAdapterDispatcher()
	mockAdapter := chainmaker.MockNilChainMakerAdapter()
	log := logger.GetLogger(logger.ModuleAdapter)
	ad.SetLog(log)
	ad.Register(mockAdapter)
	resp, err := ad.Invoke("chainIDNotExist", event.NewCommitTransactionEvent("crossID", "chainID", []byte{}))
	require.Equal(t, err, fmt.Errorf("can not find adapter for chain[chainIDNotExist]"))
	require.Nil(t, resp)
}

func TestChainAdapterDispatcher_QueryByTxKey(t *testing.T) {
	// test nil
	ad := GetChainAdapterDispatcher()
	log := logger.GetLogger(logger.ModuleAdapter)
	ad.SetLog(log)
	resp, err := ad.QueryByTxKey("chainIDNotExist", "txKey")
	require.Equal(t, err, fmt.Errorf("can not find adapter for chain[chainIDNotExist]"))
	require.Nil(t, resp)
}

func TestChainAdapterDispatcher_GetChainIDs(t *testing.T) {
	// test nil
	ad := GetChainAdapterDispatcher()
	log := logger.GetLogger(logger.ModuleAdapter)
	ad.SetLog(log)
	ids := ad.GetChainIDs()
	require.Equal(t, ids, []string{"chainID"})

	// test set ID
	mockAdapter := chainmaker.MockNilChainMakerAdapter()
	ad.Register(mockAdapter)
	ids = ad.GetChainIDs()
	require.Equal(t, ids, []string{"chainID"})
}

func TestChainAdapterDispatcher_Query(t *testing.T) {
	// test nil
	ad := GetChainAdapterDispatcher()
	log := logger.GetLogger(logger.ModuleAdapter)
	ad.SetLog(log)
	resp, err := ad.Query("chainIDNotExist", []byte{})
	require.Equal(t, err, fmt.Errorf("can not find adapter for chain[chainIDNotExist]"))
	require.Nil(t, resp)
}

func TestChainAdapterDispatcher_SetLog(t *testing.T) {
	ad := GetChainAdapterDispatcher()
	log := logger.GetLogger(logger.ModuleAdapter)
	ad.SetLog(log)
}
