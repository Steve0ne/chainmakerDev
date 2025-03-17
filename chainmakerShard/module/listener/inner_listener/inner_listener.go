/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package inner_listener

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"chainmaker.org/chainmaker-cross/channel"
	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/handler"
	"chainmaker.org/chainmaker-cross/logger"
	"go.uber.org/zap"
)

// InnerListener is listener to handle chain call event
type InnerListener struct {
	sync.Mutex                                // lock
	exchange     *channel.InnerChannel        // 消息队列
	logger       *zap.SugaredLogger           // log
	eventHandler handler.EventHandler         // 消息处理逻辑的接口
	contexts     *event.ProofResponseContexts // 消息证明
	cancel       context.CancelFunc           // 退出函数
	started      bool                         // 是否启动的tag
}

// NewInnerListener create new inner listener
func NewInnerListener() *InnerListener {
	return &InnerListener{
		exchange: channel.GetInnerChannel(),
		logger:   logger.GetLogger(logger.ModuleInnerListener),
		contexts: event.GetProofResponseContexts(),
	}
}

// ListenStart inner listener server start
func (l *InnerListener) ListenStart() error {
	l.Lock()
	defer l.Unlock()
	if l.started {
		return errors.New("this inner listener has been started")
	}
	if l.eventHandler == nil {
		eveHandler, exist := handler.GetEventHandlerTools().GetHandler(handler.ChainCall)
		if !exist {
			return errors.New("can not find handler to hand this event")
		}
		l.eventHandler = eveHandler
	}
	ctx, cancel := context.WithCancel(context.Background())
	l.cancel = cancel
	// 启动监听
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				l.logger.Warn("Module inner-listener stopped")
				return
			case <-time.After(conf.LogWritePeriod):
				// 打印日志
				l.logger.Info("inner listener is running periodically!")
			case eve := <-l.exchange.GetChan():
				// 获取到事件
				l.logger.Info("inner listener receive event from channel, will handle it")
				// 启动一个独立协程进行处理
				go l.Handle(eve)
			}
		}
	}(ctx)
	l.started = true
	return nil
}

// Stop inner listener server stop
func (l *InnerListener) Stop() error {
	if !l.started {
		return fmt.Errorf("this inner listener has not started")
	}
	l.cancel()
	l.started = false
	return nil
}

// GetExchange return the channel of exchange data
func (l *InnerListener) GetExchange() chan *event.TransactionEventContext {
	return l.exchange.GetChan()
}

// Handle handle the event which will invoke adapter
func (l *InnerListener) Handle(eveCtx *event.TransactionEventContext) {
	key := eveCtx.GetKey()
	_, err := l.eventHandler.Handle(eveCtx, true)
	if err != nil {
		l.contexts.DoneError(key, err.Error())
	}
	l.contexts.Remove(key)
}
