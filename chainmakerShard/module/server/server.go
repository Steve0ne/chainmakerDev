/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	"errors"
	"sync"

	"chainmaker.org/chainmaker-cross/adapter"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/handler"
	"chainmaker.org/chainmaker-cross/listener"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/prover"
	"chainmaker.org/chainmaker-cross/router"
	"chainmaker.org/chainmaker-cross/store"
	"chainmaker.org/chainmaker-cross/transaction"
	"go.uber.org/zap"
)

// Server is the cross chain server
type Server struct {
	sync.Mutex                                        // lock
	started           bool                            // 是否启动的标志
	logger            *zap.SugaredLogger              // log
	stateDB           store.StateDB                   // 存储
	transactionMgr    *transaction.Manager            // 跨链消息事物管理模块
	listenerMgr       *listener.Manager               // 监听服务
	routerDispatcher  *router.RouterDispatcher        // 路由管理服务
	proverDispatcher  *prover.ProverDispatcher        // 验证管理服务
	adapterDispatcher *adapter.ChainAdapterDispatcher // 转接器管理服务
	eventHandlers     *handler.EventHandlerTools      // 跨链事件消息处理函数
}

// NewServer create new cross chain server
func NewServer() *Server {
	stateDB := store.InitStateDB()
	transactionMgr := transaction.InitManager(stateDB)
	adapterDispatcher := adapter.InitAdapters()
	event.InitLog(logger.GetLogger(logger.ModuleDefault))
	return &Server{
		started:           false,
		logger:            logger.GetLogger(logger.ModuleServer),
		stateDB:           stateDB,
		transactionMgr:    transactionMgr,
		listenerMgr:       listener.InitListener(),
		routerDispatcher:  router.InitRouters(adapterDispatcher.GetChainIDs()),
		proverDispatcher:  prover.InitProvers(),
		adapterDispatcher: adapterDispatcher,
		eventHandlers:     handler.InitEventHandlers(stateDB, transactionMgr.GetEventChan()),
	}
}

// GetTransactionMgr return TransactionMgr
func (s *Server) GetTransactionMgr() *transaction.Manager {
	return s.transactionMgr
}

// GetListenerMgr return ListenerMgr
func (s *Server) GetListenerMgr() *listener.Manager {
	return s.listenerMgr
}

// GetStateDB return state database
func (s *Server) GetStateDB() store.StateDB {
	return s.stateDB
}

// GetRouterDispatcher return RouterDispatcher
func (s *Server) GetRouterDispatcher() *router.RouterDispatcher {
	return s.routerDispatcher
}

// GetProverDispatcher return ProverDispatcher
func (s *Server) GetProverDispatcher() *prover.ProverDispatcher {
	return s.proverDispatcher
}

// GetAdapterDispatcher return ChainAdapterDispatcher
func (s *Server) GetAdapterDispatcher() *adapter.ChainAdapterDispatcher {
	return s.adapterDispatcher
}

// GetEventHandlers return EventHandlerTools
func (s *Server) GetEventHandlers() *handler.EventHandlerTools {
	return s.eventHandlers
}

// Start cross chain server start
func (s *Server) Start() error {
	var log = logger.GetLogger(logger.ModuleServer)
	s.Lock()
	defer s.Unlock()
	if s.started {
		log.Error("this server has been started")
		return errors.New("this server has been started")
	}
	log.Info("--- start transaction manager ---")
	err := s.transactionMgr.Start()
	if err != nil {
		log.Error("start transaction manager error:", err)
		return err
	}
	log.Info("--- start transaction manager over ---")
	log.Info("--- start listener manager ---")
	err = s.listenerMgr.Start()
	if err != nil {
		log.Error("start listener manager error:", err)
		return err
	}
	log.Info("--- start listener manager over ---")
	s.beenStarted()
	return nil
}

// Stop stop the cross chain server
func (s *Server) Stop() error {
	s.Lock()
	defer s.Unlock()
	if !s.started {
		return errors.New("this server has not been started")
	}
	s.transactionMgr.Stop()
	s.stateDB.Close()
	if err := s.listenerMgr.Stop(); err != nil {
		// 打印err
		s.logger.Errorf("stop proxy server error", err)
	}
	s.beenStopped()
	return nil
}

func (s *Server) beenStarted() {
	s.started = true
}

func (s *Server) beenStopped() {
	s.started = false
}
