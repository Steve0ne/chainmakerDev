/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package adapter

import (
	"fmt"

	"chainmaker.org/chainmaker-cross/adapter/chainmaker"
	"chainmaker.org/chainmaker-cross/adapter/fabric"
	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/logger"
	"go.uber.org/zap"
)

type Provider string

const (
	ChainMakerProvider 	Provider = "chainmaker"
	FabricProvider 		Provider = "fabric"
	EthProvider        	Provider = "ETH"
)


// InitAdapters read adapter config and init instance of dispatcher
func InitAdapters() *ChainAdapterDispatcher {
	var log = logger.GetLogger(logger.ModuleAdapter)
	dispatcher.SetLog(log)
	log.Infof("get [%d] adapters to connect...", len(conf.Config.AdapterConfigs))
	for _, config := range conf.Config.AdapterConfigs {
		innerConfig := config
		adapterCfgPath := finalAdapterCfgPath(innerConfig)
		if adapterImpl, err := createAdapter(innerConfig, adapterCfgPath, log); err == nil {
			dispatcher.Register(adapterImpl)
		} else {
			log.Errorf("adapter config [%s] create adapter error:[%v]", adapterCfgPath, err)
		}
	}
	return dispatcher
}

// createAdapter create instance of adapter
func createAdapter(adapterCfg *conf.AdapterConfig, adapterCfgPath string, log *zap.SugaredLogger) (ChainAdapter, error) {
	adapterProvider := Provider(adapterCfg.Provider)
	if adapterProvider == ChainMakerProvider {
		log.Infof("create chainmaker adapter, chain id: [%s]", adapterCfg.ChainID)
		return chainmaker.NewChainMakerAdapter(adapterCfg.ChainID, adapterCfgPath, adapterCfg.ProofContract, log)
	}
	if adapterProvider == FabricProvider {
		log.Infof("create fabric adapter, chain id: [%s]", adapterCfg.ChainID)
		return fabric.NewFabricAdapter(adapterCfg.ChainID, adapterCfgPath, adapterCfg.ProofContract, log)
	}
	panic(fmt.Sprintf("can not find adapters for %v", adapterProvider))
}

// finalAdapterCfgPath return final config path of adapter
func finalAdapterCfgPath(adapterCfg *conf.AdapterConfig) string {
	return conf.FinalCfgPath(adapterCfg.ConfigPath)
}
