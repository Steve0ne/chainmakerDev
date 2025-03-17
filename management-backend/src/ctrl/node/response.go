/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package node

import (
	"management_backend/src/db/chain_participant"
)

type NodeView struct {
	Id           int64
	OrgName      string
	OrgId        string
	NodeName     string
	NodeId       string
	NodeType     int
	NodeAddr     string
	NodePort     string
	UpdateType   string
	CreateTime   int64
	LinkNodeList []LinkNode
}

type LinkNode struct {
	LinkNodeName string
	LinkNodeType int
}

const fullNode = "FULL"

func NewNodeView(node *chain_participant.NodeWithChainOrg) *NodeView {
	return &NodeView{
		Id:         node.Id,
		OrgId:      node.OrgId,
		OrgName:    node.OrgName,
		NodeName:   node.NodeName,
		NodeType:   node.Type,
		NodeId:     node.NodeId,
		NodeAddr:   node.NodeIp,
		NodePort:   node.NodePort,
		UpdateType: fullNode,
	}
}

func NewNodeViewWithLinkNode(node chain_participant.NodeWithChainOrg, nodeList []LinkNode) *NodeView {
	return &NodeView{
		Id:           node.Id,
		OrgId:        node.OrgId,
		OrgName:      node.OrgName,
		NodeName:     node.NodeName,
		NodeType:     node.Type,
		NodeId:       node.NodeId,
		NodeAddr:     node.NodeIp,
		NodePort:     node.NodePort,
		UpdateType:   fullNode,
		CreateTime:   node.CreatedAt.Unix(),
		LinkNodeList: nodeList,
	}
}
