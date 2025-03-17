/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package listener

import (
	"chainmaker.org/chainmaker-cross/handler"
	"chainmaker.org/chainmaker-cross/listener/channel_listener"
	"chainmaker.org/chainmaker-cross/listener/inner_listener"
	"chainmaker.org/chainmaker-cross/listener/web_listener"
)

// Listener is listener
type Listener interface {

	// ListenStart start server
	ListenStart() error

	// Stop stop server
	Stop() error
}

var manager *Manager

func init() {
	manager = &Manager{
		dispatcher: handler.GetEventHandlerTools(),
		listeners:  make([]Listener, 0),
	}
}

// InitListener init all the listeners
func InitListener() *Manager {
	manager.InitListeners()
	return manager
}

// Manager is manager which will dispatch message
type Manager struct {
	dispatcher *handler.EventHandlerTools // 负责跨链消息的分发
	listeners  []Listener                 // 封装本地监听服务
}

// InitListeners init all the listeners
func (m *Manager) InitListeners() {
	cl := channel_listener.NewChannelListener()
	il := inner_listener.NewInnerListener()
	wl := web_listener.NewWebListener()
	m.listeners = append(m.listeners, cl, il, wl)
}

// Start all the listener start
func (m *Manager) Start() error {
	for _, l := range m.listeners {
		err := l.ListenStart()
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop stop all the listeners
func (m *Manager) Stop() error {
	for _, l := range m.listeners {
		if err := l.Stop(); err != nil {
			return err
		}
	}
	return nil
}
