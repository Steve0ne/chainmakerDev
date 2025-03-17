/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/channel"
	"chainmaker.org/chainmaker-cross/event"
)

// ChannelRouter is router which will communication with other cross chain proxy
type ChannelRouter struct {
	chainIDs []string                     // 转发代理支持的 chainID
	ch       *channel.NetChannel          // 跨链代理之间的连接
	contexts *event.ProofResponseContexts // 交易验证数据
}

// NewChannelRouter create new channel router
func NewChannelRouter(chainIDs []string, netCh *channel.NetChannel) *ChannelRouter {
	return &ChannelRouter{
		chainIDs: chainIDs,
		ch:       netCh,
		contexts: event.GetProofResponseContexts(),
	}
}

// GetType return the type of router
func (c *ChannelRouter) GetType() RouterType {
	return ChannelRouterType
}

// GetChainIDs return the chain ids
func (c *ChannelRouter) GetChainIDs() []string {
	return c.chainIDs
}

// Invoke put the event into channel, and wait response until completed
func (c *ChannelRouter) Invoke(eve *eventproto.TransactionEvent, waitTime time.Duration) (*event.ProofResponse, error) {
	// 创建返回对象
	proofResponse := event.NewProofResponse(eve.GetCrossID(), eve.GetChainID(), eve.OpFunc)
	context := event.NewProofResponseContext(proofResponse)
	// 注册该上下文到集合中
	c.contexts.Register(context)
	// 将事务事件放入队列
	err := c.ch.Deliver(event.NewTransactionEventContext(context.GetKey(), eve))
	if err != nil {
		// 从缓存中移除，防止内存膨胀
		c.contexts.Remove(context.GetKey())
		return proofResponse, err
	}
	// 等待结果
	proofResponse.Wait(waitTime)
	return proofResponse, nil
}
