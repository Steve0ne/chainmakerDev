/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"fmt"
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
)

type RouterType int

const (
	InnerRouterType RouterType = iota
	ChannelRouterType
	PreferenceIndex = 0
)

// Router is interface for router
type Router interface {

	// GetType return type of router
	GetType() RouterType

	// GetChainIDs return chain's array which be support by this router
	GetChainIDs() []string

	// Invoke invoke the transaction event
	Invoke(eve *eventproto.TransactionEvent, waitTime time.Duration) (*event.ProofResponse, error)
}

// Routers is the entry for router
type Routers struct {
	inner Router   // 直连消息通道
	rs    []Router // 转发消息通道
}

// NewRouters create new routers
func NewRouters() *Routers {
	return &Routers{
		rs: make([]Router, 0),
	}
}

// Add add router to routers
func (rs *Routers) Add(router Router) error {
	ty := router.GetType()
	if ty == InnerRouterType {
		rs.inner = router
		return nil
	} else if ty == ChannelRouterType {
		rs.rs = append(rs.rs, router)
		return nil
	}
	return fmt.Errorf("can not support router type -> [%v]", ty)
}

// Support return true if inner router is not nil or channel router is not empty
func (rs *Routers) Support() bool {
	return rs.InnerSupport() || rs.ChannelSupport()
}

// InnerSupport return true if inner router is not nil
func (rs *Routers) InnerSupport() bool {
	return rs.inner != nil
}

// ChannelSupport return true if channel router is not empty
func (rs *Routers) ChannelSupport() bool {
	return rs.rs != nil
}

// GetInnerRouter return inner router
func (rs *Routers) GetInnerRouter() (Router, bool) {
	if rs.InnerSupport() {
		return rs.inner, true
	}
	return nil, false
}

// GetChannelRouter return channel router
func (rs *Routers) GetChannelRouter() (Router, bool) {
	if rs.ChannelSupport() {
		// 优先选择第一个
		return rs.rs[PreferenceIndex], true
	}
	return nil, false
}
