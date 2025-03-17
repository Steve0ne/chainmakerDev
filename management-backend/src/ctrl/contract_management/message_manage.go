/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_management

import (
	"encoding/json"
	"management_backend/src/db/chain"
	dbcontract "management_backend/src/db/contract"
	"strconv"
	"strings"
)

var (
	MessageTplMap = map[int]string{
		BLOCK_UPDATE: "StartOrgId申请修改链配置，其中，区块最大容量由OldBlockTxCapacity笔修改为BlockTxCapacity笔，出块间隔由OldBlockIntervalms修改为BlockIntervalms," +
			"交易过期时长由OldTxTimeouts修改为TxTimeouts",
		INIT_CONTRACT:     "StartOrgId申请部署ContractName合约",
		UPGRADE_CONTRACT:  "StartOrgId申请将ContractName合约由OldContractVersion版本升级到ContractVersion版本",
		FREEZE_CONTRACT:   "StartOrgId申请冻结ContractName合约",
		UNFREEZE_CONTRACT: "StartOrgId申请解冻ContractName合约",
		REVOKE_CONTRACT:   "StartOrgId申请废止ContractName合约",
		PERMISSION_UPDATE: "StartOrgId申请修改链权限，在ChainId链上修改链权限Type规则为Rule,更新组织列表为OrgList,更新角色列表为RoleList",
	}
)

const (
	START_ORG_ID          = "StartOrgId"
	CHAIN_ID              = "ChainId"
	BLOCK_TX_CAPACITY     = "BlockTxCapacity"
	BLOCK_INTERVAL        = "BlockInterval"
	TX_TIMEOUT            = "TxTimeout"
	OLD_BLOCK_TX_CAPACITY = "OldBlockTxCapacity"
	OLD_BLOCK_INTERVAL    = "OldBlockInterval"
	OLD_TX_TIMEOUT        = "OldTxTimeout"
	CONTRACT_NAME         = "ContractName"
	OLD_CONTRACT_VERSION  = "OldContractVersion"
	CONTRACT_VERSION      = "ContractVersion"
	TYPE                  = "Type"
	RULE                  = "Rule"
	ORG_LIST              = "OrgList"
	ROLE_LIST             = "RoleList"
)

func getMessage(chainId, startOrgId, paramJson string, voteType int) (string, error) {
	var message string
	var paramMap map[string]interface{}
	err := json.Unmarshal([]byte(paramJson), &paramMap)
	if err != nil {
		return "", err
	}
	if voteType == BLOCK_UPDATE {
		message = MessageTplMap[BLOCK_UPDATE]
		message = replaceConmmon(message, chainId, startOrgId)

		chainInfo, err := chain.GetChainByChainId(chainId)
		if err != nil {
			return "", err
		}

		message = strings.Replace(message, OLD_BLOCK_TX_CAPACITY, strconv.Itoa(int(chainInfo.BlockTxCapacity)), -1)
		message = strings.Replace(message, OLD_BLOCK_INTERVAL, strconv.Itoa(int(chainInfo.BlockInterval)), -1)
		message = strings.Replace(message, OLD_TX_TIMEOUT, strconv.Itoa(int(chainInfo.TxTimeout)), -1)

		message = strings.Replace(message, BLOCK_TX_CAPACITY, strconv.FormatFloat(paramMap[BLOCK_TX_CAPACITY].(float64), 'f', -1, 64), -1)
		message = strings.Replace(message, BLOCK_INTERVAL, strconv.FormatFloat(paramMap[BLOCK_INTERVAL].(float64), 'f', -1, 64), -1)
		message = strings.Replace(message, TX_TIMEOUT, strconv.FormatFloat(paramMap[TX_TIMEOUT].(float64), 'f', -1, 64), -1)
	}
	if voteType == INIT_CONTRACT {
		message = MessageTplMap[INIT_CONTRACT]
		message = replaceConmmon(message, chainId, startOrgId)
		message = strings.Replace(message, CONTRACT_NAME, paramMap[CONTRACT_NAME].(string), -1)
	}
	if voteType == UPGRADE_CONTRACT {
		message = MessageTplMap[UPGRADE_CONTRACT]
		message = replaceConmmon(message, chainId, startOrgId)
		message = strings.Replace(message, CONTRACT_NAME, paramMap[CONTRACT_NAME].(string), -1)

		contractInfo, err := dbcontract.GetContractByName(paramMap[CHAIN_ID].(string), paramMap[CONTRACT_NAME].(string))
		if err != nil {
			return "", err
		}
		message = strings.Replace(message, OLD_CONTRACT_VERSION, contractInfo.Version, -1)
		message = strings.Replace(message, CONTRACT_VERSION, paramMap[CONTRACT_VERSION].(string), -1)
	}
	if voteType == FREEZE_CONTRACT {
		message = MessageTplMap[FREEZE_CONTRACT]
		message = replaceConmmon(message, chainId, startOrgId)
		message = strings.Replace(message, CONTRACT_NAME, paramMap[CONTRACT_NAME].(string), -1)
	}
	if voteType == UNFREEZE_CONTRACT {
		message = MessageTplMap[UNFREEZE_CONTRACT]
		message = replaceConmmon(message, chainId, startOrgId)
		message = strings.Replace(message, CONTRACT_NAME, paramMap[CONTRACT_NAME].(string), -1)
	}
	if voteType == REVOKE_CONTRACT {
		message = MessageTplMap[REVOKE_CONTRACT]
		message = replaceConmmon(message, chainId, startOrgId)
		message = strings.Replace(message, CONTRACT_NAME, paramMap[CONTRACT_NAME].(string), -1)
	}
	return message, nil
}

func replaceConmmon(message, chainId, startOrgId string) string {
	message = strings.Replace(message, START_ORG_ID, startOrgId, -1)
	message = strings.Replace(message, CHAIN_ID, chainId, -1)
	return message
}
