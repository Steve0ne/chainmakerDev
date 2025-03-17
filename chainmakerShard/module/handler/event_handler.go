/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package handler

import (
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/store"
	"go.uber.org/zap"
)

type HandlerType int32

const (
	ChainCall HandlerType = iota
	TransactionProcess
	CrossProcess
	CrossSearch
)

type EventHandler interface {

	// SetStateDB set state database
	SetStateDB(stateDB store.StateDB)

	// SetLogger set logger
	SetLogger(logger *zap.SugaredLogger)

	// GetType return the type of handler
	GetType() HandlerType

	// Handle handle event
	Handle(eve event.Event, syncWait bool) (interface{}, error)
}
