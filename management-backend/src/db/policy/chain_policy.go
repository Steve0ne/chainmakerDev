/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package policy

import (
	"strconv"
	"strings"

	"management_backend/src/db/common"
	"management_backend/src/db/connection"
	loggers "management_backend/src/logger"
)

var (
	log = loggers.GetLogger(loggers.ModuleDb)
)

const (
	Majority = iota
	Any
	Self
	All
	Forbidden
	Percentage
)

const (
	AdminRole = iota
	ClientRole
	AllRole
)
const roleSelected = 1

func CreateChainPolicy(chainPolicy *common.ChainPolicy, chainPolicyOrgList []*common.ChainPolicyOrg) error {

	dbChainPolicy, err := GetChainPolicy(chainPolicy.ChainId, chainPolicy.Type)
	if err == nil {
		dbChainPolicy.PolicyType = chainPolicy.PolicyType
		dbChainPolicy.RoleType = chainPolicy.RoleType
		dbChainPolicy.PercentNum = chainPolicy.PercentNum
		if err = connection.DB.Debug().Model(chainPolicy).Where("id = ?", dbChainPolicy.Id).
			UpdateColumns(getUpdateColumns(dbChainPolicy)).Error; err != nil {
			log.Error("UpdateChainPolicy failed: " + err.Error())
			return err
		}
		for _, chainPolicyOrg := range chainPolicyOrgList {
			err = connection.DB.Debug().Model(chainPolicyOrg).
				Where("chain_policy_id = ? AND org_id = ?", dbChainPolicy.Id, chainPolicyOrg.OrgId).
				UpdateColumns(GetUpdateColumns(chainPolicyOrg)).Error
			if err != nil {
				log.Error("UpdateChainPolicyOrg failed: " + err.Error())
				return err
			}
		}
	} else {
		err = SaveChainPolicy(chainPolicy)
		for _, chainPolicyOrg := range chainPolicyOrgList {
			chainPolicyOrg.ChainPolicyId = chainPolicy.Id
			err = SaveChainPolicyOrg(chainPolicyOrg)
		}
		return err
	}

	return nil
}

func GetPassedVoteCnt(orgCnt int, p *common.ChainPolicy) int {
	//策略类型 0:Majority; 1:Any; 2:Self; 3:All 4:Forbidden; 5:percentage
	switch p.PolicyType {
	case Majority: // Majority
		// 大于，不包括等于
		return (orgCnt / 2) + 1
	case Any: // Any
		return 1
	case Self: // Self
		return 1
	case All: // All
		return orgCnt
	case Forbidden: // Forbidden
		return 0
	case Percentage: // percentage
		// 大于或者等于的条件均满足
		sep := "/"
		arr := strings.Split(p.PercentNum, sep)
		if len(arr) == 1 {
			i, err := strconv.Atoi(arr[0])
			if err != nil {
				return -1
			}
			return i

		}

		numerator, err := strconv.Atoi(arr[0])
		if err != nil {
			return -1
		}
		denominator, err := strconv.Atoi(arr[1])
		if err != nil {
			return -1
		}
		percent := float64(numerator) / float64(denominator)
		needCnt := percent * float64(orgCnt)
		if float64(int(needCnt)) < needCnt {
			return int(needCnt) + 1
		}
		return int(needCnt)
	}
	return -1
}

func getUpdateColumns(chainPolicy *common.ChainPolicy) map[string]interface{} {
	columns := make(map[string]interface{})
	columns["type"] = chainPolicy.Type
	columns["policy_type"] = chainPolicy.PolicyType
	columns["role_type"] = chainPolicy.RoleType
	return columns
}

func SaveChainPolicy(chainPolicy *common.ChainPolicy) error {
	if err := connection.DB.Create(&chainPolicy).Error; err != nil {
		log.Error("Save chainPolicy Failed: " + err.Error())
		return err
	}
	return nil
}

func GetChainPolicy(chainId string, chainPolicyType int) (*common.ChainPolicy, error) {
	var chainPolicy common.ChainPolicy
	err := connection.DB.Where("chain_id = ? AND type = ?", chainId, chainPolicyType).Find(&chainPolicy).Error
	if err != nil {
		log.Error("GetChainPolicy Failed: " + err.Error())
		return nil, err
	}
	return &chainPolicy, nil
}

func GetChainPolicyByChainId(chainId string) ([]*common.ChainPolicy, error) {
	var chainPolicy []*common.ChainPolicy
	if err := connection.DB.Where("chain_id = ?", chainId).Find(&chainPolicy).Error; err != nil {
		log.Error("GetChainPolicyByChainId Failed: " + err.Error())
		return nil, err
	}
	return chainPolicy, nil
}

type UserRole struct {
	Role     int
	Selected int
}

func GetRoleList(chainId string, opType int) ([]*UserRole, error) {
	var (
		policy         common.ChainPolicy
		err            error
		adminSelected  int
		clientSelected int
		roleList       []*UserRole
	)

	err = connection.DB.Model(&common.ChainPolicy{}).
		Select("role_type").
		Where("chain_id = ? and type = ?", chainId, opType).
		Find(&policy).Error

	if err != nil {
		log.Error("GetRoleList Failed: " + err.Error())
		return nil, err
	}

	switch policy.RoleType {
	case AdminRole:
		adminSelected = roleSelected
	case ClientRole:
		clientSelected = roleSelected
	case AllRole:
		adminSelected = roleSelected
		clientSelected = roleSelected
	}

	var (
		admin  UserRole
		client UserRole
	)
	// 构建admin
	admin = UserRole{
		Role:     AdminRole,
		Selected: adminSelected,
	}
	client = UserRole{
		Role:     ClientRole,
		Selected: clientSelected,
	}
	roleList = []*UserRole{&admin, &client}
	return roleList, nil
}
