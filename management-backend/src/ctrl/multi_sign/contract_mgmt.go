/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package multi_sign

import (
	"chainmaker.org/chainmaker/common/v2/evmutils"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"

	"management_backend/src/ctrl/ca"
	"management_backend/src/ctrl/common"
	"management_backend/src/ctrl/contract_management"
	"management_backend/src/db"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/contract"
	"management_backend/src/sync"
	"management_backend/src/utils"
)

const (
	VOTING = iota
	NO_VOTING
)

const TxHandleTimeout = 10

const NULL = "null"

type contractOpType int

const (
	contractOpTypeFreeze contractOpType = 2 + iota
	contractOpTypeUnfreeze
	contractOpTypeRevoke
)

func ContractInstallModify(parameters string, orgList []string, roleType, installType int) error {
	var contractInstallBody contract_management.InstallContractParams
	err := json.Unmarshal([]byte(parameters), &contractInstallBody)
	if err != nil {
		log.Errorf("Unmarshal parameters to contractInstallBody err:, %s", err)
		return err
	}
	chainId := contractInstallBody.ChainId
	contractName := contractInstallBody.ContractName
	dbContract, err := contract.GetContractByName(chainId, contractName)
	if err != nil {
		//newError := common.CreateError(common.ErrorInstallContract, "没有可用的合约")
		return err
	}
	id, userId, hash, err := ca.ResolveUploadKey(contractInstallBody.CompileSaveKey)
	if err != nil {
		//newError := entity.NewError(entity.ErrorContractInstall, "该合约源文件key错误")
		return err
	}
	upload, err := db.GetUploadByIdAndUserIdAndHash(id, userId, hash)
	if err != nil {
		//newError := entity.NewError(entity.ErrorContractInstall, "该合约源文件错误")
		return err
	}
	var newContractStatus dbcommon.ContractStatus
	// 同步发布合约负责将该合约发布
	dbContract.Version = contractInstallBody.ContractVersion
	err = installContract(dbContract, upload, convertToPbKeyValues(&contractInstallBody), orgList, roleType, installType)
	if err != nil {
		if installType == contract_management.INIT_CONTRACT {
			newContractStatus = dbcommon.ContractInitFailure
		} else if installType == contract_management.UPGRADE_CONTRACT {
			newContractStatus = dbcommon.ContractUpgradeFailure
		}
		_ = contract.UpdateContractStatus(dbContract.Id, int(newContractStatus), NO_VOTING)
		return err
	}
	if installType == contract_management.INIT_CONTRACT {
		newContractStatus = dbcommon.ContractInitOK
	} else if installType == contract_management.UPGRADE_CONTRACT {
		newContractStatus = dbcommon.ContractUpgradeOK
	}
	// 修改当前合约的状态
	err = contract.UpdateContractStatus(dbContract.Id, int(newContractStatus), NO_VOTING)
	if err != nil {
		return err
	}

	var methodStr string
	var functionType int
	if contractInstallBody.RuntimeType == contract_management.EVM {
		methodStr, functionType, err = contract_management.GetEvmMethodsByAbi(contractInstallBody.EvmAbiSaveKey)
		if err != nil {
			return errors.New("get evmmethods from Abi err")
		}
	} else {
		methodJson, err := json.Marshal(contractInstallBody.Methods)
		if err != nil {
			return err
		}

		methodStr = string(methodJson)
		if methodStr == NULL {
			methodStr = ""
		}
	}

	contractInfo := &dbcommon.Contract{
		Name:            contractInstallBody.ContractName,
		Methods:         methodStr,
		SourceSaveKey:   contractInstallBody.CompileSaveKey,
		EvmAbiSaveKey:   contractInstallBody.EvmAbiSaveKey,
		EvmFunctionType: functionType,
	}
	err = contract.UpdateContractMethodByName(contractInfo)
	if err != nil {
		return err
	}
	return nil
}

func convertToPbKeyValues(body *contract_management.InstallContractParams) []*pbcommon.KeyValuePair {
	keyValues := body.Parameters
	if len(keyValues) > 0 {
		pbKvs := make([]*pbcommon.KeyValuePair, 0)
		for _, kv := range keyValues {
			pbKvs = append(pbKvs, &pbcommon.KeyValuePair{
				Key:   kv.Key,
				Value: []byte(kv.Value),
			})
		}
		return pbKvs
	}
	return []*pbcommon.KeyValuePair{}
}

func installContract(contract *dbcommon.Contract, upload *dbcommon.Upload,
	keyValues []*pbcommon.KeyValuePair, orgList []string, roleType, installType int) error {
	sdkClientPool := sync.GetSdkClientPool()
	if sdkClientPool == nil {
		newError := common.CreateError(common.ErrorChainNotSub)
		return newError
	}
	sdkClient := sdkClientPool.SdkClients[contract.ChainId]
	chainClient := sdkClient.ChainClient
	var content string
	content = utils.Base64Encode(upload.Content)
	if contract.RuntimeType == 5 {
		content = string(upload.Content)
	}
	var (
		payload *pbcommon.Payload
		err     error
	)

	if keyValues == nil {
		keyValues = []*pbcommon.KeyValuePair{}
	}
	var contractName string
	contractName = contract.Name
	if contract.RuntimeType == contract_management.EVM {
		contractName = hex.EncodeToString(evmutils.Keccak256([]byte(contractName)))[24:]

		if contract.EvmFunctionType == contract_management.CONSTRUCTOR {
			keyValues, err = contract_management.GetConstructorKeyValuePair(sdkClient.SdkConfig.UserCert, contract.EvmAbiSaveKey)
			if err != nil {
				return err
			}
		}
	}

	if installType == contract_management.INIT_CONTRACT {
		// 新建合约
		if !contract.CanInstall() {
			log.Error("contract cann't install")
			return errors.New("contract cann't install")
		}
		payload, err = chainClient.CreateContractCreatePayload(
			contractName, contract.Version, content, pbcommon.RuntimeType(contract.RuntimeType), keyValues)
		if err != nil {
			return err
		}
	} else if installType == contract_management.UPGRADE_CONTRACT {
		// 升级合约
		if !contract.CanUpgrade() {
			log.Error("contract cann't upgrade")
			return errors.New("contract cann't upgrade")
		}
		payload, err = chainClient.CreateContractUpgradePayload(
			contract.Name, contract.Version, content, pbcommon.RuntimeType(contract.RuntimeType), keyValues)
		if err != nil {
			return err
		}
	}

	endorsements, err := GetEndorsements(payload, orgList, roleType)
	if err != nil {
		return err
	}

	// 发送创建合约请求
	resp, err := chainClient.SendContractManageRequest(payload, endorsements, TxHandleTimeout, true)
	if err != nil {
		//return entity.CreateError(entity.ErrorHandleFailure, err.Error())
		return err
	}
	// 判断结果
	if resp.Code != pbcommon.TxStatusCode_SUCCESS {
		// 失败
		//return entity.NewError(entity.ErrorHandleFailure, "install contract failure")
		return errors.New("install contract failure")
	}

	if resp.ContractResult == nil {
		return fmt.Errorf("contract result is nil")
	}

	if resp.ContractResult != nil && resp.ContractResult.Code != 0 {
		//return entity.NewError(entity.ErrorHandleFailure, resp.ContractResult.Message)
		return errors.New("install contract failure")
	}

	return nil
}

func ContractFreezeModify(parameters string, orgList []string, roleType int) error {
	var contractFreezeBody contract_management.FreezeContractParams
	err := json.Unmarshal([]byte(parameters), &contractFreezeBody)
	if err != nil {
		log.Errorf("Unmarshal parameters to contractInstallBody err:, %s", err)
		return err
	}

	chainId := contractFreezeBody.ChainId
	contractName := contractFreezeBody.ContractName
	dbContract, err := contract.GetContractByName(chainId, contractName)
	if err != nil {
		//newError := entity.NewError(entity.ErrorInstallContract, "没有可用的合约")
		return err
	}
	// 检查之前合约状态，合约必须处理初始化成功、升级成功状态、解冻成功状态才可以进行冻结
	if !dbContract.CanFreeze() {
		// 不可以进行冻结操作
		//newError := entity.NewError(entity.ErrorContractInstall, "该合约不能冻结")
		log.Error("contract cann't freeze")
		return errors.New("contract cann't freeze")
	}

	if err := mgmtContract(chainId, contractName, contractOpTypeFreeze, orgList, roleType); err != nil {
		newErr := contract.UpdateContractStatus(dbContract.Id, int(dbcommon.ContractFreezeFailure), NO_VOTING)
		// 状态更新为冻结失败
		if newErr != nil {
			return newErr
		}
		return err
	}

	return contract.UpdateContractStatus(dbContract.Id, int(dbcommon.ContractFreezeOK), NO_VOTING)
}

func ContractUnfreezeModify(parameters string, orgList []string, roleType int) error {
	var contractUnFreezeBody contract_management.FreezeContractParams
	err := json.Unmarshal([]byte(parameters), &contractUnFreezeBody)
	if err != nil {
		log.Errorf("Unmarshal parameters to contractInstallBody err:, %s", err)
		return err
	}

	chainId := contractUnFreezeBody.ChainId
	contractName := contractUnFreezeBody.ContractName
	dbContract, err := contract.GetContractByName(chainId, contractName)
	if err != nil {
		//newError := entity.NewError(entity.ErrorInstallContract, "没有可用的合约")
		return err
	}
	// 检查之前合约状态，合约必须处理初始化成功、升级成功状态、解冻成功状态才可以进行冻结
	if !dbContract.CanUnfreeze() {
		// 不可以进行冻结操作
		log.Error("contract cann't unfreeze")
		return errors.New("contract cann't unfreeze")
	}

	if err := mgmtContract(chainId, contractName, contractOpTypeUnfreeze, orgList, roleType); err != nil {
		newErr := contract.UpdateContractStatus(dbContract.Id, int(dbcommon.ContractUnfreezeFailure), NO_VOTING)
		// 状态更新为冻结失败
		if newErr != nil {
			return newErr
		}
		return err
	}

	return contract.UpdateContractStatus(dbContract.Id, int(dbcommon.ContractUnfreezeOK), NO_VOTING)
}

func ContractRevokeModify(parameters string, orgList []string, roleType int) error {
	var contractRevokeBody contract_management.FreezeContractParams
	err := json.Unmarshal([]byte(parameters), &contractRevokeBody)
	if err != nil {
		log.Errorf("Unmarshal parameters to contractInstallBody err:, %s", err)
		return err
	}

	chainId := contractRevokeBody.ChainId
	contractName := contractRevokeBody.ContractName
	dbContract, err := contract.GetContractByName(chainId, contractName)
	if err != nil {
		return err
	}
	// 检查之前合约状态，合约必须处理初始化成功、升级成功状态、解冻成功状态才可以进行冻结
	if !dbContract.CanRevoke() {
		// 不可以进行注销操作
		log.Error("contract cann't revoke")
		return errors.New("contract cann't revoke")
	}

	if err := mgmtContract(chainId, contractName, contractOpTypeRevoke, orgList, roleType); err != nil {
		// 状态更新为冻结失败
		newErr := contract.UpdateContractStatus(dbContract.Id, int(dbcommon.ContractRevokeFailure), NO_VOTING)
		if newErr != nil {
			return newErr
		}
		return err
	}

	return contract.UpdateContractStatus(dbContract.Id, int(dbcommon.ContractRevokeOK), NO_VOTING)
}

func mgmtContract(chainId, contractName string, opType contractOpType, orgList []string, roleType int) error {
	sdkClientPool := sync.GetSdkClientPool()
	if sdkClientPool == nil {
		newError := common.CreateError(common.ErrorChainNotSub)
		return newError
	}
	sdkClient := sdkClientPool.SdkClients[chainId]
	chainClient := sdkClient.ChainClient
	var (
		payload *pbcommon.Payload
		err     error
	)
	if opType == contractOpTypeFreeze {
		payload, err = chainClient.CreateContractFreezePayload(contractName)
	} else if opType == contractOpTypeUnfreeze {
		payload, err = chainClient.CreateContractUnfreezePayload(contractName)
	} else if opType == contractOpTypeRevoke {
		payload, err = chainClient.CreateContractRevokePayload(contractName)
	}
	if err != nil {
		//return common.CreateError(entity.ErrorHandleFailure, "sdk client is nil")
		return err
	}

	endorsements, err := GetEndorsements(payload, orgList, roleType)
	if err != nil {
		//return entity.NewError(entity.ErrorHandleFailure, err.Error())
		return err
	}

	// 发送创建合约请求
	resp, err := chainClient.SendContractManageRequest(payload, endorsements, TxHandleTimeout, true)
	if err != nil {
		//return entity.NewError(entity.ErrorHandleFailure, err.Error())
		return err
	}
	// 判断结果
	if resp.Code != pbcommon.TxStatusCode_SUCCESS {
		// 失败
		//return entity.NewError(entity.ErrorHandleFailure, "mgmt contract failure")
		return err
	}
	return nil
}
