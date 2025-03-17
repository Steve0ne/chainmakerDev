/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"testing"

	"chainmaker.org/chainmaker-cross/channel"
	"chainmaker.org/chainmaker-cross/event"
	"github.com/stretchr/testify/require"
)

func TestNewRouters(t *testing.T) {
	event.InitLog(getLogger())
	routers := NewRouters()
	require.NotNil(t, routers)
	// 创建innerRouter
	var chainIDs = []string{"chain1"}
	innerRouter := GetInnerRouter()
	innerRouter.Init(chainIDs)
	err := routers.Add(innerRouter)
	require.Nil(t, err)
	// 创建channelRouter
	connection := newConnectionMock()
	netChannel := channel.NewNetChannel(connection)
	err = netChannel.Init()
	require.Nil(t, err)
	channelRouter := NewChannelRouter(chainIDs, netChannel)
	err = routers.Add(channelRouter)
	require.Nil(t, err)
}

func TestRouters(t *testing.T) {
	event.InitLog(getLogger())
	routers := NewRouters()
	require.NotNil(t, routers)
	// 创建innerRouter
	var chainIDs = []string{"chain1"}
	innerRouter := GetInnerRouter()
	innerRouter.Init(chainIDs)
	err := routers.Add(innerRouter)
	require.Nil(t, err)
	// 创建channelRouter
	connection := newConnectionMock()
	netChannel := channel.NewNetChannel(connection)
	err = netChannel.Init()
	require.Nil(t, err)
	channelRouter := NewChannelRouter(chainIDs, netChannel)
	err = routers.Add(channelRouter)
	require.Nil(t, err)
	// check
	newInnerRouter, exist := routers.GetInnerRouter()
	require.Equal(t, innerRouter, newInnerRouter)
	require.True(t, exist)
	newChannelRouter, exist := routers.GetChannelRouter()
	require.Equal(t, channelRouter, newChannelRouter)
	require.True(t, exist)
	require.True(t, routers.InnerSupport())
	require.True(t, routers.ChannelSupport())
	require.True(t, routers.Support())
}
