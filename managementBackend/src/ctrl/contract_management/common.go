/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_management

import (
	bcx509 "chainmaker.org/chainmaker/common/v2/crypto/x509"
	"chainmaker.org/chainmaker/common/v2/evmutils"
	pbcommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"management_backend/src/ctrl/ca"
	"math/rand"
	"strconv"
	"strings"

	"management_backend/src/ctrl/common"
	"management_backend/src/db"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
	dbcontract "management_backend/src/db/contract"
	"management_backend/src/db/policy"
	"management_backend/src/db/relation"
	loggers "management_backend/src/logger"
	"management_backend/src/sync"
)

const (
	NODE_ADDR_UPDATE = iota
	TRUST_ROOT_UPDATE
	CONSENSUS_EXT_DELETE
	BLOCK_UPDATE
	INIT_CONTRACT
	UPGRADE_CONTRACT
	FREEZE_CONTRACT
	UNFREEZE_CONTRACT
	REVOKE_CONTRACT
	PERMISSION_UPDATE
)

const (
	UPDATE_CONFIG = iota
	UPDATE_AUTH
	UPDATE_OTHER
)

const (
	FUNCTION = iota
	CONSTRUCTOR
)

const NULL = "null"

const EVM = 5

var log = loggers.GetLogger(loggers.ModuleWeb)

func SaveVote(chainId, reason, paramJson string, voteType, configType int) (string, error) {
	sdkClientPool := sync.GetSdkClientPool()
	if sdkClientPool == nil {
		newError := common.CreateError(common.ErrorChainNotSub)
		return "", newError
	}
	var orgId string
	if sdkClientPool.SdkClients[chainId] == nil {
		newError := common.CreateError(common.ErrorChainNotSub)
		return "", newError
	}
	orgId = sdkClientPool.SdkClients[chainId].SdkConfig.OrgId
	orgName, err := chain_participant.GetOrgNameByOrgId(orgId)
	if err != nil {
		newError := common.CreateError(common.ErrorGetOrgName)
		return "", newError
	}

	VoteDetailMessage, err := getMessage(chainId, orgName, paramJson, voteType)
	if err != nil {
		newError := common.CreateError(common.ErrorGetMessage)
		return "", newError
	}

	chainPolicy, err := policy.GetChainPolicy(chainId, voteType)
	if err != nil {
		newError := common.CreateError(common.ErrorGetChainPolicy)
		return "", newError
	}

	orgList, err := relation.GetChainOrgList(chainId)
	if err != nil {
		newError := common.CreateError(common.ErrorGetOrg)
		return "", newError
	}
	multiId := strconv.Itoa(rand.Int())
	for _, org := range orgList {
		vote := &dbcommon.VoteManagement{
			MultiId:      multiId,
			ChainId:      chainId,
			StartOrgId:   orgId,
			StartOrgName: orgName,
			VoteOrgId:    org.OrgId,
			VoteOrgName:  org.OrgName,
			VoteType:     voteType,
			PolicyType:   chainPolicy.PolicyType,
			VoteResult:   0,
			VoteStatus:   0,
			Reason:       reason,
			VoteDetail:   VoteDetailMessage,
			Params:       paramJson,
			ConfigStatus: configType,
		}
		err = db.CreateVote(vote)
		if err != nil {
			newError := common.CreateError(common.ErrorCreateVote)
			return "", newError
		}
	}
	return orgId, nil
}

func UpdateMultiSignStatus(contract *dbcommon.Contract) error {
	err := dbcontract.UpdateContractMultiSignStatus(contract)
	if err != nil {
		log.Error("update vote status failed:", err.Error())
		return err
	}
	return nil
}

func GetEvmMethodsByAbi(abiKey string) (string, int, error) {
	id, userId, hash, err := ca.ResolveUploadKey(abiKey)
	if err != nil {
		return "", -1, err
	}
	upload, err := db.GetUploadByIdAndUserIdAndHash(id, userId, hash)
	if err != nil {
		return "", -1, err
	}

	myAbi, err := abi.JSON(strings.NewReader(string(upload.Content)))
	if err != nil {
		return "", -1, err
	}
	var methods = make([]*Method, 0)

	// 0：正常方法 1：构造函数
	functionType := FUNCTION
	if len(myAbi.Constructor.Inputs) > 0 {
		functionType = CONSTRUCTOR
	}
	for methodName, methodVale := range myAbi.Methods {
		method := &Method{}
		method.MethodName = methodName

		var methodKeyStr string
		inputs := methodVale.Inputs
		for _, input := range inputs {
			methodKeyStr = methodKeyStr + input.Name + ","
		}

		methodKeyStr = strings.TrimRight(methodKeyStr, ",")
		method.MethodKey = methodKeyStr
		methods = append(methods, method)
	}

	methodJson, err := json.Marshal(methods)
	if err != nil {
		return "", -1, err
	}

	methodStr := string(methodJson)
	if methodStr == NULL {
		methodStr = ""
	}
	return methodStr, functionType, nil
}

func MakeAddrAndSkiFromCrtBytes(crtBytes []byte) (string, string, string, error) {
	blockCrt, _ := pem.Decode(crtBytes)
	crt, err := bcx509.ParseCertificate(blockCrt.Bytes)
	if err != nil {
		return "", "", "", err
	}

	ski := hex.EncodeToString(crt.SubjectKeyId)
	addrInt, err := evmutils.MakeAddressFromHex(ski)
	if err != nil {
		return "", "", "", err
	}

	fmt.Sprintf("0x%s", addrInt.AsStringKey())

	return addrInt.String(), fmt.Sprintf("0x%x", addrInt.AsStringKey()), ski, nil
}

func GetConstructorKeyValuePair(crtBytes []byte, abiKey string) ([]*pbcommon.KeyValuePair, error) {
	_, _, client1AddrSki, err :=
		MakeAddrAndSkiFromCrtBytes(crtBytes)
	if err != nil {
		return nil, err
	}
	addrInt, err := evmutils.MakeAddressFromHex(client1AddrSki)
	if err != nil {
		return nil, err
	}
	addr := evmutils.BigToAddress(addrInt)

	id, userId, hash, err := ca.ResolveUploadKey(abiKey)
	if err != nil {
		return nil, err
	}
	upload, err := db.GetUploadByIdAndUserIdAndHash(id, userId, hash)
	if err != nil {
		return nil, err
	}

	myAbi, err := abi.JSON(strings.NewReader(string(upload.Content)))
	if err != nil {
		return nil, err
	}

	dataByte, err := myAbi.Pack("", addr)
	if err != nil {
		return nil, err
	}

	data := hex.EncodeToString(dataByte)
	pairs := []*pbcommon.KeyValuePair{
		{
			Key:   "data",
			Value: []byte(data),
		},
	}

	return pairs, nil
}
