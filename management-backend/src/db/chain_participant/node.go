/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_participant

import (
	"github.com/jinzhu/gorm"

	"management_backend/src/db/common"
	"management_backend/src/db/connection"
)

const (
	NODE_CONSENSUS = 0
	NODE_COMMON    = 1
)

func CreateNode(node *common.Node) error {
	if err := connection.DB.Create(&node).Error; err != nil {
		log.Error("Save node Failed: " + err.Error())
		return err
	}
	return nil
}

func GetNodeByNodeName(nodeName string) (*common.Node, error) {
	var node common.Node
	if err := connection.DB.Where("node_name = ?", nodeName).Find(&node).Error; err != nil {
		log.Error("GetNodeByNodeName Failed: " + err.Error())
		return nil, err
	}
	return &node, nil
}

func GetNodeByNodeId(nodeId string) (*common.Node, error) {
	var node common.Node
	if err := connection.DB.Where("node_id = ?", nodeId).Find(&node).Error; err != nil {
		log.Error("GetNodeByNodeId Failed: " + err.Error())
		return nil, err
	}
	return &node, nil
}

func GetConsensusNodeByNodeName(nodeName string) (*common.Node, error) {
	var node common.Node
	if err := connection.DB.Where("node_name = ? AND type = ?", nodeName, NODE_CONSENSUS).Find(&node).Error; err != nil {
		log.Error("GetNodeByNodeName Failed: " + err.Error())
		return nil, err
	}
	return &node, nil
}

type NodeIds struct {
	NodeId string `gorm:"column:NodeId"`
}

func GetNodeIds() ([]*NodeIds, error) {
	sql := "SELECT node_id AS NodeId FROM " + common.TableNode
	var nodeIds []*NodeIds
	connection.DB.Raw(sql).Scan(&nodeIds)
	return nodeIds, nil
}

type NodeWithChainOrg struct {
	common.Node
	OrgId   string `json:"org_id"`
	OrgName string
}

func GetNodeListByChainId(chainId string, nodeName string, offset int, limit int) (int64, []*NodeWithChainOrg, error) {
	var (
		count        int64
		nodeList     []*NodeWithChainOrg
		err          error
		nodeSelector *gorm.DB
	)

	nodeSelector = connection.DB.Select("node.*, chain.org_name, chain.org_id").Table(common.TableChainOrgNode+" chain").
		Joins("LEFT JOIN "+common.TableNode+" node on chain.node_name = node.node_name").
		Where("chain.chain_id = ?", chainId)

	if nodeName != "" {
		if err = nodeSelector.Where("node.node_name = ?", nodeName).Find(&nodeList).Error; err != nil {
			log.Error("GetNodeListByChainId Failed: " + err.Error())
			return count, nodeList, err
		}
		return int64(len(nodeList)), nodeList, err
	}

	if err = nodeSelector.Count(&count).Error; err != nil {
		log.Error("GetNodeListByChainId Failed: " + err.Error())
		return count, nodeList, err
	}
	if err = nodeSelector.Offset(offset).Limit(limit).Find(&nodeList).Error; err != nil {
		log.Error("GetNodeListByChainId Failed: " + err.Error())
		return count, nodeList, err
	}

	return count, nodeList, err
}

func GetNodeInfo(chainId string, nodeId int) (NodeWithChainOrg, error) {
	var (
		nodeInfo NodeWithChainOrg
		err      error
	)

	err = connection.DB.Select("node.*, chain.org_name, chain.org_id").Table(common.TableChainOrgNode+" chain").
		Joins("LEFT JOIN "+common.TableNode+" node on chain.node_name = node.node_name").
		Where("chain.chain_id = ? and node.id = ?", chainId, nodeId).Find(&nodeInfo).Error

	if err != nil {
		log.Error("GetNodeInfo Failed: " + err.Error())
	}

	return nodeInfo, err
}

func GetLinkNodeList(chainId string, nodeId int) []*NodeWithChainOrg {
	var (
		nodeList []*NodeWithChainOrg
		err      error
	)
	err = connection.DB.Select("node.*, chain.org_name, chain.org_id").Table(common.TableChainOrgNode+" chain").
		Joins("LEFT JOIN "+common.TableNode+" node on chain.node_name = node.node_name").
		Where("chain.chain_id = ? and node.id != ?", chainId, nodeId).Find(&nodeList).Error

	if err != nil {
		log.Error("GetLinkNodeList Failed: " + err.Error())
	}
	return nodeList
}

func GetConsensusNodeNameList(chainId string) []string {
	// 获取某个交易的共识节点列表。
	// 这个方法应该不对，如果共识节点有变更，该方法只能获取当前的共识节点，获取不了历史上某个块或者交易当时的共识节点列表
	// 应该去查询链上当前交易发生时的链状态数据，以获取准确数据
	var (
		nodeList []string
		err      error
		orgList  []common.Org
	)

	err = connection.DB.Table(common.TableChainOrgNode+" org_node").Select("org_node.org_name").
		Joins("LEFT JOIN "+common.TableNode+" node on org_node.node_name = node.node_name").
		Where("org_node.chain_id = ? and node.type = 0", chainId).
		Scan(&orgList).Error

	if err != nil {
		log.Error("GetConsensusNodeNameList Failed: " + err.Error())
	}
	for _, org := range orgList {
		nodeList = append(nodeList, org.OrgName)
	}
	return nodeList
}
