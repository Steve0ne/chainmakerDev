/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net_libp2p

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"reflect"

	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/net"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	ma "github.com/multiformats/go-multiaddr"
	"go.uber.org/zap"
)

const (
	DefaultAddress     = "/ip4/0.0.0.0/tcp/0"
	WriteChannelLength = 1024
	ReadChannelLength  = 1024
)

// LibP2pNode is node which wrapper libp2p
type LibP2pNode struct {
	host.Host                    // p2p网络中的节点概念
	readChan  chan net.Message   // 读通道
	writeChan chan net.Message   // 写通道
	pid       protocol.ID        // 通信的协议号
	delimit   byte               // 消息分割符
	log       *zap.SugaredLogger // log
}

// ID return the id of libp2p node
func (l *LibP2pNode) ID() string {
	return l.Host.ID().Pretty()
}

// Listen node server start
func (l *LibP2pNode) Listen() (chan net.Message, error) {
	var err error
	if l.readChan != nil {
		return l.readChan, nil
	}
	ch := make(chan net.Message, ReadChannelLength)
	l.SetStreamHandler(l.pid, func(s network.Stream) {
		// Create a buffer stream for non blocking read and write.
		// stream 's' will stay open until you close it.
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		go func() {
			err = readData(rw, ch, l.delimit)
		}()
	})
	go func() {
		writeData(l)
	}()
	return ch, err
}

// Write write message to node
func (l *LibP2pNode) Write(msg net.Message) error {
	l.log.Debug("node write data")
	//m, ok := msg.(*LibP2pMessage)
	//if !ok {
	//	l.log.Error("unsupported libp2p message type", reflect.TypeOf(msg))
	//}
	l.writeChan <- msg
	return nil
}

// Stop node server stop
func (l *LibP2pNode) Stop() error {
	err := l.Close()
	if err != nil {
		l.log.Error(err)
		return err
	}
	return nil
}

// NewDummyHost create new dummy host node
func NewDummyHost(pid protocol.ID, delimit byte) (*LibP2pNode, error) {
	ctx := context.Background()
	node, err := libp2p.New(ctx,
		libp2p.ListenAddrStrings(DefaultAddress),
	)
	if err != nil {
		return nil, err
	}
	readChan := make(chan net.Message, WriteChannelLength)
	writeChan := make(chan net.Message, WriteChannelLength)
	return &LibP2pNode{
		node, readChan, writeChan, pid, delimit, logger.GetLogger(logger.ModuleNet)}, nil
}

// NewMonitorHost creates a LibP2P host with known peer ID listening on the given address.
func NewMonitorHost(address, privKeyFile string, pid protocol.ID, delimit byte) (*LibP2pNode, error) {
	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, err := prepareKey(privKeyFile)
	if err != nil {
		return nil, err
	}

	// set up options
	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(address),
		libp2p.Identity(priv),
	}

	// TODO support secret connection

	// create new host
	basicHost, err := libp2p.New(context.Background(), opts...) // TODO, listen on repeat port
	if err != nil {
		return nil, err
	}

	// configure default ping protocol
	pingService := &ping.PingService{Host: basicHost}
	basicHost.SetStreamHandler(ping.ID, pingService.PingHandler)

	multiAddr, err := hostToMultiAddr(basicHost)
	if err != nil {
		return nil, err
	}
	log := logger.GetLogger(logger.ModuleNet)
	if multiAddr != nil {
		log.Info(fmt.Sprintf("channel start listen %s\n", multiAddr))
	}
	readChan := make(chan net.Message, WriteChannelLength)
	writeChan := make(chan net.Message, WriteChannelLength)
	return &LibP2pNode{basicHost, readChan, writeChan, pid, delimit, log}, nil
}

// prepareKey read private key from listener config
func prepareKey(keyFile string) (crypto.PrivKey, error) {
	keyBytes, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}
	privateKey, err := PrivateKeyFromPEM(keyBytes, nil)
	return privateKey, err
}

// hostToMultiAddr build host multiaddress
func hostToMultiAddr(host host.Host) (ma.Multiaddr, error) {
	log := logger.GetLogger(logger.ModuleNet)
	hostAddr, err := ma.NewMultiaddr(fmt.Sprintf("/p2p/%s", host.ID().Pretty()))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := host.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)

	return fullAddr, nil
}

// readData read channel listener stream inputs, dispatch every event to handler
func readData(rw *bufio.ReadWriter, ch chan net.Message, delimit byte) error {
	var err error
	log := logger.GetLogger(logger.ModuleNet)
	handle := NewLibP2PSteam(rw, delimit)
	err = handle.ReadStream(ch)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// writeData write channel listen stream write inputs, write data back to stream
func writeData(l *LibP2pNode) {
	//var err error
	for {
		// read channel
		msg, ok := <-l.writeChan
		if !ok {
			l.log.Error("load date from write-channel failed")
			return
		}
		// create new stream
		s, err := l.NewStream(context.Background(), peer.ID(msg.GetNodeID()), l.pid)
		if err != nil {
			l.log.Error("create stream error: ", err)
			continue
		}
		rw := bufio.NewReadWriter(nil, bufio.NewWriter(s))
		handle := NewLibP2PSteam(rw, l.delimit)
		m, ok := msg.(*LibP2pMessage)
		if !ok {
			l.log.Error("unsupported libp2p message type", reflect.TypeOf(msg))
			continue
		}
		err = handle.WriteStream(m)
		if err != nil {
			l.log.Error("write stream error: ", err)
			s.Reset()
			continue
		}
		l.log.Debug("new data write to stream")

		// close every write connection
		s.Close()
	}
}
