/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"context"
	"management_backend/src/db/connection"
	"strconv"
	"sync"
	"time"

	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	sdkconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"

	"management_backend/src/config"
	"management_backend/src/db/chain"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/policy"
	"management_backend/src/db/relation"
	"management_backend/src/entity"
	"management_backend/src/utils"
)

const (
	ADMIN = iota
	CLIENT
	ALL
)

const (
	NO_SELECTED = iota
	SELECTED
)

type SdkClientPool struct {
	SdkClients map[string]*SdkClient
}

type SdkClient struct {
	lock        sync.Mutex
	ChainId     string
	SdkConfig   *entity.SdkConfig
	ChainClient *sdk.ChainClient
}

const (
	BlockUpdate      = "CHAIN_CONFIG-BLOCK_UPDATE"
	PermissionUpdate = "CHAIN_CONFIG-PERMISSION_UPDATE"
	InitContract     = "CONTRACT_MANAGE-INIT_CONTRACT"
	UpgradeContract  = "CONTRACT_MANAGE-UPGRADE_CONTRACT"
	FreezeContract   = "CONTRACT_MANAGE-FREEZE_CONTRACT"
	UnfreezeContract = "CONTRACT_MANAGE-UNFREEZE_CONTRACT"
	RevokeContract   = "CONTRACT_MANAGE-REVOKE_CONTRACT"
)

var (
	ResourceNameMap = map[string]int{
		BlockUpdate:      3,
		InitContract:     4,
		UpgradeContract:  5,
		FreezeContract:   6,
		UnfreezeContract: 7,
		RevokeContract:   8,
		PermissionUpdate: 9}

	ResourceNameValueMap = map[int]string{
		0: syscontract.ChainConfigFunction_NODE_ID_UPDATE.String(),
		1: syscontract.ChainConfigFunction_TRUST_ROOT_UPDATE.String(),
		2: syscontract.ChainConfigFunction_CONSENSUS_EXT_UPDATE.String(),
		3: BlockUpdate,
		4: InitContract,
		5: UpgradeContract,
		6: FreezeContract,
		7: UnfreezeContract,
		8: RevokeContract,
		9: PermissionUpdate}

	RuleMap = map[string]int{"MAJORITY": 0, "ANY": 1, "SELF": 2, "ALL": 3, "FORBIDDEN": 4, "PERCENTAGE": 4}
	RoleMap = map[string]int{"admin": 0, "client": 1}

	RuleValueMap = map[int]string{0: "MAJORITY", 1: "ANY", 2: "SELF", 3: "ALL", 4: "FORBIDDEN", 5: "PERCENTAGE"}
	RoleValueMap = map[int]string{0: "admin", 1: "client"}
)

func NewSdkClient(sdkConfig *entity.SdkConfig) (*SdkClient, error) {

	client, err := CreateSdkClientWithChainId(sdkConfig)
	if err != nil {
		return nil, err
	}
	return &SdkClient{ChainId: sdkConfig.ChainId, ChainClient: client, SdkConfig: sdkConfig}, nil
}

func NewSdkClientPool(sdkClient *SdkClient) *SdkClientPool {
	sdkClients := make(map[string]*SdkClient)
	sdkClients[sdkClient.ChainId] = sdkClient
	return &SdkClientPool{
		SdkClients: sdkClients,
	}
}

func (sdkClient *SdkClient) Load() {
	chainId := sdkClient.ChainId
	log.Debugf("[WEB] begin to load chain's information, [chain:%s] ", chainId)
	// update chain info at times
	go sdkClient.loadChainAtFixedTime()
	// get max block height for this chain
	maxBlockHeight := chain.GetMaxBlockHeight(chainId)
	// start listener
	go blockListenStart(sdkClient, maxBlockHeight)
}

func (sdkClient *SdkClient) loadChainAtFixedTime() {
	// update first, which update at times after
	ticker := time.NewTicker(time.Second * time.Duration(config.GlobalConfig.WebConf.LoadPeriodSeconds))
	for {
		<-ticker.C
		LoadChainRefInfos(sdkClient)
	}
}

func LoadChainRefInfos(sdkClient *SdkClient) {
	loadOrgInfo(sdkClient)
	loadNodeInfo(sdkClient)
	loadChainInfo(sdkClient)
	loadChainErrorLog(sdkClient)
}

func loadChainInfo(sdkClient *SdkClient) *dbcommon.Chain {
	sdkClient.lock.Lock()
	defer sdkClient.lock.Unlock()

	chainClient := sdkClient.ChainClient
	chainConfig, err := chainClient.GetChainConfig()
	if err != nil {
		log.Error("[SDK] Get Chain Config Failed : " + err.Error())
		return nil
	}

	var chainInfo dbcommon.Chain
	chainInfo.ChainId = chainConfig.ChainId
	chainInfo.BlockInterval = chainConfig.Block.BlockInterval
	chainInfo.BlockTxCapacity = chainConfig.Block.BlockTxCapacity
	chainInfo.TxTimeout = chainConfig.Block.TxTimeout
	chainInfo.Consensus = chainConfig.Consensus.Type.String()
	chainInfo.Version = chainConfig.Version
	chainInfo.Sequence = strconv.FormatUint(chainConfig.Sequence, 10)
	err = chain.UpdateChainInfo(&chainInfo)
	if err != nil {
		log.Error("[SDK] Update Chain Config Failed : " + err.Error())
		return nil
	}
	var roleType int

	chainOrgList, err := relation.GetChainOrgList(chainConfig.ChainId)
	if err != nil {
		log.Error("GetChainOrgList: " + err.Error())
	}

	resourcePolicyList := chainConfig.ResourcePolicies
	resourcePolicyList = addConfigPolicy(resourcePolicyList)

	for _, resourcePolicy := range resourcePolicyList {
		resourceName := resourcePolicy.ResourceName

		resourceType, ok := ResourceNameMap[resourceName]
		if !ok {
			continue
		}

		if len(resourcePolicy.Policy.RoleList) == 1 {
			roleType = RoleMap[resourcePolicy.Policy.RoleList[0]]
		} else {
			roleType = ALL
		}
		chainPolicy := &dbcommon.ChainPolicy{
			ChainId:    chainConfig.ChainId,
			Type:       resourceType,
			PolicyType: RuleMap[resourcePolicy.Policy.Rule],
			RoleType:   roleType,
			PercentNum: resourcePolicy.Policy.Rule,
		}
		var chainPolicyOrgList []*dbcommon.ChainPolicyOrg
		if len(resourcePolicy.Policy.OrgList) == 0 {
			for _, chainOrg := range chainOrgList {
				chainPolicyOrg := &dbcommon.ChainPolicyOrg{
					OrgName: chainOrg.OrgName,
					OrgId:   chainOrg.OrgId,
					Status:  SELECTED,
				}
				chainPolicyOrgList = append(chainPolicyOrgList, chainPolicyOrg)
			}
		} else {
			orgIdMap := make(map[string]string)
			for _, orgId := range resourcePolicy.Policy.OrgList {
				orgIdMap[orgId] = orgId
			}
			for _, chainOrg := range chainOrgList {
				chainPolicyOrg := &dbcommon.ChainPolicyOrg{
					OrgName: chainOrg.OrgName,
					OrgId:   chainOrg.OrgId,
					Status:  NO_SELECTED,
				}
				if orgIdMap[chainOrg.OrgId] != "" {
					chainPolicyOrg.Status = SELECTED
				}
				chainPolicyOrgList = append(chainPolicyOrgList, chainPolicyOrg)
			}

		}
		err := policy.CreateChainPolicy(chainPolicy, chainPolicyOrgList)
		if err != nil {
			log.Error("[SDK] Save chainPolicy Failed : " + err.Error())
		}
	}

	return &chainInfo
}

func addConfigPolicy(resourcePolicyList []*sdkconfig.ResourcePolicy) []*sdkconfig.ResourcePolicy {
	for resourceName, _ := range ResourceNameMap {
		add := true
		for _, resourcePolicy := range resourcePolicyList {
			if resourceName == resourcePolicy.ResourceName {
				add = false
				break
			}
		}
		if add {
			policyInfo := &accesscontrol.Policy{
				Rule:     "MAJORITY",
				OrgList:  nil,
				RoleList: []string{"admin"},
			}
			resourcePolicy := &sdkconfig.ResourcePolicy{
				ResourceName: resourceName,
				Policy:       policyInfo,
			}
			resourcePolicyList = append(resourcePolicyList, resourcePolicy)
		}
	}

	// permission update 的操作不太适合普通用户去修改，因此在此写死，在此默认是all，admin；
	permissionPolicy := &accesscontrol.Policy{
		Rule:     "ALL",
		OrgList:  nil,
		RoleList: []string{"admin"},
	}
	permissionResource := &sdkconfig.ResourcePolicy{
		ResourceName: PermissionUpdate,
		Policy:       permissionPolicy,
	}
	resourcePolicyList = append(resourcePolicyList, permissionResource)

	return resourcePolicyList
}

func loadNodeInfo(sdkClient *SdkClient) {
	sdkClient.lock.Lock()
	defer sdkClient.lock.Unlock()

	chainClient := sdkClient.ChainClient
	chainInfo, err := chainClient.GetChainInfo()
	if err != nil {
		log.Error("[SDK] Get Chain Info Failed : " + err.Error())
		if err.Error() != "connections are busy" {
			chainAddNode(sdkClient.SdkConfig.ChainId)
		}
		return
	}
	nodeList := chainInfo.NodeList
	for _, node := range nodeList {
		var nodeName string
		dbNode, err := chain_participant.GetNodeByNodeId(node.NodeId)
		if err != nil {
			log.Error("[SDK] Get Node Info Failed : " + err.Error())
			nodeName = node.NodeId
		}
		if dbNode != nil {
			nodeName = dbNode.NodeName
		}
		orgNodeList, err := relation.GetOrgNodeByNodeId(node.NodeId)
		if err != nil {
			log.Error("[SDK] Get Org Node Info Failed : " + err.Error())
		}
		chainOrgList, err := relation.GetChainOrgList(sdkClient.SdkConfig.ChainId)
		if err != nil {
			log.Error("[SDK] Get Chain Org Info Failed : " + err.Error())
			break
		}
		var orgId string
		var orgName string
		if len(orgNodeList) > 0 {
			for _, orgNode := range orgNodeList {
				for _, chainOrg := range chainOrgList {
					if orgNode.OrgId == chainOrg.OrgId {
						orgId = chainOrg.OrgId
						orgName = chainOrg.OrgName
					}
				}
			}
		}
		chainOrgNode := &dbcommon.ChainOrgNode{
			ChainId:  sdkClient.SdkConfig.ChainId,
			OrgId:    orgId,
			OrgName:  orgName,
			NodeId:   node.NodeId,
			NodeName: nodeName,
		}
		err = relation.CreateChainOrgNode(chainOrgNode)
		if err != nil {
			log.Error("CreateChainOrgNode err : " + err.Error())
		}
	}
}

func chainAddNode(chainId string) {
	chainOrgList, err := relation.GetChainOrgList(chainId)
	if err != nil {
		log.Error("[SDK] Get Chain Org Info Failed : " + err.Error())
		return
	}
	for _, chainOrg := range chainOrgList {
		orgNodeList, err := relation.GetOrgNode(chainOrg.OrgId)
		if err != nil {
			log.Error("[SDK] Get Org Node Info Failed : " + err.Error())
			return
		}
		for _, orgNode := range orgNodeList {
			chainOrgNode := &dbcommon.ChainOrgNode{
				ChainId:  chainId,
				OrgId:    orgNode.OrgId,
				OrgName:  orgNode.OrgName,
				NodeId:   orgNode.NodeId,
				NodeName: orgNode.NodeName,
			}
			err = relation.CreateChainOrgNode(chainOrgNode)
			if err != nil {
				log.Error("CreateChainOrgNode err : " + err.Error())
			}
		}
	}
}

func loadOrgInfo(sdkClient *SdkClient) {
	sdkClient.lock.Lock()
	defer sdkClient.lock.Unlock()

	chainClient := sdkClient.ChainClient
	chainConfig, err := chainClient.GetChainConfig()
	if err != nil {
		log.Error("[SDK] Get Chain config Failed : " + err.Error())
		return
	}
	trustRoots := chainConfig.TrustRoots

	orgIdMap := make(map[string]string)
	for _, trustRoot := range trustRoots {
		var orgName string
		orgName, err = chain_participant.GetOrgNameByOrgId(trustRoot.OrgId)
		if err != nil {
			orgName = trustRoot.OrgId
			log.Error("[SDK] Get Org Name Failed : " + err.Error())
		}
		orgIdMap[trustRoot.OrgId] = trustRoot.OrgId
		chainOrg := &dbcommon.ChainOrg{
			ChainId: chainConfig.ChainId,
			OrgId:   trustRoot.OrgId,
			OrgName: orgName,
		}
		err = relation.CreateChainOrg(chainOrg)
		if err != nil {
			log.Error("CreateChainOrg err : " + err.Error())
		}
	}
}

func loadChainErrorLog(sdkClient *SdkClient) {
	chainId := sdkClient.ChainId
	chainInfo, err := chain.GetChainByChainId(chainId)
	if err != nil {
		log.Error("failed get chain info")
		return
	}
	// 不允许监控 则返回
	if chainInfo.Monitor == 0 {
		return
	}

	host := utils.GetHostFromAddress(sdkClient.SdkConfig.Remote)
	err = PullChainErrorLog(host)
	if err != nil {
		log.Error("failed fetch chain error log")
	}
	return
}

func blockListenStart(sdkClient *SdkClient, maxBlockHeight int64) {
	sdkClient.lock.Lock()
	chainClient := sdkClient.ChainClient
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var startBlock int64
	if maxBlockHeight > 0 {
		startBlock = maxBlockHeight - 1
	} else {
		startBlock = 0
	}
	c, err := chainClient.SubscribeBlock(ctx, startBlock, -1, true, false)
	if err != nil {
		log.Error("[Sync Block] Get Block By SDK failed: " + err.Error())
	}
	sdkClient.lock.Unlock()
	pool := NewPool(1)
	go pool.Run()
	for {
		select {
		case block, ok := <-c:
			if !ok {
				log.Error("Chan Is Closed")
				updateChainNoWork(sdkClient.ChainId)
				return
			}

			blockInfo, ok := block.(*common.BlockInfo)
			if !ok {
				log.Error("The Data Type Error")
				updateChainNoWork(sdkClient.ChainId)
				return
			}
			pool.EntryChan <- NewTask(storageBlock(blockInfo))
		case <-ctx.Done():
			return
		}
	}
}

func updateChainNoWork(chainId string) {
	var chainInfo dbcommon.Chain
	chainInfo.Status = connection.NO_WORK
	chainInfo.ChainId = chainId
	err := chain.UpdateChainStatus(&chainInfo)
	if err != nil {
		log.Error("[SDK] Update Chain Config Failed : " + err.Error())
	}
}

func storageBlock(blockInfo *common.BlockInfo) func() error {
	return func() error {
		err := ParseBlockToDB(blockInfo)
		if err != nil {
			log.Error("Storage Block Failed: " + err.Error())
			return err
		}
		return nil
	}
}

// addSdkClient add SDKClient
func (pool *SdkClientPool) AddSdkClient(sdkClient *SdkClient) error {

	sdkClients := pool.SdkClients
	sdkClients[sdkClient.ChainId] = sdkClient
	pool.SdkClients = sdkClients
	return nil
}

func (pool *SdkClientPool) LoadChains() {
	// one chain one goroutine
	sdkClients := pool.SdkClients
	for _, sdkClient := range sdkClients {
		go sdkClient.Load()
	}
}
