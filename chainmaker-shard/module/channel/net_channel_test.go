/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/event/coder"
	"chainmaker.org/chainmaker-cross/net"
	"chainmaker.org/chainmaker-cross/net/net_libp2p"
	"chainmaker.org/chainmaker-cross/utils"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/stretchr/testify/require"
)

const privatePEMStr = "-----BEGIN EC PRIVATE KEY-----\n" +
	"MHcCAQEEIKwz2DjJZWtxnHNtMroZEjr28dqihVAPcyrmQp/5XeGRoAoGCCqGSM49\n" +
	"AwEHoUQDQgAERrAT6tdtd95gd6zh2/qics3s/h7eUtAGLdPhuGstl1j9iqSXI/0i\n" +
	"9qyw+HHd2IkkLdQX6ivmaovGS0vbSDpRCA==\n" +
	"-----END EC PRIVATE KEY-----\n"
const localAddress = "/ip4/127.0.0.1/tcp/12345/p2p/QmSVuE22fYPeTYg4iUJJbVmwtaumYavSApR3PSm4Gb6upV"
const testLocal = "/ip4/127.0.0.1/tcp/12345"
const testPid = protocol.ID("/p2p/1.0.0")
const testDelimit = '\n'

func TestNetChannel_Init(t *testing.T) {
	// mock host
	node, err := mockHost(t)
	require.NoError(t, err)
	require.NotNil(t, node)
	_, err = node.Listen()
	require.NoError(t, err)

	//mock connection
	connection, err := net_libp2p.NewLibP2pConnection(localAddress, testPid, testDelimit, 10, 1000)
	require.NoError(t, err)
	require.NotNil(t, connection)

	// test get channel
	nc := NewNetChannel(connection)
	require.NotNil(t, nc)

	err = nc.Init()
	require.Nil(t, err)
}

func TestNetChannelInitReadData(t *testing.T) {
	// mock host
	node, err := mockHost(t)
	require.NoError(t, err)
	require.NotNil(t, node)
	_, err = node.Listen()
	require.NoError(t, err)

	//mock connection
	connection, err := net_libp2p.NewLibP2pConnection(localAddress, testPid, testDelimit, 10, 1000)
	require.NoError(t, err)
	require.NotNil(t, connection)

	// test get channel
	n := NewNetChannel(connection)
	require.NotNil(t, n)

	dataChan, err := n.connection.ReadData()
	require.Nil(t, err)
	go func() {

	}()
	// 启动监听，用于读取数据并打印，不做启动事情
	go func() {
		for {
			select {
			case msg := <-dataChan:
				// 从通道中读到数据
				if len(msg.GetPayload()) < MinDataLength {
					// 打印错误信息
					n.log.Error("receive data is illegal")
				} else {
					n.log.Debugf("receive data length = %v", len(msg.GetPayload()))
					receivedData, err := utils.Base64DecodeToBytes(string(msg.GetPayload()))
					if err != nil {
						n.log.Error("base64 decode data failed, ", err)
						continue
					}
					eventTy, marshalTy := event.EventType(receivedData[coder.EventTyIndex]), event.MarshalType(receivedData[coder.MarshalTyIndex])
					if eventTy == event.ProofRespEventType {
						if eveCoder, exist := n.coders.GetDefaultCoder(eventTy); exist {
							// 处理
							if marshalTy == event.BinaryMarshalType {
								eve, err := eveCoder.UnmarshalFromBinary(receivedData)
								if err != nil {
									// 打印错误信息
									n.log.Error("unmarshal receive data failed, ", err)
								} else {
									if resp, ok := eve.(*event.ProofResponse); ok {
										// 填充结果
										if resp.Code == event.SuccessResp {
											n.log.Infof("cross[%s]->chain[%s]->key[%s] response is success",
												resp.CrossID, resp.ChainID, resp.Key)
											// 操作成功，填充结果
											n.contexts.DoneByProofResp(resp)
										} else {
											n.log.Errorf("cross[%s]->chain[%s]->key[%s] response is failed",
												resp.CrossID, resp.ChainID, resp.Key)
											n.contexts.DoneError(resp.Key, resp.Msg)
										}
									}
								}
							}
						}
					} else {
						n.log.Errorf("")
					}
				}
			case <-time.After(conf.LogWritePeriod):
				// 打印日志，表明在正常活着
				n.log.Info("net channel is running periodically!")
			}
		}
	}()
}

func TestNetChannel(t *testing.T) {
	// mock host
	node, err := mockHost(t)
	require.NoError(t, err)
	require.NotNil(t, node)
	_, err = node.Listen()
	require.NoError(t, err)

	//mock connection
	connection, err := net_libp2p.NewLibP2pConnection(localAddress, testPid, testDelimit, 10, 1000)
	require.NoError(t, err)
	require.NotNil(t, connection)

	// test get channel
	nc := NewNetChannel(connection)
	require.NotNil(t, nc)

	// test get channel type
	cType := nc.GetChanType()
	require.Equal(t, cType, NetTransmissionChan)

	// test deliver
	err = nc.Deliver(&event.TransactionEventContext{Key: "test", Event: &event.TransactionEvent{}})
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
	node, err := net_libp2p.NewMonitorHost(testLocal, path, testPid, testDelimit)
	require.NoError(t, err)
	t.Log("local libp2p peer id: ", node.Host.ID().Pretty())
	return node, nil
}
