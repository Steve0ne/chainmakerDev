/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package relation

import (
	"management_backend/src/db/common"
	"management_backend/src/db/connection"
)

func CreateOrgNode(orgNode *common.OrgNode) error {
	if err := connection.DB.Create(&orgNode).Error; err != nil {
		log.Error("Save orgNode Failed: " + err.Error())
		return err
	}
	return nil
}

func GetOrgNode(orgId string) ([]*common.OrgNode, error) {
	var orgNodes []*common.OrgNode

	db := connection.DB
	db = db.Where("org_id = ?", orgId)
	if err := db.Find(&orgNodes).Error; err != nil {
		log.Error("QueryOrgCaCert Failed: " + err.Error())
		return nil, err
	}
	return orgNodes, nil
}

func GetOrgNodeByNodeId(nodeId string) ([]*common.OrgNode, error) {
	var orgNodes []*common.OrgNode

	db := connection.DB
	db = db.Where("node_id = ?", nodeId)
	if err := db.Find(&orgNodes).Error; err != nil {
		log.Error("GetOrgNodeByNodeId Failed: " + err.Error())
		return nil, err
	}
	return orgNodes, nil
}
