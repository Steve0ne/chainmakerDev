/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"context"
	"fmt"
	"testing"
	"time"

	"chainmaker.org/chainmaker-cross/sdk/builder"
	"github.com/stretchr/testify/require"
)

func Test_UserGuide(t *testing.T) {
	//生成CrossSDK实例
	crossSDK, err := NewCrossSDK(WithConfigFile("./config/template/cross_chain_sdk.yml"))
	require.NoError(t, err)
	require.NotNil(t, crossSDK)
	//构造跨链事件
	crossEvent, err := crossSDK.GenCrossEvent(NewCrossTxBuildCtx(
		"chain1",
		0,
		builder.NewContract("balance_003", "Plus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
		builder.NewContract("balance_003", "Minus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
	), NewCrossTxBuildCtx(
		"chain2",
		1,
		builder.NewContract("balance_003", "GetProof", builder.NewParamsWithMap(map[string]string{"number": "1"})),
		builder.NewContract("balance_003", "Minus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
		builder.WithUseProofKey(true),
	))
	require.NoError(t, err)
	ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
	//发送跨链事件

	res, err := crossSDK.SendCrossEvent(crossEvent, "http://localhost:8080", true, WithContextOpt(ctx), WithSyncStrategyOpt(SyncStrategy{
		0, 2 * time.Second, 1 * time.Second,
	}))
	require.NoError(t, err)
	if res != nil {
		fmt.Printf("%+v\n", *res)
	}
}

func TestCrossChainMakerAndChainMaker(t *testing.T) {
	//生成CrossSDK实例
	crossSDK, err := NewCrossSDK(WithConfigFile("/root/workspace/ChainMaker/chainmaker-cross-chain/release/config/cross_chain_sdk.yml"))
	require.NoError(t, err)
	require.NotNil(t, crossSDK)
	//构造跨链事件
	crossEvent, err := crossSDK.GenCrossEvent(NewCrossTxBuildCtx(
		"chain1", 0,
		builder.NewContract("BalanceStable", "Minus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
		builder.NewContract("BalanceStable", "Reset", nil),
	), NewCrossTxBuildCtx(
		"chain2", 1,
		builder.NewContract("BalanceStable", "Minus", builder.NewParamsNoKeys("123", "123", "123")),
		builder.NewContract("BalanceStable", "Reset", builder.NewParams(builder.NewKV("para1", "123"), builder.NewKV("para2", "123"))),
	))
	require.NoError(t, err)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//发送跨链事件
	res, err := crossSDK.SendCrossEvent(crossEvent, "https://localhost:8080", true, WithContextOpt(ctx), WithSyncStrategyOpt(SyncStrategy{
		0, 2 * time.Second, 1 * time.Second,
	}))
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestCrossChainMakerAndFabric(t *testing.T) {
	//生成CrossSDK实例
	crossSDK, err := NewCrossSDK(WithConfigFile("/root/workspace/ChainMaker/cross-chain/release/config/cross_chain_sdk.yml"))
	require.NoError(t, err)
	require.NotNil(t, crossSDK)
	//构造跨链事件
	crossEvent, err := crossSDK.GenCrossEvent(
		NewCrossTxBuildCtx(
			"chain1", 0,
			builder.NewContract("Balance", "Plus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
			builder.NewContract("Balance", "Reset", nil),
		),
		NewCrossTxBuildCtx(
			"mychannel", 1,
			builder.NewContract("fabcar", "ChangeCarOwner", builder.NewParamsNoKeys("CAR0", "OWNER")),
			builder.NewContract("fabcar", "QueryAllCars", nil),
		),
	)
	require.NoError(t, err)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	//发送跨链事件
	res, err := crossSDK.SendCrossEvent(crossEvent, "https://localhost:8080", true, WithContextOpt(ctx), WithSyncStrategyOpt(SyncStrategy{
		0, 20 * time.Second, 20 * time.Second,
	}))
	require.NoError(t, err)
	require.NotNil(t, res)
}
