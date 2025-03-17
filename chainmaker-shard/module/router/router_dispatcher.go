/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"errors"
	"fmt"
	"sync"
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
	"go.uber.org/zap"
)

var dispatcher *RouterDispatcher

func init() {
	dispatcher = &RouterDispatcher{
		routers: make(map[string]*Routers),
	}
}

// RouterDispatcher dispatcher of router
type RouterDispatcher struct {
	sync.RWMutex                     // lock
	routers      map[string]*Routers // 能够连接到指定 chainID 的节点路由
	logger       *zap.SugaredLogger  // log
}

// GetDispatcher return the instance of RouterDispatcher
func GetDispatcher() *RouterDispatcher {
	return dispatcher
}

// SetLogger set logger
func (d *RouterDispatcher) SetLogger(logger *zap.SugaredLogger) {
	d.logger = logger
}

// Register add router to router dispatcher
func (d *RouterDispatcher) Register(router Router) error {
	d.Lock()
	defer d.Unlock()
	chainIDs := router.GetChainIDs()
	if chainIDs == nil {
		return errors.New("chainIDs is empty")
	}
	for _, chainID := range chainIDs {
		if rs, exist := d.routers[chainID]; exist {
			// 已经存在该chainID的处理
			if err := rs.Add(router); err != nil {
				d.logger.Warnf("register chain[%v] router failed, ", chainID, err)
			}
		} else {
			// 不存在则创建
			routers := NewRouters()
			d.routers[chainID] = routers
			if err := routers.Add(router); err != nil {
				d.logger.Warnf("register chain[%v] router failed, ", chainID, err)
			}
		}
	}
	return nil
}

// Invoke is entry for handle of transaction event
func (d *RouterDispatcher) Invoke(eve *eventproto.TransactionEvent, waitTime time.Duration) (*event.ProofResponse, error) {
	var (
		router Router
		exist  bool
	)
	chainID := eve.GetChainID()
	// 判断是否支持本地路由
	if router, exist = d.getInnerRouter(chainID); !exist {
		// 判断是否支持远端路由
		if router, exist = d.getChannelRouter(chainID); !exist {
			d.logger.Errorf("can not find router to handle chain[%s]", chainID)
			return nil, fmt.Errorf("can not find router to handle chain[%v]", chainID)
		} else {
			d.logger.Infof("find channel router to handle event for chain[%s]", chainID)
		}
	} else {
		d.logger.Infof("find inner router to handle event for chain[%s]", chainID)
	}
	return router.Invoke(eve, waitTime)
}

func (d *RouterDispatcher) getChannelRouter(chainID string) (*ChannelRouter, bool) {
	d.RLock()
	defer d.RUnlock()
	if routers, exist := d.routers[chainID]; exist {
		r, exi := routers.GetChannelRouter()
		if exi {
			if chr, ok := r.(*ChannelRouter); ok {
				return chr, true
			}
		}
	}
	return nil, false
}

func (d *RouterDispatcher) getInnerRouter(chainID string) (*InnerRouter, bool) {
	d.RLock()
	defer d.RUnlock()
	if routers, exist := d.routers[chainID]; exist {
		r, exi := routers.GetInnerRouter()
		if exi {
			if chr, ok := r.(*InnerRouter); ok {
				return chr, true
			}
		}
	}
	return nil, false
}
