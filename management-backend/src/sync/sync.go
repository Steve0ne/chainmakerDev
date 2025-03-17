/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"management_backend/src/entity"
	loggers "management_backend/src/logger"
)

var (
	log           = loggers.GetLogger(loggers.ModuleWeb)
	sdkClientPool *SdkClientPool
)

func SubscribeChain(sdkConfig *entity.SdkConfig) error {
	var err error
	sdkClient, err := NewSdkClient(sdkConfig)
	if err != nil {
		log.Error("创建sdkClient失败: ", err.Error())
		return err
	}
	_, err = sdkClient.ChainClient.GetChainConfig()
	if err != nil {
		log.Error("订阅链失败: ", err.Error())
		return err
	}
	sdkClientPool = GetSdkClientPool()
	if sdkClientPool != nil {
		err = sdkClientPool.AddSdkClient(sdkClient)
		if err != nil {
			log.Error("[WEB] AddSdkClient err : ", err.Error())
			return err
		}
	} else {
		sdkClientPool = NewSdkClientPool(sdkClient)
	}

	LoadChainRefInfos(sdkClient)
	sdkClientPool.LoadChains()
	return nil
}

func GetSdkClientPool() *SdkClientPool {
	if sdkClientPool == nil {
		sdkClients := make(map[string]*SdkClient)
		sdkClientPool = &SdkClientPool{
			SdkClients: sdkClients,
		}
	}
	return sdkClientPool
}
