/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package conf

import (
	"fmt"

	"github.com/pkg/errors"
)

type CfgHandle func(cfg *Config) (*Config, error)

type CfgDecorator func(CfgHandle) CfgHandle

func chainParamNameSet(fn CfgHandle) CfgHandle {
	return func(cfg *Config) (*Config, error) {
		for _, v := range cfg.ConfigLists {
			switch v.ChainType {
			case "chainmaker":
				v.TransactionContractName = "CROSS_TRANSACTION"
				v.TransactionExecuteMethod = "EXECUTE"
				v.TransactionCommitMethod = "COMMIT"
				v.TransactionRollbackMethod = "ROLLBACK"

				v.TransactionExecuteDataKey = "execData"
				v.TransactionRollbackDataKey = "rollbackData"

				v.BusinessCrossIDKey = "crossID"
				v.BusinessProofKey = "proofKey"
				v.BusinessContractNameKey = "contract"
				v.BusinessMethodKey = "method"
				v.BusinessParamsKey = "params"
			case "fabric":
				v.TransactionContractName = "Transaction"
				v.TransactionExecuteMethod = "Execute"
				v.TransactionCommitMethod = "Commit"
				v.TransactionRollbackMethod = "Rollback"

				v.TransactionExecuteDataKey = "executeData"
				v.TransactionRollbackDataKey = "rollbackData"

				v.BusinessCrossIDKey = "crossID"
				v.BusinessProofKey = "proofKey"
				v.BusinessContractNameKey = "contract"
				v.BusinessMethodKey = "method"
				v.BusinessParamsKey = "params"
			}
		}
		if fn == nil {
			return cfg, nil
		}
		return fn(cfg)
	}
}

func configTemFileSet(fn CfgHandle) CfgHandle {
	return func(cfg *Config) (*Config, error) {
		for _, v := range cfg.ConfigLists {
			switch v.ChainType {
			case "chainmaker":
				// TODO 删除临时文件
				_, path, err := createChainMakerSdkTmpFile()
				if err != nil {
					return cfg, errors.WithMessage(err, "create chainmaker sdk tmp file error:")
				}
				v.ChainConfigTemplatePath = path
			case "fabric":
				// TODO 删除临时文件
				_, path, err := createFabricSdkTmpFile()
				if err != nil {
					return cfg, errors.WithMessage(err, "create fabric sdk tmp file error:")
				}
				v.ChainConfigTemplatePath = path
			default:
				return cfg, fmt.Errorf("unrecognized chain type: %s", v.ChainType)
			}
		}
		if fn == nil {
			return cfg, nil
		}
		return fn(cfg)
	}
}
