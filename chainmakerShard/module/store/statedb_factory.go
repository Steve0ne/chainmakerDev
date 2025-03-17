/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"fmt"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/store/kvdb"
	"chainmaker.org/chainmaker-cross/store/kvdb/factory"
	storetypes "chainmaker.org/chainmaker-cross/store/types"
)

// InitStateDB init state database
func InitStateDB() StateDB {
	var stateDB StateDB
	dbProvider := storetypes.StateDBProvider(conf.Config.StorageConfig.Provider)
	switch dbProvider {
	case storetypes.LevelDB:
		dbProvider, err := factory.NewKvDBProvider(dbProvider, conf.Config.StorageConfig.LevelDB)
		if err == nil {
			stateDB = kvdb.NewKvStateDB(dbProvider)
		} else {
			panic(fmt.Sprintf("init statedb config failed, %v", err))
		}
	case storetypes.Memory:
		dbProvider, err := factory.NewKvDBProvider(dbProvider, nil)
		if err == nil {
			stateDB = kvdb.NewKvStateDB(dbProvider)
		} else {
			panic(fmt.Sprintf("init statedb config failed, %v", err))
		}
	default:
		dbProvider, err := factory.NewKvDBProvider(storetypes.Memory, nil)
		if err == nil {
			stateDB = kvdb.NewKvStateDB(dbProvider)
		} else {
			panic(fmt.Sprintf("init statedb config failed, %v", err))
		}
	}
	return stateDB
}
