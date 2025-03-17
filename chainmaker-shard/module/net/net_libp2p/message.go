/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net_libp2p

import (
	"time"

	"chainmaker.org/chainmaker-cross/net"
	"github.com/google/uuid"
	"github.com/libp2p/go-libp2p-core/peer"
)

var _ net.Message = (*LibP2pMessage)(nil)

// LibP2pMessage is message which wrapper libp2p protocol
type LibP2pMessage struct {
	//ClientVersion       string  `protobuf:"bytes,1,opt,name=clientVersion,proto3" json:"clientVersion,omitempty"`
	Timestamp int64   `protobuf:"varint,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"` // 时间戳
	ID        string  `protobuf:"bytes,3,opt,name=id,proto3" json:"id,omitempty"`                // 消息的UUID
	Gossip    bool    `protobuf:"varint,4,opt,name=gossip,proto3" json:"gossip,omitempty"`       // 是否使用 gossip 扩散协议，当前未启动
	NodeId    peer.ID `protobuf:"bytes,5,opt,name=nodeId,proto3" json:"nodeID,omitempty"`        // 消息发送者的ID
	//NodePubKey          []byte  `protobuf:"bytes,6,opt,name=nodePubKey,proto3" json:"nodePubKey,omitempty"`
	//Sign                []byte 	`protobuf:"bytes,7,opt,name=sign,proto3" json:"sign,omitempty"`
	Payload []byte `protobuf:"bytes,7,opt,name=data,proto3" json:"payload,omitempty"` // 传输数据主体
}

// NewLibP2pMessage create new message
func NewLibP2pMessage(nodeId string, payload []byte, gossip bool) (*LibP2pMessage, error) {
	peerId, err := peer.Decode(nodeId)
	if err != nil {
		return nil, err
	}
	return &LibP2pMessage{
		time.Now().Unix(),
		uuid.New().String(),
		gossip,
		peerId,
		payload,
	}, nil
}

// GetNodeID return node of message
func (m *LibP2pMessage) GetNodeID() string {
	return m.NodeId.Pretty()
}

// GetPayload return the payload of message
func (m *LibP2pMessage) GetPayload() []byte {
	return m.Payload
}
