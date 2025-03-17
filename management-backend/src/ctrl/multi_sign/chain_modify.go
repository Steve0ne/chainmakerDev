/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package multi_sign

import (
	"encoding/json"

	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	sdkcommon "chainmaker.org/chainmaker/pb-go/v2/common"

	"management_backend/src/ctrl/common"
	"management_backend/src/ctrl/overview"
	"management_backend/src/db/chain"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/sync"
)

func ChainAuthModify(parameters string, orgList []string, roleType int) error {
	var chainAuthBody overview.ModifyChainAuthParams
	err := json.Unmarshal([]byte(parameters), &chainAuthBody)
	if err != nil {
		log.Errorf("Unmarshal parameters to chainAuthBody err:, %s", err)
		return err
	}

	chainId := chainAuthBody.ChainId
	sdkClientPool := sync.GetSdkClientPool()
	if sdkClientPool == nil {
		newError := common.CreateError(common.ErrorChainNotSub)
		return newError
	}
	var roleList []string
	var policyOrgList []string
	for _, role := range chainAuthBody.RoleList {
		roleList = append(roleList, sync.RoleValueMap[role.Role])
	}
	for _, org := range chainAuthBody.OrgList {
		policyOrgList = append(policyOrgList, org.OrgId)
	}

	var rule string
	if chainAuthBody.Rule == 5 {
		rule = chainAuthBody.PercentNum
	} else {
		rule = sync.RuleValueMap[chainAuthBody.Rule]
	}

	policy := &accesscontrol.Policy{
		Rule:     rule,
		OrgList:  policyOrgList,
		RoleList: roleList,
	}

	sdkClient := sdkClientPool.SdkClients[chainId]
	chainClient := sdkClient.ChainClient
	payload, err := chainClient.CreateChainConfigPermissionUpdatePayload(
		sync.ResourceNameValueMap[chainAuthBody.Type], policy)
	if err != nil {
		newError := common.CreateError(common.ErrorCreateChainConfigPermissionUpdatePayload)
		return newError
	}

	endorsements, err := GetEndorsements(payload, orgList, roleType)
	if err != nil {
		newError := common.CreateError(common.ErrorMergeSign)
		return newError
	}
	resp, err := chainClient.SendContractManageRequest(payload, endorsements, TxHandleTimeout, true)
	if err != nil {
		log.Error("invoke contract err : %s", err.Error())
		return err
	}
	if resp.Code != sdkcommon.TxStatusCode_SUCCESS {
		log.Errorf("Send ChainConfigUpdate failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
		newError := common.CreateError(common.ErrorUpdateChainConfig)
		return newError
	}

	return nil
}

func ChainConfigModify(parameters string, orgList []string, roleType int) error {
	var chainConfigBody overview.ModifyChainConfigParams
	err := json.Unmarshal([]byte(parameters), &chainConfigBody)
	if err != nil {
		log.Errorf("Unmarshal parameters to chainConfigBody err:, %s", err)
		return err
	}
	chainId := chainConfigBody.ChainId
	// 检查与之前配置是否一致
	chainInfo, err := chain.GetChainByChainId(chainId)
	if err != nil {
		log.Error("GetChainByChainId err : " + err.Error())
		return err
	}
	if equalChainConfigs(&chainConfigBody, chainInfo) {
		// 相等表示不需要更新，直接返回正常即可
		return nil
	}

	sdkClientPool := sync.GetSdkClientPool()
	if sdkClientPool == nil {
		newError := common.CreateError(common.ErrorChainNotSub)
		return newError
	}
	sdkClient := sdkClientPool.SdkClients[chainId]
	chainClient := sdkClient.ChainClient

	payload, err := chainClient.CreateChainConfigBlockUpdatePayload(
		true, chainConfigBody.TxTimeout, chainConfigBody.BlockTxCapacity,
		10, chainConfigBody.BlockInterval)
	if err != nil {
		newErr := common.CreateError(common.ErrorCreateChainConfigBlockUpdatePayload)
		return newErr
	}

	endorsements, err := GetEndorsements(payload, orgList, roleType)
	if err != nil {
		newErr := common.CreateError(common.ErrorMergeSign)
		return newErr
	}

	resp, err := chainClient.SendChainConfigUpdateRequest(payload, endorsements, -1, true)
	if err != nil {
		newErr := common.CreateError(common.ErrorUpdateChainConfig)
		return newErr
	}
	// 判断结果
	if resp.Code != sdkcommon.TxStatusCode_SUCCESS {
		newError := common.CreateError(common.ErrorUpdateChainConfig)
		return newError
	}
	return nil
}

func equalChainConfigs(chainBody *overview.ModifyChainConfigParams, dbChain *dbcommon.Chain) bool {
	return chainBody.TxTimeout == dbChain.TxTimeout &&
		chainBody.BlockTxCapacity == dbChain.BlockTxCapacity &&
		chainBody.BlockInterval == dbChain.BlockInterval
}
