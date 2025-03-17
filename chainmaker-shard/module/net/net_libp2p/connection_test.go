/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net_libp2p

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const localAddress = "/ip4/127.0.0.1/tcp/12345/p2p/QmSVuE22fYPeTYg4iUJJbVmwtaumYavSApR3PSm4Gb6upV"

func TestConnection(t *testing.T) {
	// Prepare environment
	node, err := mockHost(t)
	require.NoError(t, err)

	// Test listen
	_, err = node.Listen()
	require.NoError(t, err)

	// Test new connection
	connection, err := NewLibP2pConnection(localAddress, testPid, testDelimit, 10, 1000)
	require.NoError(t, err)

	// Test write data
	var cnt int
	go func() {
		//for {
		cnt++
		// auto write data
		msg := NewLibP2pMessage(connection.PeerID(), []byte(fmt.Sprintf("This is test!, num is:%d", cnt)), false)
		err = connection.WriteData(msg)
		t.Log(fmt.Sprintf("write data! %d", cnt))
		if err != nil {
			t.Log("err is: ", err)
		}

		time.Sleep(time.Second * 3)
		if cnt > 100 {
			return
		}
	}()

	// Test read data
	ch, err := connection.ReadData()
	require.NoError(t, err)
	go func() {
		//for {
		resp := <-ch
		strs := strings.Split(string(resp.GetPayload()), ":")
		if len(strs) < 2 {
			//continue
		}
		t.Log("Receive response!, num is:", strs[1])
		time.Sleep(time.Second * 2)
		//}
	}()
	// hold func
	//select {}
}

func TestMultiConnection(t *testing.T) {
	// Prepare environment
	node, err := mockHost(t)
	require.NoError(t, err)

	// Test listen
	_, err = node.Listen()
	require.NoError(t, err)

	// Test new connection
	var cnt int
	for {
		// loop create and reset peer
		connection, err := NewLibP2pConnection(localAddress, testPid, testDelimit, 10, 1000)
		require.NoError(t, err)

		// read data
		_, err = connection.ReadData()
		require.NoError(t, err)

		// new message
		msg := NewLibP2pMessage(connection.PeerID(), []byte(fmt.Sprintf("This is test!, num is:%d", cnt)), false)

		// write data
		err = connection.WriteData(msg)
		require.NoError(t, err)

		cnt++
		if cnt >= 3 {
			break
		}
	}
}

func TestResetConnection(t *testing.T) {
	// Prepare environment
	node, err := mockHost(t)
	require.NoError(t, err)

	// Test listen
	_, err = node.Listen()
	require.NoError(t, err)

	// Test new connection
	connection, err := NewLibP2pConnection(localAddress, testPid, testDelimit, 10, 1000)
	require.NoError(t, err)

	// Test write data
	var cnt int
	go func() {
		//for {
		cnt++
		// auto write data
		msg := NewLibP2pMessage(connection.PeerID(), []byte(fmt.Sprintf("This is test!, num is:%d", cnt)), false)
		t.Log(fmt.Sprintf("write data! %d", cnt))
		err = connection.WriteData(msg)
		t.Log("current connection peer id: ", connection.PeerID())
		time.Sleep(time.Second * 3)
		if cnt > 10 {
			return
		}
		//}
	}()

	// Test read data
	_, err = connection.ReadData()
	require.NoError(t, err)

}

func TestTwoConnect(t *testing.T) {
	// Prepare environment
	node, err := mockHost(t)
	require.NoError(t, err)

	// Test listen
	_, err = node.Listen()
	require.NoError(t, err)

	_, err = NewLibP2pConnection(localAddress, testPid, testDelimit, 10, 1000)
	require.NoError(t, err)
	//select {}
}

func TestTwoListener(t *testing.T) {
	// Prepare environment
	node, err := mockHost(t)
	require.NoError(t, err)

	// Test listen
	_, err = node.Listen()
	require.NoError(t, err)

	connection, err := NewLibP2pConnection(localAddress, testPid, testDelimit, 10, 1000)
	require.NoError(t, err)
	_, err = connection.ReadData()
	require.NoError(t, err)
	//select {}
}
