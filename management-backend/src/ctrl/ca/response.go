/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"management_backend/src/db/chain_participant"

	dbcommon "management_backend/src/db/common"
)

// 这里的用来区分请求
const (
	CERT_FOR_SIGN = iota
	KEY_FOR_SIGN
	CERT_FOR_TLS
	KEY_FOR_TLS
)

const (
	SIGNUSE = iota
	TLSUSE
)

type CertDetailView struct {
	CertDetail string
}

func NewCertDetailView(certDetail string) *CertDetailView {

	certDetailView := CertDetailView{
		CertDetail: certDetail,
	}

	return &certDetailView
}

type CertListView struct {
	Id         int64
	UserName   string
	OrgName    string
	NodeName   string
	CertUse    int
	CertType   int
	NodeIp     string
	NodePort   string
	CreateTime int64
}

func NewCertListView(certs []*dbcommon.Cert) []interface{} {
	certViews := arraylist.New()
	for _, cert := range certs {
		var keyUse int
		var certUse int
		var nodeIp string
		var nodePort string

		if cert.CertType == chain_participant.CONSENSUS || cert.CertType == chain_participant.COMMON {
			node, err := chain_participant.GetNodeByNodeName(cert.NodeName)
			if err == nil {
				nodeIp = node.NodeIp
				nodePort = node.NodePort
			}
		}

		if cert.CertUse == SIGNUSE {
			keyUse = KEY_FOR_SIGN
			certUse = CERT_FOR_SIGN
		} else if cert.CertUse == TLSUSE {
			keyUse = KEY_FOR_TLS
			certUse = CERT_FOR_TLS
		}
		certView := CertListView{
			Id:         cert.Id,
			UserName:   cert.CertUserName,
			OrgName:    cert.OrgName,
			NodeName:   cert.NodeName,
			CertUse:    certUse,
			CertType:   cert.CertType,
			NodeIp:     nodeIp,
			NodePort:   nodePort,
			CreateTime: cert.CreatedAt.Unix(),
		}

		keyView := CertListView{
			Id:         cert.Id,
			UserName:   cert.CertUserName,
			OrgName:    cert.OrgName,
			NodeName:   cert.NodeName,
			CertUse:    keyUse,
			CertType:   cert.CertType,
			NodeIp:     nodeIp,
			NodePort:   nodePort,
			CreateTime: cert.CreatedAt.Unix(),
		}
		certViews.Add(keyView)
		certViews.Add(certView)
	}

	return certViews.Values()
}
