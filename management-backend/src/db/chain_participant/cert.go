/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_participant

import (
	"management_backend/src/db/common"
	"management_backend/src/db/connection"
	loggers "management_backend/src/logger"
)

//证书角色
const (
	ORG_CA    = 1
	ADMIN     = 2
	CLIENT    = 3
	CONSENSUS = 4
	COMMON    = 5
)

//证书类型
const (
	ORG  = 0
	NODE = 1
	USER = 2
)

//证书用途
const (
	SIGN = 0
	TLS  = 1
)

var (
	log = loggers.GetLogger(loggers.ModuleDb)
)

func CreateCert(cert *common.Cert) error {
	if err := connection.DB.Create(&cert).Error; err != nil {
		log.Error("Save cert Failed: " + err.Error())
		return err
	}
	return nil
}

func GetOrgCaCert(orgId string) (*common.Cert, error) {
	var cert common.Cert
	if err := connection.DB.Where("org_id = ? AND cert_type = ?", orgId, ORG_CA).Find(&cert).Error; err != nil {
		log.Error("QueryOrgCaCert Failed: " + err.Error())
		return nil, err
	}
	return &cert, nil
}

func GetUserTlsCert(userName string) (*common.Cert, error) {
	var cert common.Cert
	if err := connection.DB.Where("cert_user_name = ?  AND cert_use = ?", userName, TLS).
		Find(&cert).Error; err != nil {
		log.Error("QueryOrgCaCert Failed: " + err.Error())
		return nil, err
	}
	return &cert, nil
}

func GetUserCertByOrgId(orgId string, certType int) (*common.Cert, error) {
	var cert common.Cert
	if err := connection.DB.Where("org_id = ? AND cert_type = ? AND cert_use = ?", orgId, certType, TLS).
		Limit(1).Find(&cert).Error; err != nil {
		log.Error("GetAdminUserCertByOrgId Failed: " + err.Error())
		return nil, err
	}
	return &cert, nil
}

func GetOrgCaCertCount(orgId string) (int64, error) {
	var count int64
	var cert common.Cert
	if err := connection.DB.Where("org_id = ? AND cert_type = ?", orgId, ORG_CA).Model(&cert).
		Count(&count).Error; err != nil {
		log.Error("GetOrgCaCertCount Failed: " + err.Error())
		return 0, err
	}
	return count, nil
}

func GetOrgCaCertCountBydOrgIdAndOrgName(orgId, orgName string) (int64, error) {
	var count int64
	var cert common.Cert
	if err := connection.DB.Where("(org_id = ? OR org_name = ?) AND cert_type = ?", orgId, orgName, ORG_CA).
		Model(&cert).Count(&count).Error; err != nil {
		log.Error("GetOrgCaCertCount Failed: " + err.Error())
		return 0, err
	}
	return count, nil
}

func GetNodeCertCount(nodeName string) (int64, error) {
	var count int64
	var cert common.Cert
	if err := connection.DB.Where("Node_name = ? AND (cert_type = ? OR cert_type = ?) ", nodeName, CONSENSUS, COMMON).
		Model(&cert).Count(&count).Error; err != nil {
		log.Error("GetOrgCaCertCount Failed: " + err.Error())
		return 0, err
	}
	return count, nil
}

func GetNodeCert(nodeName string) ([]*common.Cert, error) {
	var certs []*common.Cert
	if err := connection.DB.Where("Node_name = ? AND (cert_type = ? OR cert_type = ?) ", nodeName, CONSENSUS, COMMON).
		Find(&certs).Error; err != nil {
		log.Error("GetOrgCaCertCount Failed: " + err.Error())
		return nil, err
	}
	return certs, nil
}

func GetUserCertCount(userName string) (int64, error) {
	var count int64
	var cert common.Cert
	if err := connection.DB.Where("cert_user_name = ? AND (cert_type = ? OR cert_type = ?) ", userName, ADMIN, CLIENT).
		Model(&cert).Count(&count).Error; err != nil {
		log.Error("GetOrgCaCertCount Failed: " + err.Error())
		return 0, err
	}
	return count, nil
}

func GetCertById(id int64) (*common.Cert, error) {
	var cert common.Cert
	if err := connection.DB.Where("id = ?", id).Find(&cert).Error; err != nil {
		log.Error("QueryOrgCaCert Failed: " + err.Error())
		return nil, err
	}
	return &cert, nil
}

func GetCertList(pageNum int64, pageSize int, certType int, orgName, nodeName, userName string) (
	[]*common.Cert, int64, error) {
	var certs []*common.Cert

	db := connection.DB
	if orgName != "" {
		db = db.Where("org_name = ?", orgName)
	}
	if nodeName != "" {
		db = db.Where("node_name = ?", nodeName)
	}
	if userName != "" {
		db = db.Where("cert_user_name = ?", userName)
	}
	if certType == ORG {
		db = db.Where("cert_type = ?", ORG_CA)
	}
	if certType == NODE {
		db = db.Where("cert_type = ? OR cert_type = ?", CONSENSUS, COMMON)
	}
	if certType == USER {
		db = db.Where("cert_type = ? OR cert_type = ?", ADMIN, CLIENT)
	}

	offset := pageNum * int64(pageSize/2)
	if err := db.Order("id DESC").Offset(offset).Limit(pageSize / 2).Find(&certs).Error; err != nil {
		log.Error("GetCertList Failed: " + err.Error())
		return nil, 0, err
	}
	var count int64
	if err := db.Model(&certs).Count(&count).Error; err != nil {
		log.Error("GetCertListCount Failed: " + err.Error())
		return nil, 0, err
	}

	return certs, count, nil
}

func GetUserCertList(orgId string) ([]*common.Cert, int64, error) {
	var count int64
	var certs []*common.Cert

	db := connection.DB
	if orgId != "" {
		db = db.Where("org_id = ?", orgId)
	}
	db = db.Where("cert_type = ? OR cert_type = ?", ADMIN, CLIENT)

	if err := db.Find(&certs).Error; err != nil {
		log.Error("QueryOrgCaCert Failed: " + err.Error())
		return nil, 0, err
	}

	if err := db.Model(&certs).Count(&count).Error; err != nil {
		log.Error("GetOrgCaCertCount Failed: " + err.Error())
		return nil, 0, err
	}
	return certs, count, nil
}

func GetSignUserCertList(orgId string) ([]*common.Cert, int64, error) {
	var count int64
	var certs []*common.Cert

	db := connection.DB
	if orgId != "" {
		db = db.Where("org_id = ?", orgId)
	}
	db = db.Where("cert_type = ? OR cert_type = ?", ADMIN, CLIENT)
	db = db.Where("cert_use = ?", SIGN)
	if err := db.Find(&certs).Error; err != nil {
		log.Error("QueryOrgCaCert Failed: " + err.Error())
		return nil, 0, err
	}

	if err := db.Model(&certs).Count(&count).Error; err != nil {
		log.Error("GetOrgCaCertCount Failed: " + err.Error())
		return nil, 0, err
	}
	return certs, count, nil
}
