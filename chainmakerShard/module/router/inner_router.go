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

var innerRouter *InnerRouter

func init() {
	innerRouter = &InnerRouter{
		chainIDs: make([]string, 0),
	}
}

// GetInnerRouter return the instance of inner router
func GetInnerRouter() *InnerRouter {
	return innerRouter
}

// InnerRouter is a router for inner message
type InnerRouter struct {
	chainIDs []string                     // 直连的chainID
	ch       *channel.InnerChannel        // 消息队列
	contexts *event.ProofResponseContexts // 处理结果的验证返回
}

// Init init inner router
func (i *InnerRouter) Init(chainIDs []string) {
	i.chainIDs = append(i.chainIDs, chainIDs...)
	i.ch = channel.GetInnerChannel()
	i.contexts = event.GetProofResponseContexts()
}

// GetChainIDs return chain ids
func (i *InnerRouter) GetChainIDs() []string {
	return i.chainIDs
}

// GetType return type of router
func (i *InnerRouter) GetType() RouterType {
	return InnerRouterType
}

// Invoke call chain call handler and wait response until completed
func (i *InnerRouter) Invoke(eve *eventproto.TransactionEvent, waitTime time.Duration) (*event.ProofResponse, error) {
	// 创建返回对象
	proofResponse := event.NewProofResponse(eve.GetCrossID(), eve.GetChainID(), eve.OpFunc)
	context := event.NewProofResponseContext(proofResponse)
	// 注册该上下文到集合中
	i.contexts.Register(context)
	// 将事务事件放入队列即可
	i.ch.Deliver(event.NewTransactionEventContext(context.GetKey(), eve))
	// 等待结果
	proofResponse.Wait(waitTime)
	return proofResponse, nil
}
