/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"chainmaker.org/chainmaker-cross/event"
)

const (
	InnerChannelLength int = 1024
)

var innerChannel *InnerChannel

func init() {
	innerChannel = &InnerChannel{
		ch: make(chan *event.TransactionEventContext, InnerChannelLength),
	}
}

// GetInnerChannel return instance of inner channel
func GetInnerChannel() *InnerChannel {
	return innerChannel
}

// InnerChannel inner channel
type InnerChannel struct {
	ch chan *event.TransactionEventContext // 消息通道
}

// GetChanType return type of channel
func (c *InnerChannel) GetChanType() TransmissionChanType {
	return InnerTransmissionChan
}

// Deliver write transaction event context to channel
func (c *InnerChannel) Deliver(eve *event.TransactionEventContext) error {
	c.ch <- eve
	return nil
}

// GetChan return object of channel
func (c *InnerChannel) GetChan() chan *event.TransactionEventContext {
	return c.ch
}
