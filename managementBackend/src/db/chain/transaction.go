/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain

import (
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"

	"management_backend/src/db/common"
	"management_backend/src/db/connection"
)

func GetTxById(id uint64) (*common.Transaction, error) {
	var tx common.Transaction
	if err := connection.DB.Model(tx).Where("id = ?", id).Find(&tx).Error; err != nil {
		log.Error("[DB] QueryTxById Failed: " + err.Error())
		return nil, err
	}
	return &tx, nil
}

func GetTxByTxId(chainId string, txId string) (*common.Transaction, error) {
	var tx common.Transaction
	if err := connection.DB.Model(tx).Where("chain_id = ?", chainId).
		Where("binary tx_id = ?", txId).Find(&tx).Error; err != nil {
		log.Error("[DB] QueryTxById Failed: " + err.Error())
		return nil, err
	}
	return &tx, nil
}

func GetTxByContractName(chainId string, contractName string, evmAddress string) (*common.Transaction, error) {
	var tx common.Transaction

	if err := connection.DB.Model(tx).Where("chain_id = ?", chainId).
		Where("contract_name = ?", contractName).
		Or("contract_name = ?", evmAddress).
		Where("contract_method = ?", syscontract.ContractManageFunction_INIT_CONTRACT.String()).
		Order("id DESC").Limit(1).
		Find(&tx).Error; err != nil {
		log.Error("[DB] GetTxByContractName Failed: " + err.Error())
		return nil, err
	}
	return &tx, nil
}

func GetTxList(chainId string, offset int, limit int, blockHeight *int64, contractName string) (
	int64, []*common.Transaction, error) {
	var (
		count  int64
		txList []*common.Transaction
		err    error
	)

	txSelector := connection.DB.Model(&common.Transaction{}).Where("chain_id = ?", chainId)

	if blockHeight != nil {
		txSelector = txSelector.Where("block_height = ?", *blockHeight)
	}

	if contractName != "" {
		txSelector = txSelector.Where("contract_name = ?", contractName)
	}

	if err = txSelector.Count(&count).Error; err != nil {
		log.Error("GetTxList Failed: " + err.Error())
		return count, txList, err
	}

	// 交易排序的逻辑 按时间戳倒序，应该按序号倒序
	if err = txSelector.Order("timestamp desc").Offset(offset).Limit(limit).Find(&txList).Error; err != nil {
		log.Error("GetTxList Failed: " + err.Error())
		return count, txList, err
	}
	return count, txList, err
}

func GetTxNumByChainId(chainId string) (int64, error) {
	var txNum int64
	if err := connection.DB.Model(&common.Transaction{}).Where("chain_id = ?", chainId).Count(&txNum).Error; err != nil {
		log.Error("GetTxNumByChainId Failed: " + err.Error())
		return 0, err
	}
	return txNum, nil
}
