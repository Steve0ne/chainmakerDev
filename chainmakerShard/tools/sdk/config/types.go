/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"fmt"

	"chainmaker.org/chainmaker-cross/conf"
)

// Config Local config lists
type Config struct {
	ConfigLists `mapstructure:"configs"`
	Http        *conf.HttpTransport `mapstructure:"http"`
}

func (c *Config) Init() {
	c.Overload(configTemFileSet, chainParamNameSet)
}

func defaultHandle(cfg *Config) (*Config, error) {
	return cfg, nil
}

func (c *Config) Overload(Decors ...CfgDecorator) (*Config, error) {
	if len(Decors) == 0 {
		return c, nil
	}
	fn := handle(defaultHandle, Decors...)
	return fn(c)
}

func handle(fn CfgHandle, Decors ...CfgDecorator) CfgHandle {
	for _, d := range Decors {
		fn = d(fn)
	}
	return fn
}

// ConfigLists cross config array
type ConfigLists []*CrossChainConf

// ConvertToMap convert the list to Map
func (ccl *ConfigLists) ConvertToMap() *ConfigMap {
	ccM := make(ConfigMap)
	for _, v := range *ccl {
		switch v.ChainType {
		case "chainmaker":
			// TODO 删除临时文件
			_, path, err := createChainMakerSdkTmpFile()
			if err != nil {
				panic(fmt.Sprintf("create chainmaker sdk tmp file error: %s", err))
			}
			v.ChainConfigTemplatePath = path
		case "fabric":
			// TODO 删除临时文件
			_, path, err := createFabricSdkTmpFile()
			if err != nil {
				panic(fmt.Sprintf("create fabric sdk tmp file error: %s", err))
			}
			v.ChainConfigTemplatePath = path
		default:
			panic(fmt.Sprintf("unrecognized chain type: %s", v.ChainType))
		}

		// 设置默认的系统合约参数
		// 事物合约参数
		v.TransactionContractName = "TransactionStable"
		v.TransactionExecuteMethod = "Execute"
		v.TransactionCommitMethod = "Commit"
		v.TransactionRollbackMethod = "Rollback"

		v.TransactionExecuteDataKey = "executeData"
		v.TransactionRollbackDataKey = "rollbackData"

		v.BusinessCrossIDKey = "crossID"
		v.BusinessContractNameKey = "contractName"
		v.BusinessMethodKey = "method"
		v.BusinessParamsKey = "params"

		// 检查参数覆盖 TODO

		// 装载配置
		ccM[v.ChainID] = v
	}
	return &ccM
}

// ConfigMap Local config map
type ConfigMap map[string]*CrossChainConf

// GetCrossChainConfByChainID load cross chain config by chainID
func (ccm *ConfigMap) GetCrossChainConfByChainID(chainID string) *CrossChainConf {
	ccConf, ok := (*ccm)[chainID]
	if ok {
		return ccConf
	}
	return nil
}

// CrossChainConf the struct of CrossChainConf
type CrossChainConf struct {
	ChainID                 string `mapstructure:"chain_id"` // 跨链sdk chainID
	ChainType               string `mapstructure:"chain_type"`
	OrgID                   string `mapstructure:"org_id"`
	SignKeyPath             string `mapstructure:"sign_key_path"`
	SignCrtPath             string `mapstructure:"sign_crt_path"`
	ChainConfigTemplatePath string `mapstructure:"chain_config_template_path"`

	//ChainClientConfigPath      string `mapstructure:"chain_client_config_path"`      // 配置文件路径
	TransactionContractName    string                 `mapstructure:"transaction_contract_name"`     // 事物合约名，一般一条链复用一个事务合约，合约提供通用的 执行、确认、回滚 方法
	TransactionExecuteMethod   string                 `mapstructure:"transaction_execute_method"`    // 事物合约 执行方法 名
	TransactionCommitMethod    string                 `mapstructure:"transaction_commit_method"`     // 事物合约 确认方法 名
	TransactionRollbackMethod  string                 `mapstructure:"transaction_rollback_method"`   // 事物合约 回滚方法 名
	TransactionExecuteDataKey  string                 `mapstructure:"transaction_execute_data_key"`  // 调用事物合约执行方法，执行数据入参的键
	TransactionRollbackDataKey string                 `mapstructure:"transaction_rollback_data_key"` // 调用事物合约执行方法，回滚数据入参的键
	BusinessCrossIDKey         string                 `mapstructure:"business_cross_id_key"`         // 跨链交易ID的键
	BusinessProofKey           string                 `mapstructure:"business_proof_key"`            // 跨链交易proofKey的键
	BusinessContractNameKey    string                 `mapstructure:"business_contract_name_key"`    // 事务合约执行跨合约调用时，解析业务合约 合约名 的键
	BusinessMethodKey          string                 `mapstructure:"business_method_key"`           // 事务合约执行跨合约调用时，解析业务合约 合约方法 的键
	BusinessParamsKey          string                 `mapstructure:"business_params_key"`           // 事务合约执行跨合约调用时，解析业务合约 合约入参 的键
	ExtraParams                map[string]interface{} `mapstructure:"extra_params"`                  // 额外入参
}

func (c *CrossChainConf) GetExtraParamsByKey(key string) (interface{}, error) {
	value, ok := c.ExtraParams[key]
	if !ok {
		return nil, fmt.Errorf("can't fond key in extra params")
	}
	return value, nil
}
