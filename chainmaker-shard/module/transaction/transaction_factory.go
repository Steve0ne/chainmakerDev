/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package transaction

import (
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/store"
)

// InitManager init the Manager
func InitManager(stateDB store.StateDB) *Manager {
	manager := GetTransactionManager()
	manager.SetStateDB(stateDB)
	manager.SetLogger(logger.GetLogger(logger.ModuleTransactionMgr))
	return manager
}
