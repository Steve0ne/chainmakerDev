/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"errors"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/helper"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/relation"
	"management_backend/src/entity"
	"management_backend/src/utils"
)

type GenerateCertHandler struct{}

func (generateCertHandler *GenerateCertHandler) LoginVerify() bool {
	return true
}

func (generateCertHandler *GenerateCertHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGenerateCertHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	orgId := params.OrgId
	orgName := params.OrgName
	nodeName := params.NodeName
	userName := params.UserName

	if params.CertType == ORG_CERT {
		err := generateOrgCert(orgId, orgName, userName, nodeName, COUNTRY, LOCALITY, PROVINCE, ctx)
		if err != nil {
			return
		}
	} else if params.CertType == NODE_CERT {
		regexpStr := "^(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\." +
			"(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\:" +
			"(6553[0-5]|655[0-2]\\d|65[0-4]\\d{2}|6[0-4]\\d{3}|[0-5]\\d{4}|[1-9]\\d{0,3})$"
		if ok, _ := regexp.MatchString(regexpStr, params.NodeIp); !ok {
			log.Error("ip format err")
			common.ConvergeFailureResponse(ctx, common.ErrorIpFormat)
			return
		}
		err := generateNodeCert(orgId, orgName, userName, nodeName, COUNTRY, LOCALITY, PROVINCE, params.NodeIp,
			params.NodeRole, ctx)
		if err != nil {
			return
		}

	} else if params.CertType == USER_CERT {
		err := generateUserCert(orgId, orgName, userName, nodeName, COUNTRY, LOCALITY, PROVINCE, params.UserRole, ctx)
		if err != nil {
			return
		}
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

func generateOrgCert(orgId, orgName, userName, nodeName, country, locality, province string, ctx *gin.Context) error {
	count, err := chain_participant.GetOrgCaCertCountBydOrgIdAndOrgName(orgId, orgName)
	if err != nil {
		log.Error("ErrorCreateKey err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	if count > 0 {
		log.Error("orgCa has generated")
		common.ConvergeFailureResponse(ctx, common.ErrorCertExisted)
		return errors.New("orgCa has generated")
	}
	privKey, privKeyStr, err := createPrivKey()
	if err != nil {
		log.Error("ErrorCreateKey err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	certType := chain_participant.ORG_CA
	cACertificateConfig := &utils.CACertificateConfig{
		PrivKey:            privKey,
		HashType:           crypto.HASH_TYPE_SHA256,
		Country:            country,
		Locality:           locality,
		Province:           province,
		OrganizationalUnit: ORG_OU,
		Organization:       orgId,
		CommonName:         "ca." + orgId,
		ExpireYear:         EXPIREYEAR,
		Sans:               sans,
	}
	certPem, err := utils.CreateCACertificate(cACertificateConfig)
	if err != nil {
		log.Error("createCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	org := &dbcommon.Org{
		OrgId:   orgId,
		OrgName: orgName,
	}
	err = chain_participant.CreateOrg(org)
	if err != nil {
		log.Error("CreateOrg err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	err = saveCert(privKeyStr, certPem, SIGN_CERT, certType, orgId, orgName, userName, nodeName)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	return nil
}

func generateNodeCert(orgId, orgName, userName, nodeName, country, locality, province, ip string, nodeRole int,
	ctx *gin.Context) error {
	caCount, err := chain_participant.GetOrgCaCertCount(orgId)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	if caCount < 1 {
		common.ConvergeFailureResponse(ctx, common.ErrorOrgNoExisted)
		return errors.New("orgCa no existed")
	}

	var certPem string
	var tlsCertPem string
	var certType int
	var ou string
	var privKeyStr string
	var tlsPrivKeyStr string

	count, err := chain_participant.GetNodeCertCount(nodeName)
	if err != nil {
		log.Error("GetNodeCertCount err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	if count > 0 {
		log.Error("nodeCert has generated")
		common.ConvergeFailureResponse(ctx, common.ErrorCertExisted)
		return errors.New("orgCa has generated")
	}
	var nodeType int
	if nodeRole == CONSENSUS_NODE_ROLE {
		certType = chain_participant.CONSENSUS
		ou = CONSENSUS_NODE_OU
		nodeType = NODE_CONSENSUS
	} else if nodeRole == COMMON_NODE_ROLE {
		certType = chain_participant.COMMON
		ou = COMMON_NODE_OU
		nodeType = NODE_COMMON
	}
	certPem, privKeyStr, err = IssueCert(country, locality, province, ou, orgId, nodeName+".sign."+orgId)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	tlsCertPem, tlsPrivKeyStr, err = IssueCert(country, locality, province, ou, orgId, nodeName+".tls."+orgId)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	err = saveCert(tlsPrivKeyStr, tlsCertPem, TLS_CERT, certType, orgId, orgName, userName, nodeName)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	nodeId, err := helper.GetLibp2pPeerIdFromCert([]byte(tlsCertPem))
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	node := &dbcommon.Node{
		NodeId:   nodeId,
		NodeName: nodeName,
		NodeIp:   ip[:strings.Index(ip, ":")],
		NodePort: ip[strings.Index(ip, ":")+1:],
		Type:     nodeType,
	}
	err = chain_participant.CreateNode(node)
	if err != nil {
		log.Error("CreateNode err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
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
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	err = saveCert(privKeyStr, certPem, SIGN_CERT, certType, orgId, orgName, userName, nodeName)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	return nil
}

func generateUserCert(orgId, orgName, userName, nodeName, country, locality, province string, userRole int,
	ctx *gin.Context) error {
	caCount, err := chain_participant.GetOrgCaCertCount(orgId)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	if caCount < 1 {
		common.ConvergeFailureResponse(ctx, common.ErrorOrgNoExisted)
		return errors.New("orgCa no existed")
	}
	count, err := chain_participant.GetUserCertCount(userName)
	if err != nil {
		log.Error("ErrorCreateKey err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	if count > 0 {
		log.Error("userCert has generated")
		common.ConvergeFailureResponse(ctx, common.ErrorCertExisted)
		return errors.New("orgCa has generated")
	}

	var certPem string
	var tlsCertPem string
	var certType int
	var ou string
	var privKeyStr string
	var tlsPrivKeyStr string

	if userRole == ADMIN_USER_ROLE {
		certType = chain_participant.ADMIN
		ou = ADMIN_USER_OU
	} else if userRole == CLIENT_USER_ROLE {
		certType = chain_participant.CLIENT
		ou = CLIENT_USER_OU
	}
	certPem, privKeyStr, err = IssueCert(country, locality, province, ou, orgId, userName+".sign."+orgId)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	tlsCertPem, tlsPrivKeyStr, err = IssueCert(country, locality, province, ou, orgId, userName+".tls."+orgId)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	err = saveCert(tlsPrivKeyStr, tlsCertPem, TLS_CERT, certType, orgId, orgName, userName, nodeName)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	err = saveCert(privKeyStr, certPem, SIGN_CERT, certType, orgId, orgName, userName, nodeName)
	if err != nil {
		log.Error("CreateCACertificate err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	return nil
}
