/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract

import (
	common "management_backend/src/db/common"
	"management_backend/src/db/connection"
	loggers "management_backend/src/logger"
)

var (
	log = loggers.GetLogger(loggers.ModuleDb)
)

type ContractStatistics struct {
	Id               int64
	ContractName     string
	ContractVersion  string
	ContractOperator string
	TxNum            int
	Timestamp        int64
}

func CreateContract(contract *common.Contract) error {
	// 此处是创建，而非更新
	if err := connection.DB.Create(&contract).Error; err != nil {
		log.Error("Save contract Failed: " + err.Error())
		return err
	}
	return nil
}

func GetContractById(chainId string, id uint64) (*common.Contract, error) {
	var contract common.Contract
	if err := connection.DB.Model(contract).Where("chain_id = ?", chainId).Where("id = ?", id).
		Find(&contract).Error; err != nil {
		log.Error("GetContractById Failed: " + err.Error())
		return nil, err
	}
	return &contract, nil
}

func GetContract(id uint64) (*common.Contract, error) {
	var contract common.Contract
	if err := connection.DB.Model(contract).Where("id = ?", id).
		Find(&contract).Error; err != nil {
		log.Error("GetContractBy Failed: " + err.Error())
		return nil, err
	}
	return &contract, nil
}

func GetContractByChainId(pageNum int64, pageSize int, chainId, contractName string) (
	[]*common.Contract, int64, error) {
	var contracts []*common.Contract

	db := connection.DB
	if contractName != "" {
		db = db.Where("name = ?", contractName)
	}

	db = db.Where("chain_id = ?", chainId)

	offset := pageNum * int64(pageSize)
	if err := db.Order("id DESC").Offset(offset).Limit(pageSize).Find(&contracts).Error; err != nil {
		log.Error("GetContractByChainId Failed: " + err.Error())
		return nil, 0, err
	}
	var count int64
	if err := db.Model(&contracts).Count(&count).Error; err != nil {
		log.Error("GetContractByChainIdCount Failed: " + err.Error())
		return nil, 0, err
	}
	return contracts, count, nil
}

func GetContractList(chainId string) ([]*common.Contract, error) {
	var contracts []*common.Contract
	if err := connection.DB.Where("chain_id = ? AND multi_sign_status = ? "+
		"AND (contract_status = ? OR contract_status = ? OR contract_status = ?)",
		chainId, common.NO_VOTING, common.ContractInitOK, common.ContractUnfreezeOK, common.ContractUpgradeOK).
		Find(&contracts).Error; err != nil {
		log.Error("GetContractList Failed: " + err.Error())
		return nil, err
	}
	return contracts, nil
}

func GetContractByName(chainId string, name string) (*common.Contract, error) {
	var contract common.Contract
	if err := connection.DB.Model(contract).Where("chain_id = ?", chainId).Where("name = ?", name).
		Find(&contract).Error; err != nil {
		log.Error("GetContractByName Failed: " + err.Error())
		return nil, err
	}
	return &contract, nil
}

func GetContractStatisticsList(chainId string, contractName string, offset int, limit int) (
	int64, []*ContractStatistics, error) {
	var (
		count        int64
		contractList []*ContractStatistics
		err          error
	)

	if err = connection.DB.Model(&common.Contract{}).Where("chain_id = ?", chainId).
		Count(&count).Error; err != nil {
		log.Error("GetContractList Failed: " + err.Error())
		return count, contractList, err
	}

	contractSelector := connection.DB.Table(common.TableContract+" contract").Order("id").
		Select("contract.id as id, "+
			"contract.name as contract_name, "+
			"contract.version as contract_version, "+
			"contract.contract_operator, "+
			"contract.timestamp, "+
			"count(tx.id) as tx_num").
		Joins("LEFT JOIN "+common.TableTransaction+" tx on contract.name = tx.contract_name or contract.evm_address = tx.contract_name ").
		Where("contract.chain_id = ?", chainId).
		Group("contract.id")

	if contractName != "" {
		count = 1
		contractSelector = contractSelector.Where("contract.name = ? or contract.evm_address = ?", contractName, contractName)
	}

	if err = contractSelector.Order("contract.create_at desc").Offset(offset).Limit(limit).
		Scan(&contractList).Error; err != nil {
		log.Error("GetContractStatisticsList Failed: " + err.Error())
		return count, contractList, err
	}
	return count, contractList, err
}

func UpdateContractMultiSignStatus(contract *common.Contract) error {
	if err := connection.DB.Debug().Model(contract).Where("name = ?", contract.Name).
		UpdateColumn("multi_sign_status", contract.MultiSignStatus).
		Error; err != nil {
		log.Error("UpdateContractColumns multi_sign_status failed: " + err.Error())
		return err
	}
	return nil
}

func UpdateContractMethod(contract *common.Contract) error {
	if err := connection.DB.Debug().Model(contract).Where("id = ?", contract.Id).
		UpdateColumn("methods", contract.Methods).
		UpdateColumn("evm_abi_save_key", contract.EvmAbiSaveKey).
		UpdateColumn("evm_function_type", contract.EvmFunctionType).
		Error; err != nil {
		log.Error("UpdateContractColumns methods failed: " + err.Error())
		return err
	}
	return nil
}

func UpdateContractMethodByName(contract *common.Contract) error {
	if err := connection.DB.Debug().Model(contract).Where("name = ?", contract.Name).
		UpdateColumn("methods", contract.Methods).
		UpdateColumn("source_save_key", contract.SourceSaveKey).
		UpdateColumn("evm_abi_save_key", contract.EvmAbiSaveKey).
		UpdateColumn("evm_function_type", contract.EvmFunctionType).
		Error; err != nil {
		log.Error("UpdateContractColumns methods failed: " + err.Error())
		return err
	}
	return nil
}

func GetContractCountByChainId(chainId string) (int64, error) {
	var count int64
	if err := connection.DB.Model(&common.Contract{}).Where("chain_id = ?", chainId).
		Count(&count).Error; err != nil {
		log.Error("GetContractCountByChainId Failed: " + err.Error())
		return 0, err
	}
	return count, nil
}

func UpdateContractStatus(id int64, status int, voteStatus int) error {
	var contract = &common.Contract{}
	if err := connection.DB.Debug().Model(contract).Where("id = ?", id).
		UpdateColumn("contract_status", status).
		UpdateColumn("multi_sign_status", voteStatus).Error; err != nil {
		log.Error("UpdateContractColumns failed: " + err.Error())
		return err
	}
	return nil
}

//func UpdateInstallContractStatus(id int64, status int, voteStatus int, txId string) error {
//	var contract = &common.Contract{}
//	if err := connection.DB.Debug().Model(contract).Where("id = ?", id).
//		UpdateColumn("contract_status", status).
//		UpdateColumn("tx_id", txId).
//		UpdateColumn("multi_sign_status", voteStatus).Error; err != nil {
//		log.Error("UpdateContractColumns failed: " + err.Error())
//		return err
//	}
//	return nil
//}
