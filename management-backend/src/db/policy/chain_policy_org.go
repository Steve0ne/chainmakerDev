/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package policy

import (
	"management_backend/src/db/common"
	"management_backend/src/db/connection"
)

func SaveChainPolicyOrg(chainPolicyOrg *common.ChainPolicyOrg) error {
	if err := connection.DB.Create(&chainPolicyOrg).Error; err != nil {
		log.Error("Save chainPolicyOrg Failed: " + err.Error())
		return err
	}
	return nil
}

func GetUpdateColumns(chainPolicyOrg *common.ChainPolicyOrg) map[string]interface{} {
	columns := make(map[string]interface{})
	columns["status"] = chainPolicyOrg.Status
	return columns
}

func GetOrgListByPolicyType(chainId string, opType int) ([]*common.ChainPolicyOrg, error) {
	var (
		orgs []*common.ChainPolicyOrg
		err  error
	)

	orgSelector := connection.DB.Select("policy_org.org_id, org.org_name, policy_org.status").
		Table(common.TableChainPolicy+" policy").
		Joins("LEFT JOIN "+common.TableChainPolicyOrg+" policy_org on policy.id = policy_org.chain_policy_id").
		Joins("LEFT JOIN "+common.TableOrg+" org on policy_org.org_id = org.org_id").
		Where("policy.chain_id = ? and policy.type = ?", chainId, opType)

	if err = orgSelector.Find(&orgs).Error; err != nil {
		log.Error("GetOrgListByPolicyType Failed: " + err.Error())
		return orgs, err
	}
	return orgs, nil

}
