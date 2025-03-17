/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net_libp2p

import (
	"bufio"
	"context"
	"time"

	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/net"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/backoff"
	"github.com/Rican7/retry/strategy"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	ma "github.com/multiformats/go-multiaddr"
)

// LibP2pConnection is connection by the way of libp2p
type LibP2pConnection struct {
	pid     protocol.ID // 协议号，由配置文件定义
	address string      // 连接的目标地址
	host    *LibP2pNode // 连接的本地节点
	peerID  peer.ID     // 连接对端的网络身份ID
	//stream 	network.Stream		// 读写IO流
	//rw 		*bufio.ReadWriter	// 读写IO流
	delimit           byte // 数据分割符
	reconnectLimit    int  // 连接断开重连次数
	reconnectInterval int  // 重连间隔，单位毫秒
}

// NewLibP2pConnection create new libp2p connection
func NewLibP2pConnection(address string, pid protocol.ID, delimit byte, reconnectLimit, reconnectInterval int) (*LibP2pConnection, error) {
	log := logger.GetLogger(logger.ModuleNet)
	//ctx := context.Background()
	node, err := NewDummyHost(pid, delimit)
	if err != nil {
		return nil, err
	}
	host := node
	maddr, err := ma.NewMultiaddr(address)
	if err != nil {
		return nil, err
	}
	addr, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, err
	}
	host.Peerstore().AddAddrs(addr.ID, addr.Addrs, peerstore.PermanentAddrTTL)
	//stream, err := host.NewStream(ctx, addr.ID, pid)
	//if err != nil {
	//	return nil, err
	//}
	//rw := bufio.NewReadWriter(bufio.NewReaderSize(stream, 1024*1024), bufio.NewWriterSize(stream, 1024*1024))
	// pack connection object
	c := &LibP2pConnection{
		pid,
		address,
		host,
		addr.ID,
		//stream,
		delimit,
		reconnectLimit,
		reconnectInterval,
	}
	// disconnection handler
	var disconnectHandler = func(network.Network, network.Conn) {
		err = retry.Retry(func(uint) error {
			err := c.ResetLibP2pConnection()
			if err != nil {
				log.Error("reconnect error: ", err)
				return err
			}
			return nil
		},
			strategy.Limit(uint(reconnectLimit)),
			strategy.Backoff(backoff.Linear(time.Duration(reconnectInterval)*time.Millisecond)),
		)
		if err != nil {
			log.Error("reconnect error: ", err)
		} else {
			log.Info("peer reconnect")
		}
	}

	// 处理断线重连
	host.Host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: disconnectHandler,
	})

	return c, nil
}

// ResetLibP2pConnection reconnction when the net been reset
func (c *LibP2pConnection) ResetLibP2pConnection() error {
	ctx := context.Background()
	maddr, err := ma.NewMultiaddr(c.address)
	if err != nil {
		return err
	}
	addr, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}
	c.host.Peerstore().AddAddrs(addr.ID, addr.Addrs, peerstore.PermanentAddrTTL)
	err = c.host.Connect(ctx, *addr)
	if err != nil {
		return err
	}
	//stream, err := c.host.NewStream(ctx, addr.ID, c.pid)
	//if err != nil {
	//	return nil
	//}
	//rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
	//c.rw = rw
	//c.stream = stream
	return nil
}

// ReadData read data from the connection
func (c *LibP2pConnection) ReadData() (chan net.Message, error) {
	ch, err := c.host.Listen()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

// WriteData write the data to connection
func (c *LibP2pConnection) WriteData(msg net.Message) error {
	ctx := context.Background()
	stream, err := c.host.NewStream(ctx, c.peerID, c.pid)
	if err != nil {
		return nil
	}
	rw := bufio.NewReadWriter(bufio.NewReaderSize(stream, 1024*1024), bufio.NewWriterSize(stream, 1024*1024))
	s := NewLibP2PSteam(rw, c.delimit)
	m := msg.(*LibP2pMessage)
	err = s.WriteStream(m)
	if err != nil {
		return err
	}
	return nil
}

// PeerID return the peerID of remote node
func (c *LibP2pConnection) PeerID() string {
	return c.host.ID()
}

// Close close the connection
func (c *LibP2pConnection) Close() error {
	return c.host.Close()
}

func (c *LibP2pConnection) GetProvider() net.ConnectionProvider {
	return net.LibP2PConnection
}
