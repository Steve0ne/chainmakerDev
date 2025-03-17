/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"testing"
	"time"

	"chainmaker.org/chainmaker-cross/channel"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/net"
	"chainmaker.org/chainmaker-cross/utils"
	"github.com/stretchr/testify/require"
)

func TestNewChannelRouter(t *testing.T) {
	event.InitLog(getLogger())
	connection := newConnectionMock()
	defer connection.Close()
	netChannel := channel.NewNetChannel(connection)
	err := netChannel.Init()
	require.Nil(t, err)
	var chainIDs = []string{"chain1", "chain2"}
	channelRouter := NewChannelRouter(chainIDs, netChannel)
	require.NotNil(t, channelRouter)
	require.Equal(t, chainIDs, channelRouter.GetChainIDs())
	require.Equal(t, ChannelRouterType, channelRouter.GetType())
}

func TestChannelRouter_Invoke(t *testing.T) {
	event.InitLog(getLogger())
	connection := newConnectionMock()
	netChannel := channel.NewNetChannel(connection)
	err := netChannel.Init()
	require.Nil(t, err)
	var chainIDs = []string{"chain1", "chain2"}
	crossID := utils.NewUUID()
	channelRouter := NewChannelRouter(chainIDs, netChannel)
	transactionEvent := event.NewExecuteTransactionEvent(crossID, "chain1", []byte(""), "", nil)
	response, err := channelRouter.Invoke(transactionEvent, time.Second)
	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, false, response.IsSuccess())
	require.Equal(t, crossID, response.GetCrossID())
	require.Equal(t, "chain1", response.GetChainID())
}

type ConnectionMock struct {
	readChan  chan net.Message
	writeChan chan net.Message
}

func newConnectionMock() *ConnectionMock {
	return &ConnectionMock{
		readChan:  make(chan net.Message, 1024),
		writeChan: make(chan net.Message, 1024),
	}
}

func (c *ConnectionMock) PeerID() string {
	return "local"
}

func (c *ConnectionMock) ReadData() (chan net.Message, error) {
	return c.readChan, nil
}

func (c *ConnectionMock) WriteData(message net.Message) error {
	c.writeChan <- message
	return nil
}

func (c *ConnectionMock) GetProvider() net.ConnectionProvider {
	return "mock"
}

func (c *ConnectionMock) Close() error {
	close(c.writeChan)
	close(c.readChan)
	return nil
}
