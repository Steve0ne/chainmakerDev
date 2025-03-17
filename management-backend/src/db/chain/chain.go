/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain

import (
	"management_backend/src/db/common"
	"management_backend/src/db/connection"
	loggers "management_backend/src/logger"
)

var (
	log = loggers.GetLogger(loggers.ModuleDb)
)

func CreateChain(chain *common.Chain) error {
	if err := connection.DB.Create(&chain).Error; err != nil {
		log.Error("[DB] Save chain Failed: " + err.Error())
		return err
	}
	return nil
}

func GetChainByChainId(chainId string) (*common.Chain, error) {
	var chain common.Chain
	if err := connection.DB.Where("chain_id = ?", chainId).Find(&chain).Error; err != nil {
		log.Error("GetChainByChainId Failed: " + err.Error())
		return nil, err
	}
	return &chain, nil
}

func GetChainById(id int64) (*common.Chain, error) {
	var chain common.Chain
	if err := connection.DB.Where("id = ?", id).Find(&chain).Error; err != nil {
		log.Error("GetChainByChainId Failed: " + err.Error())
		return nil, err
	}
	return &chain, nil
}

func GetChainByChainIdOrName(chainId, chainName string) (*common.Chain, error) {
	var chain common.Chain
	if err := connection.DB.Where("chain_id = ? || chain_name = ?", chainId, chainName).Find(&chain).Error; err != nil {
		log.Error("GetChainByChainIdOrName Failed: " + err.Error())
		return nil, err
	}
	return &chain, nil
}

func GetChainList() ([]*common.Chain, error) {
	var chains []*common.Chain
	if err := connection.DB.Order("id DESC").Find(&chains).Error; err != nil {
		log.Error("GetChainList Failed: " + err.Error())
		return nil, err
	}
	return chains, nil
}

func UpdateChainInfo(chain *common.Chain) error {
	chainId := chain.ChainId
	_, err := GetChainByChainId(chainId)
	if err != nil {
		// 插入即可
		chain.Status = connection.START
		chain.ChainName = chainId
		return CreateChain(chain)
	}
	// 修改配置，包括
	if err := connection.DB.Debug().Model(chain).Where("chain_id = ?", chain.ChainId).
		UpdateColumns(getChainUpdateColumns(chain)).Error; err != nil {
		log.Error("UpdateChainInfo failed: " + err.Error())
		return err
	}
	return nil
}

func UpdateChainStatus(chain *common.Chain) error {
	columns := make(map[string]interface{})
	columns["status"] = chain.Status
	if err := connection.DB.Debug().Model(chain).Where("chain_id = ?", chain.ChainId).
		UpdateColumns(columns).Error; err != nil {
		log.Error("UpdateChainStatus failed: " + err.Error())
		return err
	}
	return nil
}

func DeleteChain(chainId string) error {
	tx := connection.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}
	// handle chain

	if err := tx.Debug().Where("chain_id = ?", chainId).Delete(&common.Chain{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// handle chainOrg
	if err := tx.Debug().Where("chain_id = ?", chainId).Delete(&common.ChainOrg{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// handle chainOrgNode
	if err := tx.Debug().Where("chain_id = ?", chainId).Delete(&common.ChainOrgNode{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func getChainUpdateColumns(chain *common.Chain) map[string]interface{} {
	columns := make(map[string]interface{})
	columns["tx_timeout"] = chain.TxTimeout
	columns["block_tx_capacity"] = chain.BlockTxCapacity
	columns["block_interval"] = chain.BlockInterval
	columns["status"] = connection.START
	columns["version"] = chain.Version
	columns["sequence"] = chain.Sequence
	columns["consensus"] = chain.Consensus
	return columns
}
