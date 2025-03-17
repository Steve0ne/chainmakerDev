/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package common

import "chainmaker.org/chainmaker/pb-go/v2/common"

type ContractStatus int

const (
	ContractInitStored ContractStatus = iota
	ContractUpgradeStored
	ContractInitFailure
	ContractInitOK
	ContractUpgradeFailure
	ContractUpgradeOK
	ContractFreezeFailure
	ContractFreezeOK
	ContractUnfreezeFailure
	ContractUnfreezeOK
	ContractRevokeFailure
	ContractRevokeOK
	ContractMultiSign
)

func RuntimeTypeString(runtimeType int) string {
	return common.RuntimeType(runtimeType).String()
}

func (c *Contract) GetRuntimeTypeString() string {
	return RuntimeTypeString(c.RuntimeType)
}

func (c *Contract) CanUpgrade() bool {
	contractStatus := ContractStatus(c.ContractStatus)
	return contractStatus == ContractInitOK || contractStatus == ContractUpgradeOK ||
		contractStatus == ContractFreezeFailure || contractStatus == ContractUnfreezeOK ||
		contractStatus == ContractRevokeFailure || contractStatus == ContractUpgradeFailure
}

func (c *Contract) CanInstall() bool {
	contractStatus := ContractStatus(c.ContractStatus)
	return contractStatus == ContractInitStored || contractStatus == ContractInitOK ||
		contractStatus == ContractInitFailure || contractStatus == ContractUpgradeOK ||
		contractStatus == ContractUpgradeStored || contractStatus == ContractUpgradeFailure ||
		contractStatus == ContractRevokeFailure
}

func (c *Contract) CanUpgradeDeploy() bool {
	contractStatus := ContractStatus(c.ContractStatus)
	return contractStatus == ContractUpgradeStored || contractStatus == ContractUpgradeFailure
}

func (c *Contract) CanInstallDeploy() bool {
	contractStatus := ContractStatus(c.ContractStatus)
	return contractStatus == ContractInitStored || contractStatus == ContractInitFailure
}

func (c *Contract) CanFreeze() bool {
	contractStatus := ContractStatus(c.ContractStatus)
	return contractStatus == ContractInitOK || contractStatus == ContractUpgradeStored ||
		contractStatus == ContractUpgradeFailure || contractStatus == ContractUpgradeOK ||
		contractStatus == ContractFreezeFailure || contractStatus == ContractUnfreezeOK
}

func (c *Contract) CanUnfreeze() bool {
	contractStatus := ContractStatus(c.ContractStatus)
	return contractStatus == ContractFreezeOK
}

func (c *Contract) CanRevoke() bool {
	contractStatus := ContractStatus(c.ContractStatus)
	// 初始化成功过，并且未被注销
	return contractStatus >= ContractInitOK && contractStatus != ContractRevokeOK
}
