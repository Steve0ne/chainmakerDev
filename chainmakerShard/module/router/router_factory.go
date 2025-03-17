/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package router

import (
	"chainmaker.org/chainmaker-cross/channel"
	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/net"
	"chainmaker.org/chainmaker-cross/net/net_http"
	"chainmaker.org/chainmaker-cross/net/net_libp2p"
	"github.com/libp2p/go-libp2p-core/protocol"
)

// InitRouters init all the routers and return router dispatcher
func InitRouters(innerChainIDs []string) *RouterDispatcher {
	log := logger.GetLogger(logger.ModuleRouter)
	innerRouter := GetInnerRouter()
	// 初始化innerRouter
	innerRouter.Init(innerChainIDs)
	dispatcher.SetLogger(log)
	_ = dispatcher.Register(innerRouter)
	// 开始处理所有的Router
	for _, routerConfig := range conf.Config.RouterConfigs {
		routerProvider := net.ConnectionProvider(routerConfig.Provider)
		var connection net.Connection
		var err error
		if routerProvider == net.LibP2PConnection {
			if routerConfig.LibP2PRouter != nil {
				connection, err = net_libp2p.NewLibP2pConnection(
					routerConfig.LibP2PRouter.Address,
					protocol.ID(routerConfig.LibP2PRouter.ProtocolID),
					routerConfig.LibP2PRouter.GetDelimit(),
					routerConfig.LibP2PRouter.ReconnectLimit,
					routerConfig.LibP2PRouter.ReconnectInterval,
				)
				if err != nil {
					log.Errorf("connect [%s] failed", routerConfig.LibP2PRouter.Address)
				} else {
					log.Infof("connect [%s] established", routerConfig.LibP2PRouter.Address)
				}
			}
		} else if routerProvider == net.HttpConnection { //http连接
			connection, err = net_http.NewConnection(routerConfig.HttpRouter, net_http.WithRetryStrategy(net_http.RetryStrategy{
				MaxRetries: routerConfig.HttpRouter.RequestMaxRetries, Interval: routerConfig.HttpRouter.RequestRetryInterval,
			}))
		}
		if err == nil {
			// 将connection加入router
			netChannel := channel.NewNetChannel(connection)
			err := netChannel.Init()
			if err == nil {
				channelRouter := NewChannelRouter(routerConfig.GetChainIDs(), netChannel)
				err := GetDispatcher().Register(channelRouter)
				if err != nil {
					// 打印，但不处理
					log.Warn("register channel router failed, ", err)
				}
			} else {
				// 打印信息
				log.Warn("init channel router failed, ", err)
			}
		} else {
			// 打印信息
			log.Warn("create net connection failed, ", err)
		}
	}
	return dispatcher
}
