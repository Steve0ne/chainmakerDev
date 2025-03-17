/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package channel_listener

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"chainmaker.org/chainmaker-cross/net/net_libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/require"
)

const privatePEMStr = "-----BEGIN EC PRIVATE KEY-----\n" +
	"MHcCAQEEIKwz2DjJZWtxnHNtMroZEjr28dqihVAPcyrmQp/5XeGRoAoGCCqGSM49\n" +
	"AwEHoUQDQgAERrAT6tdtd95gd6zh2/qics3s/h7eUtAGLdPhuGstl1j9iqSXI/0i\n" +
	"9qyw+HHd2IkkLdQX6ivmaovGS0vbSDpRCA==\n" +
	"-----END EC PRIVATE KEY-----\n"
const localAddress = "/ip4/127.0.0.1/tcp/12345/p2p/QmSVuE22fYPeTYg4iUJJbVmwtaumYavSApR3PSm4Gb6upV"
const testLocal = "/ip4/127.0.0.1/tcp/12345"
const testPid = "/p2p/1.0.0"
const testDelimit = '\n'

func TestListenStart(t *testing.T) {
	// prepare environment
	path := preparePEM(t)
	defer os.RemoveAll(path)
	// new host
	node, err := net_libp2p.NewMonitorHost(testLocal, path, testPid, testDelimit)
	require.NoError(t, err)
	_, err = node.Listen()
	require.NoError(t, err)
}

func TestProcessReadData(t *testing.T) {
	// prepare environment
	//path := preparePEM(t)
	//defer os.RemoveAll(path)
	//// new host
	//node, err := net_libp2p.NewMonitorHost(testLocal, path, testPid, testDelimit)
	//require.NoError(t, err)
	//ch, err := node.Listen()
	//require.NoError(t, err)
	//go func() {
	//	for {
	//		x := <- ch
	//		t.Log("get data: ", x)
	//	}
	//}()
	//select {}
}

func TestChannelConnection(t *testing.T) {
	// prepare environment
	path := preparePEM(t)
	defer os.RemoveAll(path)
	// new host
	node, err := net_libp2p.NewMonitorHost(testLocal, path, testPid, testDelimit)
	require.NoError(t, err)
	_, err = node.Listen()
	require.NoError(t, err)
	// new host
	node, err = net_libp2p.NewDummyHost("/p2p/1.0.0", '\n')
	require.NoError(t, err)
	// decode address
	ctx := context.Background()
	addr, err := multiaddr.NewMultiaddr(localAddress)
	if err != nil {
		panic(err)
	}
	// send ping
	PingServer(ctx, node.Host, addr)
	// send message
	SendMsg(ctx, node.Host, addr)
}

func PingServer(ctx context.Context, node host.Host, addr multiaddr.Multiaddr) {
	p, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	if err := node.Connect(ctx, *p); err != nil {
		panic(err)
	}
	fmt.Println("sending 5 ping messages to", addr)
	ch := ping.Ping(ctx, node, p.ID)
	for i := 0; i < 5; i++ {
		res := <-ch
		fmt.Println("pinged", addr, "in", res.RTT)
	}
}

func SendMsg(ctx context.Context, node host.Host, addr multiaddr.Multiaddr) {
	p, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		panic(err)
	}
	if err := node.Connect(ctx, *p); err != nil {
		panic(err)
	}
	defer node.Close()
	fmt.Println("sending messages to", addr)
	s, err := node.NewStream(ctx, p.ID, "/p2p/1.0.0")
	if err != nil {
		return
	}
	_, err = s.Write([]byte("this is a test\n"))
	if err != nil {
		return
	}
}

func preparePEM(t *testing.T) string {
	var err error
	// 创建临时文件目录
	p, err := ioutil.TempDir(os.TempDir(), "mock")
	require.NoError(t, err)
	path := p + "/PEM"
	// 写入私钥
	err = ioutil.WriteFile(path, []byte(privatePEMStr), 0600)
	require.NoError(t, err)
	return path
}
