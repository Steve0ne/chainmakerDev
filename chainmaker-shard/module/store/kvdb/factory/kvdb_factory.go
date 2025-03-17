/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package factory

import (
	"errors"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/store/kvdb/leveldb"
	"chainmaker.org/chainmaker-cross/store/kvdb/memory"
	"chainmaker.org/chainmaker-cross/store/kvdb/types"
	storetypes "chainmaker.org/chainmaker-cross/store/types"
)

var (
	levelDBConfigError       = errors.New("create leveldb error by config")
	unsupportedProviderError = errors.New("can not support this provider")
)

// NewKvDBProvider create new kvdb provider instance by provider type
func NewKvDBProvider(provider storetypes.StateDBProvider, config interface{}) (types.KvDBProvider, error) {
	if provider == storetypes.LevelDB {
		if dbConf, ok := config.(*conf.LevelDBConfig); ok {
			return leveldb.NewLevelDBProvider(dbConf), nil
		} else {
			return nil, levelDBConfigError
		}
	} else if provider == storetypes.Memory {
		return memory.NewMemProvider(), nil
	}
	return nil, unsupportedProviderError
}
