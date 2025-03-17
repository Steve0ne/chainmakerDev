/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package overview

import (
	"github.com/emirpasic/gods/lists/arraylist"

	"management_backend/src/ctrl/contract_management"
	dbcommon "management_backend/src/db/common"
)

type GeneralDataView struct {
	TxNum       int64 `json:"TxNum"`
	BlockHeight int64 `json:"BlockHeight"`
	NodeNum     int   `json:"NodeNum"`
	OrgNum      int   `json:"OrgNum"`
	ContractNum int64 `json:"ContractNum"`
}

type AuthListView struct {
	Type       int
	PolicyType int
	PercentNum string
}

func NewAuthListView(authList []*dbcommon.ChainPolicy) []interface{} {
	authViews := arraylist.New()
	for _, auth := range authList {
		if auth.Type != contract_management.PERMISSION_UPDATE {
			authView := AuthListView{
				Type:       auth.Type,
				PolicyType: auth.PolicyType,
				PercentNum: auth.PercentNum,
			}
			authViews.Add(authView)
		}
	}
	return authViews.Values()
}

type PolicyOrgView struct {
	OrgName  string `json:"OrgName"`
	OrgId    string `json:"OrgId"`
	Selected int    `json:"Selected"`
}

func NewPolicyOrgView(org *dbcommon.ChainPolicyOrg) *PolicyOrgView {
	return &PolicyOrgView{
		OrgName:  org.OrgName,
		OrgId:    org.OrgId,
		Selected: org.Status,
	}
}

type ChainView struct {
	Id              int64
	ChainId         string
	ChainName       string
	Version         string
	Sequence        string
	BlockTxCapacity uint32
	TxTimeout       uint32
	BlockInterval   uint32
	DockerVm        int
}

func NewChainView(chain *dbcommon.Chain) *ChainView {
	chainView := ChainView{
		Id:              chain.Id,
		ChainId:         chain.ChainId,
		ChainName:       chain.ChainName,
		Version:         chain.Version,
		Sequence:        chain.Sequence,
		BlockTxCapacity: chain.BlockTxCapacity,
		TxTimeout:       chain.TxTimeout,
		BlockInterval:   chain.BlockInterval,
		DockerVm:        chain.DockerVm,
	}
	return &chainView
}
