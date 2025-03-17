/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/store"
)

// InitEventHandlers init all event handlers
func InitEventHandlers(stateDB store.StateDB, eventChan chan event.Event) *EventHandlerTools {
	// 注册ChainCallHandler
	handlerTools.Register(getChainCallHandler(stateDB))
	// 注册TransactionProcessHandler
	handlerTools.Register(getTransactionProcessHandler(stateDB))
	// 注册CrossProcessHandler
	handlerTools.Register(getCrossProcessHandler(stateDB, eventChan))
	// 注册CrossSearchHandler
	handlerTools.Register(getCrossSearchHandler(stateDB))
	return handlerTools
}

// getChainCallHandler return chain call handler
func getChainCallHandler(stateDB store.StateDB) *ChainCallHandler {
	chainCallHandler := GetChainCallHandler()
	chainCallHandler.SetStateDB(stateDB)
	chainCallHandler.SetLogger(logger.GetLogger(logger.ModuleHandler))
	return chainCallHandler
}

// getTransactionProcessHandler return transaction process handler
func getTransactionProcessHandler(stateDB store.StateDB) *TransactionProcessHandler {
	txProcessHandler := GetTransactionProcessHandler()
	txProcessHandler.SetStateDB(stateDB)
	txProcessHandler.SetLogger(logger.GetLogger(logger.ModuleHandler))
	return transactionProcessHandler
}

// getCrossProcessHandler return cross process handler
func getCrossProcessHandler(stateDB store.StateDB, eventChan chan event.Event) *CrossProcessHandler {
	crossProcessHandler := GetCrossProcessHandler()
	crossProcessHandler.SetStateDB(stateDB)
	crossProcessHandler.SetLogger(logger.GetLogger(logger.ModuleHandler))
	crossProcessHandler.SetEventChan(eventChan)
	return crossProcessHandler
}

// getCrossSearchHandler return cross search handler
func getCrossSearchHandler(stateDB store.StateDB) *CrossSearchHandler {
	crossSearchHandler := GetCrossSearchHandler()
	crossSearchHandler.Init()
	crossSearchHandler.SetStateDB(stateDB)
	crossSearchHandler.SetLogger(logger.GetLogger(logger.ModuleHandler))
	return crossSearchHandler
}
