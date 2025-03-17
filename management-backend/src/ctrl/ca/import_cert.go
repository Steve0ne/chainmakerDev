/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/entity"
	"management_backend/src/utils"
)

type ImportCertHandler struct{}

func (importCertHandler *ImportCertHandler) LoginVerify() bool {
	return true
}

func (importCertHandler *ImportCertHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindImportCertHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	caCertKey := params.CaCert
	caKey := params.CaKey
	orgId := params.OrgId
	orgName := params.OrgName
	userName := params.UserName
	signPrivKey := params.SignKey
	signCertKey := params.SignCert
	tlsPrivKey := params.TlsKey
	tlsCertKey := params.TlsCert

	if params.Type == ORG_CERT {
		err := checkOrgCert(caCertKey, caKey, ctx)
		if err != nil {
			log.Error("Check sign cert err : " + err.Error())
			return
		}
		err = importOrgCert(orgId, orgName, caCertKey, caKey, userName, ctx)
		if err != nil {
			return
		}
	} else {
		err := checkCert(signCertKey, signPrivKey, orgId, ctx)
		if err != nil {
			log.Error("Check sign cert err : " + err.Error())
			return
		}
		err = checkCert(tlsCertKey, tlsPrivKey, orgId, ctx)
		if err != nil {
			log.Error("Check tls cert err : " + err.Error())
			return
		}
		err = importUserAndNodeCert(params, ctx)
		if err != nil {
			return
		}
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

func importOrgCert(orgId, orgName, caCertKey, caKey, userName string, ctx *gin.Context) error {
	count, err := chain_participant.GetOrgCaCertCountBydOrgIdAndOrgName(orgId, orgName)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		log.Error("GetOrgCaCertCountBydOrgIdAndOrgName err : " + err.Error())
		return err
	}
	if count > 0 {
		common.ConvergeFailureResponse(ctx, common.ErrorCertExisted)
		log.Error("orgCa has generated")
		return errors.New("orgCa has generated")
	}
	err = saveUploadCert(caKey, caCertKey, orgId, orgName, userName, "", "", SIGN_CERT)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		log.Error("SaveUploadCert err : " + err.Error())
		return err
	}
	org := &dbcommon.Org{
		OrgId:   orgId,
		OrgName: orgName,
	}
	err = chain_participant.CreateOrg(org)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		log.Error("CreateOrg err : " + err.Error())
		return err
	}
	return nil
}

func importUserAndNodeCert(params *ImportCertParams, ctx *gin.Context) error {
	//certType通过证书解析出来
	tlsPrivKey := params.TlsKey
	tlsCertKey := params.TlsCert
	signPrivKey := params.SignKey
	signCertKey := params.SignCert
	orgId := params.OrgId
	orgName := params.OrgName
	userName := params.UserName

	if params.Type == NODE_CERT {
		count, err := chain_participant.GetNodeCertCount(params.NodeName)
		if err != nil {
			log.Error("GetNodeCertCount err : " + err.Error())
			common.ConvergeHandleFailureResponse(ctx, err)
			return err
		}
		if count > 0 {
			log.Error("nodeName has existed")
			common.ConvergeFailureResponse(ctx, common.ErrorCertExisted)
			return errors.New("nodeName has existed")
		}

		regexpStr := "^(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\." +
			"(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\:" +
			"(6553[0-5]|655[0-2]\\d|65[0-4]\\d{2}|6[0-4]\\d{3}|[0-5]\\d{4}|[1-9]\\d{0,3})$"
		if ok, _ := regexp.MatchString(regexpStr, params.NodeIp); !ok {
			log.Error("ip format err")
			common.ConvergeFailureResponse(ctx, common.ErrorIpFormat)
			return errors.New("ip format err")
		}

	} else if params.Type == USER_CERT {
		count, err := chain_participant.GetUserCertCount(params.UserName)
		if err != nil {
			log.Error("GetUserCertCount err : " + err.Error())
			common.ConvergeHandleFailureResponse(ctx, err)
			return err
		}
		if count > 0 {
			log.Error("userName has existed")
			common.ConvergeFailureResponse(ctx, common.ErrorCertExisted)
			return errors.New("userName has existed")
		}
	}
	err := saveUploadCert(tlsPrivKey, tlsCertKey, orgId, orgName, userName, params.NodeName, params.NodeIp, TLS_CERT)
	if err != nil {
		log.Error("SaveUploadCert err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	err = saveUploadCert(signPrivKey, signCertKey, orgId, orgName, userName, params.NodeName, params.NodeIp, SIGN_CERT)
	if err != nil {
		log.Error("SaveUploadCert err : " + err.Error())
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	return nil
}

func checkOrgCert(certKey, privKey string, ctx *gin.Context) error {
	certId, certUserId, certHash, err := ResolveUploadKey(certKey)
	if err != nil {
		return err
	}

	certUpload, err := db.GetUploadByIdAndUserIdAndHash(certId, certUserId, certHash)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	certInfo, err := utils.ParseCertificate(certUpload.Content)
	if err != nil {
		//证书格式错误
		common.ConvergeFailureResponse(ctx, common.ErrorCertContent)
		return err
	}

	var keyContent []byte

	privKeyId, privKeyUserId, privKeyHash, err := ResolveUploadKey(privKey)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	privUpload, err := db.GetUploadByIdAndUserIdAndHash(privKeyId, privKeyUserId, privKeyHash)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	keyContent = privUpload.Content

	publicKey, ok := certInfo.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		common.ConvergeHandleFailureResponse(ctx, errors.New("this cert dosen't have publicKey"))
		return errors.New("this cert dosen't have publicKey")
	}

	privateKey, err := utils.ParsePrivateKey(keyContent)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	if !publicKey.Equal(privateKey.PublicKey().ToStandardKey()) {
		common.ConvergeFailureResponse(ctx, common.ErrorCertKeyMatch)
		return errors.New("this cert dosen't match key")
	}

	return nil
}

func checkCert(certKey, privKey, orgID string, ctx *gin.Context) error {

	orgCert, err := chain_participant.GetOrgCaCert(orgID)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorOrgNoExisted)
		return err
	}

	certId, certUserId, certHash, err := ResolveUploadKey(certKey)
	if err != nil {
		return err
	}

	certUpload, err := db.GetUploadByIdAndUserIdAndHash(certId, certUserId, certHash)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	orgCertInfo, err := utils.ParseCertificate([]byte(orgCert.Cert))
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}

	certInfo, err := utils.ParseCertificate(certUpload.Content)
	if err != nil {
		//证书格式错误
		common.ConvergeFailureResponse(ctx, common.ErrorCertContent)
		return err
	}
	authKeyId := certInfo.AuthorityKeyId

	var keyContent []byte

	privKeyId, privKeyUserId, privKeyHash, err := ResolveUploadKey(privKey)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	privUpload, err := db.GetUploadByIdAndUserIdAndHash(privKeyId, privKeyUserId, privKeyHash)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	keyContent = privUpload.Content

	publicKey, ok := certInfo.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		common.ConvergeHandleFailureResponse(ctx, errors.New("this cert dosen't have publicKey"))
		return errors.New("this cert dosen't have publicKey")
	}

	privateKey, err := utils.ParsePrivateKey(keyContent)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return err
	}
	if !publicKey.Equal(privateKey.PublicKey().ToStandardKey()) {
		common.ConvergeFailureResponse(ctx, common.ErrorCertKeyMatch)
		return errors.New("this cert dosen't match key")
	}

	if !bytes.Equal(orgCertInfo.SubjectKeyId, authKeyId) {
		common.ConvergeFailureResponse(ctx, common.ErrorIssueOrg)
		return errors.New("this cert dosen't Issue by this org")
	}
	return nil
}
