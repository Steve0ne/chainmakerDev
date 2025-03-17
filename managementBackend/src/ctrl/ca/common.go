/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"strings"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/helper"
	"management_backend/src/db"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/relation"
	loggers "management_backend/src/logger"
	"management_backend/src/utils"
)

//OU字段
const (
	ORG_OU            = "root-cert"
	CONSENSUS_NODE_OU = "consensus"
	COMMON_NODE_OU    = "common"
	ADMIN_USER_OU     = "admin"
	CLIENT_USER_OU    = "client"
)

//证书类型
const (
	ORG_CERT  = 0
	NODE_CERT = 1
	USER_CERT = 2
)

//证书角色
const (
	CONSENSUS_NODE_ROLE = 0
	COMMON_NODE_ROLE    = 1
	ADMIN_USER_ROLE     = 0
	CLIENT_USER_ROLE    = 1
)

//证书用途
const (
	SIGN_CERT = 0
	TLS_CERT  = 1
)

//节点类型
const (
	NODE_CONSENSUS = 0
	NODE_COMMON    = 1
)

//证书属性
const (
	EXPIREYEAR = 9
	TLS_HOST   = "chainmaker.org"
)

const (
	COUNTRY  = "cn"
	LOCALITY = "beijing"
	PROVINCE = "beijing"
)

var (
	sans = []string{"127.0.0.1", "localhost", "chainmaker.org"}
	log  = loggers.GetLogger(loggers.ModuleWeb)
)

func createPrivKey() (crypto.PrivateKey, string, error) {
	privKey, err := utils.CreatePrivKey(crypto.ECC_NISTP256)
	if err != nil {
		return nil, "", err
	}

	privKeyStr, err := privKey.String()
	if err != nil {
		return nil, "", err
	}

	return privKey, privKeyStr, nil
}

func IssueCert(country, locality, province, ou, orgId, cn string) (string, string, error) {

	privKey, privKeyStr, err := createPrivKey()
	if err != nil {
		return "", "", err
	}
	csrConfig := &utils.CSRConfig{
		PrivKey:            privKey,
		Country:            country,
		Locality:           locality,
		Province:           province,
		OrganizationalUnit: ou,
		Organization:       orgId,
		CommonName:         cn,
	}
	csrPem, err := utils.CreateCSR(csrConfig)
	if err != nil {
		return "", "", err
	}
	csr, err := utils.ParseCsr([]byte(csrPem))
	if err != nil {
		return "", "", err
	}

	orgCaCert, err := chain_participant.GetOrgCaCert(orgId)
	if err != nil {
		return "", "", err
	}

	certInfo, err := utils.ParseCertificate([]byte(orgCaCert.Cert))
	if err != nil {
		return "", "", err
	}
	pkInfo, err := utils.ParsePrivateKey([]byte(orgCaCert.PrivateKey))
	if err != nil {
		return "", "", err
	}

	issueCertificateConfig := &utils.IssueCertificateConfig{
		HashType:         crypto.HASH_TYPE_SHA256,
		IsCA:             false,
		IssuerPrivKeyPwd: nil,
		ExpireYear:       EXPIREYEAR,
		Sans:             sans,
		Uuid:             "",
		PrivKey:          pkInfo,
		IssuerCert:       certInfo,
		Csr:              csr,
	}
	var certPem string
	certPem, err = utils.IssueCertificate(issueCertificateConfig)
	if err != nil {
		return "", "", err
	}
	return certPem, privKeyStr, nil
}

func saveCert(privKeyStr, certPemStr string, certUse int, certType int, orgId, orgName, userName,
	nodeName string) error {
	certInfo := &dbcommon.Cert{
		Cert:         certPemStr,
		PrivateKey:   privKeyStr,
		CertType:     certType,
		CertUse:      certUse,
		OrgId:        orgId,
		OrgName:      orgName,
		CertUserName: userName,
		NodeName:     nodeName,
	}
	err := chain_participant.CreateCert(certInfo)
	if err != nil {
		return err
	}
	return nil
}

func saveUploadCert(privKey, certKey, orgId, orgName, userName, nodeName, nodeIp string, certUse int) error {
	certId, certUserId, certHash, err := ResolveUploadKey(certKey)
	if err != nil {
		return err
	}

	var privContent []byte

	if privKey != "" {
		var privKeyId, privKeyUserId int64
		var privKeyHash string
		privKeyId, privKeyUserId, privKeyHash, err = ResolveUploadKey(privKey)
		if err != nil {
			return err
		}
		privUpload, err := db.GetUploadByIdAndUserIdAndHash(privKeyId, privKeyUserId, privKeyHash)
		if err != nil {
			return err
		}
		privContent = privUpload.Content
	}

	certUpload, err := db.GetUploadByIdAndUserIdAndHash(certId, certUserId, certHash)
	if err != nil {
		return err
	}

	var certType int
	var nodeType int
	certInfo, err := utils.ParseCertificate(certUpload.Content)
	if err != nil {
		return err
	}
	if certInfo.Subject.OrganizationalUnit[0] == ADMIN_USER_OU {
		certType = chain_participant.ADMIN
	}
	if certInfo.Subject.OrganizationalUnit[0] == CLIENT_USER_OU {
		certType = chain_participant.CLIENT
	}
	if certInfo.Subject.OrganizationalUnit[0] == ORG_OU {
		certType = chain_participant.ORG_CA
	}
	if certInfo.Subject.OrganizationalUnit[0] == CONSENSUS_NODE_OU {
		certType = chain_participant.CONSENSUS
		nodeType = NODE_CONSENSUS
	}
	if certInfo.Subject.OrganizationalUnit[0] == COMMON_NODE_OU {
		certType = chain_participant.COMMON
		nodeType = NODE_COMMON
	}

	err = saveCert(string(privContent), string(certUpload.Content), certUse, certType, orgId, orgName, userName, nodeName)
	if err != nil {
		return err
	}

	if (certType == chain_participant.CONSENSUS || certType == chain_participant.COMMON) && certUse == TLS_CERT {
		nodeId, err := helper.GetLibp2pPeerIdFromCert(certUpload.Content)
		if err != nil {
			return err
		}
		node := &dbcommon.Node{
			NodeId:   nodeId,
			NodeName: nodeName,
			NodeIp:   nodeIp[:strings.Index(nodeIp, ":")],
			NodePort: nodeIp[strings.Index(nodeIp, ":")+1:],
			Type:     nodeType,
		}
		err = chain_participant.CreateNode(node)
		if err != nil {
			return err
		}
		orgNode := &dbcommon.OrgNode{
			NodeId:   nodeId,
			NodeName: nodeName,
			OrgId:    orgId,
			OrgName:  orgName,
		}
		err = relation.CreateOrgNode(orgNode)
		if err != nil {
			log.Error("CreateOrgNode err : " + err.Error())
			return err
		}
	}

	return nil
}
