/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"github.com/emirpasic/gods/lists/arraylist"

	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/relation"
)

type CertNodeListView struct {
	NodeName string
}

func NewCertNodeListView(orgNodes []*dbcommon.OrgNode) []interface{} {
	nodeListView := arraylist.New()
	for _, orgNode := range orgNodes {
		certNodeView := CertNodeListView{
			NodeName: orgNode.NodeName,
		}
		nodeListView.Add(certNodeView)
	}

	return nodeListView.Values()
}

type CertOrgListView struct {
	OrgId   string
	OrgName string
}

func NewCertOrgListView(orgs []*dbcommon.Org) []interface{} {
	orgListView := arraylist.New()
	for _, org := range orgs {
		certOrgView := CertOrgListView{
			OrgId:   org.OrgId,
			OrgName: org.OrgName,
		}
		orgListView.Add(certOrgView)
	}

	return orgListView.Values()
}

func NewCertOrgListByChainIdView(chainOrg []*dbcommon.ChainOrg) []interface{} {
	orgListView := arraylist.New()
	for _, org := range chainOrg {
		certOrgView := CertOrgListView{
			OrgId:   org.OrgId,
			OrgName: org.OrgName,
		}
		orgListView.Add(certOrgView)
	}

	return orgListView.Values()
}

type CertUserListView struct {
	UserName string
}

func NewCertUserListView(certs []*dbcommon.Cert) []interface{} {
	certUserListView := arraylist.New()
	for _, cert := range certs {
		certUserView := CertUserListView{
			UserName: cert.CertUserName,
		}
		certUserListView.Add(certUserView)
	}

	return certUserListView.Values()
}

type ChainListView struct {
	Id         int64
	ChainName  string
	ChainId    string
	CreateTime int64
	OrgNum     int
	NodeNum    int
	AutoReport int
	Status     int
	Monitor    int
}

func NewChainListView(chains []*dbcommon.Chain) []interface{} {
	chainListView := arraylist.New()
	for _, chain := range chains {
		orgNum, _ := relation.GetOrgCountByChainId(chain.ChainId)
		nodeNum, _ := relation.GetNodeCountByChainId(chain.ChainId)
		chainsView := ChainListView{
			Id:         chain.Id,
			ChainName:  chain.ChainName,
			ChainId:    chain.ChainId,
			CreateTime: chain.CreatedAt.Unix(),
			OrgNum:     orgNum,
			NodeNum:    nodeNum,
			AutoReport: chain.AutoReport,
			Status:     chain.Status,
			Monitor:    chain.Monitor,
		}
		chainListView.Add(chainsView)
	}

	return chainListView.Values()
}

type DownloadZipView struct {
	File     string
	FileName string
}
