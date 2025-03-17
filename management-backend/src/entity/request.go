/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package entity

const (
	Project = "chainmaker"
	CMB     = "cmb"
	Token   = "token"

	// 文件管理

	UploadFile = "UploadFile"

	// 证书管理

	GenerateCert = "GenerateCert"
	GetCert      = "GetCert"
	GetCertList  = "GetCertList"
	ImportCert   = "ImportCert"
	DownloadCert = "DownloadCert"

	// 用户管理

	GetCaptcha     = "GetCaptcha"
	Login          = "Login"
	GetUserList    = "GetUserList"
	AddUser        = "AddUser"
	ModifyPassword = "ModifyPassword"
	Logout         = "Logout"
	EnableUser     = "EnableUser"
	DisableUser    = "DisableUser"
	ResetPassword  = "ResetPassword"

	// 链管理

	AddChain            = "AddChain"
	DeleteChain         = "DeleteChain"
	GetConsensusList    = "GetConsensusList"
	GetCertUserList     = "GetCertUserList"
	GetCertOrgList      = "GetCertOrgList"
	GetCertNodeList     = "GetCertNodeList"
	GetChainList        = "GetChainList"
	SubscribeChain      = "SubscribeChain"
	DownloadChainConfig = "DownloadChainConfig"

	// Explorer api

	GetTxList         = "GetTxList"
	GetTxDetail       = "GetTxDetail"
	GetBlockList      = "GetBlockList"
	GetBlockDetail    = "GetBlockDetail"
	GetContractList   = "GetContractList"
	GetContractDetail = "GetContractDetail"
	HomePageSearch    = "HomePageSearch"

	GetOrgList = "GetOrgList"

	GetNodeList   = "GetNodeList"
	GetNodeDetail = "GetNodeDetail"

	Vote              = "Vote"
	GetVoteManageList = "GetVoteManageList"
	GetVoteDetail     = "GetVoteDetail"

	// 合约管理

	InstallContract       = "InstallContract"
	FreezeContract        = "FreezeContract"
	UnFreezeContract      = "UnFreezeContract"
	RevokeContract        = "RevokeContract"
	UpgradeContract       = "UpgradeContract"
	GetRuntimeTypeList    = "GetRuntimeTypeList"
	GetContractManageList = "GetContractManageList"
	ContractDetail        = "ContractDetail"
	ModifyContract        = "ModifyContract"
	GetInvokeContractList = "GetInvokeContractList"
	InvokeContract        = "InvokeContract"
	ReInvokeContract      = "ReInvokeContract"
	GetInvokeRecordList   = "GetInvokeRecordList"
	GetInvokeRecordDetail = "GetInvokeRecordDetail"

	// 区块链概览

	ModifyChainAuth   = "ModifyChainAuth"
	ModifyChainConfig = "ModifyChainConfig"
	GeneralData       = "GeneralData"
	GetChainDetail    = "GetChainDetail"
	GetAuthOrgList    = "GetAuthOrgList"
	GetAuthRoleList   = "GetAuthRoleList"
	GetAuthList       = "GetAuthList"

	// 日志收集
	DownloadLogFile   = "DownloadLogFile"
	GetLogList        = "GetLogList"
	ReportLogFile     = "ReportLogFile"
	AutoReportLogFile = "AutoReportLogFile"
	PullErrorLog      = "PullErrorLog"
)
