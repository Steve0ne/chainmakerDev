/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package chainmaker

import (
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/prover"
)

func MockNilChainMakerAdapter() *ChainMakerAdapter {
	config := []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			LogLevel:     logger.INFO,
			FilePath:     "logs/default.log",
			MaxAge:       365,
			RotationTime: 1,
			LogInConsole: false,
			ShowColor:    true,
		},
	}
	logger.InitLogConfig(config)
	log := logger.GetLogger(logger.ModuleAdapter)
	pd := prover.GetProverDispatcher()
	return &ChainMakerAdapter{
		chainID:    "chainID",
		dispatcher: pd,
		sdk:        nil,
		logger:     log,
	}
}
