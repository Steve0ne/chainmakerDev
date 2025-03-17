/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_participant

import (
	"management_backend/src/db/common"
	"management_backend/src/db/connection"
)

func CreateOrg(org *common.Org) error {
	if err := connection.DB.Create(&org).Error; err != nil {
		log.Error("[DB] Save org Failed: " + err.Error())
		return err
	}
	return nil
}

func GetOrgByOrgId(orgId string) (*common.Org, error) {
	var org common.Org
	if err := connection.DB.Where("org_id = ?", orgId).Find(&org).Error; err != nil {
		log.Error("QueryOrgByOrgId Failed: " + err.Error())
		return nil, err
	}
	return &org, nil
}

func GetOrgNameByOrgId(orgId string) (string, error) {
	var org common.Org
	if err := connection.DB.Select("org_name").Where("org_id = ?", orgId).Find(&org).Error; err != nil {
		log.Error("GetOrgNameByOrgId Failed: " + err.Error())
		return "", err
	}
	return org.OrgName, nil
}

func GetOrgList() ([]*common.Org, error) {
	var orgs []*common.Org
	if err := connection.DB.Find(&orgs).Error; err != nil {
		log.Error("GetOrgList Failed: " + err.Error())
		return nil, err
	}
	return orgs, nil
}

type OrgIds struct {
	OrgId string `gorm:"column:OrgId"`
}

func GetOrgIds() ([]*OrgIds, error) {
	sql := "SELECT org_id AS OrgId FROM " + common.TableOrg
	var orgIds []*OrgIds
	connection.DB.Raw(sql).Scan(&orgIds)
	return orgIds, nil
}
