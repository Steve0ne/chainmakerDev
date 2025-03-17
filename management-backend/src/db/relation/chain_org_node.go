/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package relation

import (
	"management_backend/src/db/common"
	"management_backend/src/db/connection"
)

func CreateChainOrgNode(chainOrgNode *common.ChainOrgNode) error {
	_, err := GetChainOrgByNodeIdAndChainId(chainOrgNode.NodeId, chainOrgNode.ChainId)
	if err == nil {
		if err := connection.DB.Debug().Model(chainOrgNode).Where("chain_id = ?", chainOrgNode.ChainId).
			Where("node_id = ?", chainOrgNode.NodeId).
			UpdateColumns(updateColumns(chainOrgNode)).Error; err != nil {
			log.Error("UpdateChainOrgNode failed: " + err.Error())
			return err
		}
		return nil
	}
	if err := connection.DB.Create(&chainOrgNode).Error; err != nil {
		log.Error("Save chainOrg Failed: " + err.Error())
		return err
	}
	return nil
}

func updateColumns(chainOrgNode *common.ChainOrgNode) map[string]interface{} {
	columns := make(map[string]interface{})
	columns["org_id"] = chainOrgNode.OrgId
	columns["org_name"] = chainOrgNode.OrgName
	columns["node_name"] = chainOrgNode.NodeName
	return columns
}

func GetChainOrgByChainIdList(chainId string) ([]*common.ChainOrgNode, error) {
	var chainOrgs []*common.ChainOrgNode
	if err := connection.DB.Where("chain_id = ?", chainId).Find(&chainOrgs).Error; err != nil {
		log.Error("GetChainOrgList Failed: " + err.Error())
		return nil, err
	}
	return chainOrgs, nil
}

func GetChainOrgByNodeIdAndChainId(nodeId, chainId string) (*common.ChainOrgNode, error) {
	var chainOrgNode common.ChainOrgNode
	if err := connection.DB.Where("node_id = ? And chain_id = ?", nodeId, chainId).Find(&chainOrgNode).Error; err != nil {
		log.Error("GetChainOrgByNodeIdAndChainId Failed: " + err.Error())
		return nil, err
	}
	return &chainOrgNode, nil
}

func GetChainOrg(orgId, chainId string) ([]*common.ChainOrgNode, error) {
	var chainOrgs []*common.ChainOrgNode

	db := connection.DB
	db = db.Where("org_id = ?", orgId)
	if chainId != "" {
		db = db.Where("chain_id = ?", chainId)
	}
	if err := db.Find(&chainOrgs).Error; err != nil {
		log.Error("QueryOrgCaCert Failed: " + err.Error())
		return nil, err
	}
	return chainOrgs, nil
}

func GetNodeCountByChainId(chainId string) (int, error) {
	var chainOrgNode common.ChainOrgNode
	var count int
	if err := connection.DB.Model(&chainOrgNode).Where("chain_id = ?", chainId).Count(&count).Error; err != nil {
		log.Error("GetOrgCountByChainId Failed: " + err.Error())
		return 0, err
	}
	return count, nil
}
