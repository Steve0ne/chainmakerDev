/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net_libp2p

import (
	"io/ioutil"
	"os"
	"testing"

	"chainmaker.org/chainmaker-cross/net"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/stretchr/testify/require"
)

const privatePEMStr = "-----BEGIN EC PRIVATE KEY-----\n" +
	"MHcCAQEEIKwz2DjJZWtxnHNtMroZEjr28dqihVAPcyrmQp/5XeGRoAoGCCqGSM49\n" +
	"AwEHoUQDQgAERrAT6tdtd95gd6zh2/qics3s/h7eUtAGLdPhuGstl1j9iqSXI/0i\n" +
	"9qyw+HHd2IkkLdQX6ivmaovGS0vbSDpRCA==\n" +
	"-----END EC PRIVATE KEY-----\n"
const privateKeyHex = "0803127930770201010420ac33d838c9656b719c736d32ba19123af6f1daa285500f732ae6429ff95de191a00a06082a8648ce3d030107a1440342000446b013ead76d77de6077ace1dbfaa272cdecfe1ede52d0062dd3e1b86b2d9758fd8aa49723fd22f6acb0f871ddd889242dd417ea2be66a8bc64b4bdb483a5108"
const testLocal = "/ip4/127.0.0.1/tcp/12345"
const testPid = protocol.ID("/p2p/1.0.0")
const testDelimit = '\n'

// TestHost - go through host utilities
func TestHost(t *testing.T) {
	// Test create dummy node
	//_, err := NewDummyHost(testPid, '\n')

	// Test create mock node
	mock, err := mockHost(t)
	require.NoError(t, err)

	// Test listen
	_, err = mock.Listen()
	require.NoError(t, err)

	// Test new message
	msg := NewLibP2pMessage("QmSVuE22fYPeTYg4iUJJbVmwtaumYavSApR3PSm4Gb6upV", []byte("This is a test!"), false)

	// Test write
	err = mock.Write(msg)
	require.NoError(t, err)

	// Test stop
	err = mock.Stop()
	require.NoError(t, err)
}

// mockHost - create mock host with fixed address: /ip4/127.0.0.1/tcp/12345/p2p/QmSVuE22fYPeTYg4iUJJbVmwtaumYavSApR3PSm4Gb6upV
func mockHost(t *testing.T) (net.Peer, error) {
	var err error
	// 创建临时文件目录
	p, err := ioutil.TempDir(os.TempDir(), "mock")
	require.NoError(t, err)
	// 删除目录
	defer os.RemoveAll(p)
	path := p + "/PEM"
	// 写入私钥
	err = ioutil.WriteFile(path, []byte(privatePEMStr), 0600)
	require.NoError(t, err)
	// 生成模拟节点
	node, err := NewMonitorHost(testLocal, path, testPid, testDelimit)
	require.NoError(t, err)
	t.Log("local libp2p peer id: ", node.Host.ID().Pretty())
	return node, nil
}

// TestHostInteraction - test interactions with connection I\O
func TestHostConnectionInteractions(t *testing.T) {
	//// Prepare environment
	//node, err := mockHost(t)
	//require.NoError(t, err)
	//
	//// Test listen
	//channel, err := node.Listen()
	//require.NoError(t, err)

	// Set read data response
	//for {
	//	// loop read data
	//	msg := <- channel
	//	// parse data
	//	strs := strings.Split(string(msg.GetPayload()), ":")
	//	if len(strs) < 2 {
	//		continue
	//	}
	//	t.Log("Receive data, write response")
	//	m := NewLibP2pMessage(msg.GetNodeID(), []byte(fmt.Sprintf("This is a response!, num is:%s", strs[1])), false)
	//	err := node.Write(m)
	//	require.NoError(t, err)
	//	time.Sleep(time.Second * 3)
	//}
}
