/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/config"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"

	"management_backend/src/db/chain"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/utils"
)

/*
	区块解析结构
*/

const CHAIN_CONFIG = "CHAIN_CONFIG"

func ParseBlockToDB(blockInfo *common.BlockInfo) error {
	var (
		modBlock dbcommon.Block
		err      error
	)
	modBlock.BlockHeight = blockInfo.Block.Header.BlockHeight
	modBlock.BlockHash = hex.EncodeToString(blockInfo.Block.Header.BlockHash)
	modBlock.ChainId = blockInfo.Block.Header.ChainId
	modBlock.PreBlockHash = hex.EncodeToString(blockInfo.Block.Header.PreBlockHash)
	modBlock.ConsensusArgs = utils.Base64Encode(blockInfo.Block.Header.ConsensusArgs)
	modBlock.DagHash = hex.EncodeToString(blockInfo.Block.Header.DagHash)
	modBlock.OrgId, modBlock.ProposerId, err = parseMember(blockInfo.Block.Header.Proposer, modBlock.ChainId)
	if err != nil {
		err = fmt.Errorf("parse block proposer member failed: %s", err.Error())
	}
	//modBlock.ProposerType = blockInfo.Block.Header.Proposer.MemberType.String()
	modBlock.RwSetHash = hex.EncodeToString(blockInfo.Block.Header.RwSetRoot)
	modBlock.Timestamp = blockInfo.Block.Header.BlockTimestamp
	modBlock.TxCount = int(blockInfo.Block.Header.TxCount)
	modBlock.TxRootHash = hex.EncodeToString(blockInfo.Block.Header.TxRoot)

	//modBlock.OrgId = blockInfo.Block.Header.Proposer.OrgId

	transactions, contracts, err := parseTxsAndContracts(blockInfo)
	if err != nil {
		log.Error("parseTxsAndContracts err : " + err.Error())
		return err
	}

	err = chain.InsertBlockAndTx(&modBlock, transactions, contracts)
	if err != nil {
		log.Error("Insert Block And Tx Failed : " + err.Error())
		return err
	}

	return nil
}

func parseTxsAndContracts(blockInfo *common.BlockInfo) ([]*dbcommon.Transaction, []*dbcommon.Contract, error) {
	var transactions = make([]*dbcommon.Transaction, 0)
	var contracts = make([]*dbcommon.Contract, 0)
	var err error
	chainId := blockInfo.Block.Header.ChainId
	for _, tx := range blockInfo.Block.Txs {
		transaction := &dbcommon.Transaction{}
		transaction.BlockHeight = blockInfo.Block.Header.BlockHeight
		transaction.BlockHash = hex.EncodeToString(blockInfo.Block.Header.BlockHash)
		transaction.ChainId = tx.Payload.ChainId

		payload := tx.Payload
		transaction.TxId = payload.TxId
		transaction.TxType = payload.TxType.String()
		transaction.TxStatusCode = tx.Result.Code.String()
		transaction.Timestamp = payload.Timestamp
		//transaction.ExpirationTime = payload.ExpirationTime
		//transaction.ContractName = payload.ContractName
		transaction.ContractMethod = payload.Method
		transaction.Sequence = payload.Sequence
		//transaction.Limit = payload.Limit

		if tx.Sender != nil {
			transaction.OrgId, transaction.Sender, err = parseMember(tx.Sender.Signer, chainId)
			if err != nil {
				log.Error("parse Tx sender member info failed : " + err.Error())
			}
		}

		transaction.Endorsers = parseTxEndorsers(tx.Endorsers)

		transaction.TXResult = parseTxResult(tx.Result)
		if err != nil {
			log.Error("parseTxResult Failed: " + err.Error())
			return nil, nil, err
		}

		if payload.TxType == common.TxType_INVOKE_CONTRACT {
			transaction.ContractName = payload.ContractName
			transaction.ContractMethod = payload.Method
			ParametersBytes, err := json.Marshal(payload.Parameters)
			if err != nil {
				log.Error("Contract Parameters Marshal Failed: " + err.Error())
				return nil, nil, err
			}
			transaction.ContractParameters = string(ParametersBytes)
			for _, parameter := range payload.Parameters {
				if parameter.Key == syscontract.InitContract_CONTRACT_VERSION.String() {
					transaction.ContractVersion = string(parameter.Value)
				}
				if parameter.Key == syscontract.InitContract_CONTRACT_RUNTIME_TYPE.String() {
					transaction.ContractRuntimeType = string(parameter.Value)
				}
			}

			if tx.Result.ContractResult.Code == 0 {
				if payload.ContractName == syscontract.SystemContract_CONTRACT_MANAGE.String() {
					//合约操作 需要更新合约的状态（升级成功，冻结成功，初始化成功等）
					var contractName string
					var runtimeType int
					var contract = &dbcommon.Contract{}
					for _, parameter := range payload.Parameters {
						if parameter.Key == syscontract.InitContract_CONTRACT_NAME.String() {
							contractName = string(parameter.Value)
						}
						if parameter.Key == syscontract.InitContract_CONTRACT_RUNTIME_TYPE.String() {
							runtimeType = int(common.RuntimeType_value[string(parameter.Value)])
						}
					}

					// 这里暂时更改了交易的合约名，为的是 管理用户合约的交易也能通过查询用户合约查询到
					transaction.ContractName = contractName
					contract.Name = contractName
					contract.OrgId = tx.Sender.Signer.OrgId
					contract.RuntimeType = runtimeType
					contract.ChainId = payload.ChainId
					contract.Version = transaction.ContractVersion
					contract.MultiSignStatus = dbcommon.NO_VOTING
					contract.ContractStatus = getContractStatus(payload.Method)
					contract.Timestamp = payload.Timestamp
					contract.TxId = transaction.TxId

					contracts = append(contracts, contract)
				}
			}
		}

		transactions = append(transactions, transaction)

		// 合约为链管理类合约，则更新链配置
		if transaction.ContractName == syscontract.SystemContract_CHAIN_CONFIG.String() {
			updateChainConfig(chainId, transaction.BlockHeight, blockInfo.Block.Header.BlockTimestamp)
		}
	}
	return transactions, contracts, nil
}

func parseTxResult(result *common.Result) dbcommon.TXResult {
	var txResult dbcommon.TXResult

	txResult.ResultCode = result.Code.String()
	txResult.ResultMessage = result.Message
	txResult.RwSetHash = hex.EncodeToString(result.RwSetHash)

	txResult.ContractResult = result.ContractResult.Result
	txResult.ContractResultCode = result.ContractResult.Code
	txResult.ContractResultMessage = result.ContractResult.Message
	txResult.Gas = result.ContractResult.GasUsed

	return txResult
}

func parseTxEndorsers(endorsers []*common.EndorsementEntry) string {
	return ""
}

func getContractStatus(mgmtMethod string) int {
	switch mgmtMethod {
	case syscontract.ContractManageFunction_INIT_CONTRACT.String():
		return int(dbcommon.ContractInitOK)
	case syscontract.ContractManageFunction_UPGRADE_CONTRACT.String():
		return int(dbcommon.ContractUpgradeOK)
	case syscontract.ContractManageFunction_FREEZE_CONTRACT.String():
		return int(dbcommon.ContractFreezeOK)
	case syscontract.ContractManageFunction_UNFREEZE_CONTRACT.String():
		return int(dbcommon.ContractUnfreezeOK)
	case syscontract.ContractManageFunction_REVOKE_CONTRACT.String():
		return int(dbcommon.ContractRevokeOK)
	}
	return -1
}

func parseMember(sender *accesscontrol.Member, chainId string) (orgId string, memberId string, err error) {
	var (
		x509Cert  *x509.Certificate
		certBytes []byte
		resp      *common.CertInfos
	)

	if sender != nil {
		orgId = sender.OrgId

		switch sender.MemberType {
		case accesscontrol.MemberType_CERT:
			certBytes = sender.MemberInfo
			x509Cert, err = utils.ParseCertificate(certBytes)
			if err == nil {
				memberId = x509Cert.Subject.CommonName
			}

		case accesscontrol.MemberType_CERT_HASH:
			certBytes = sender.MemberInfo
			sdkClientPool := GetSdkClientPool()
			sdkClient := sdkClientPool.SdkClients[chainId]
			if sdkClient == nil {
				err = errors.New("ClientIsNil")
				return
			}
			resp, err = sdkClient.ChainClient.QueryCert([]string{hex.EncodeToString(sender.MemberInfo)})
			if err != nil {
				return
			}
			certBytes = resp.CertInfos[0].Cert

			x509Cert, err = utils.ParseCertificate(certBytes)
			if err == nil {
				memberId = x509Cert.Subject.CommonName
			}
		}

	} else {
		err = errors.New("SenderIsNil")
	}
	return
}

func updateChainConfig(chainId string, blockHeight uint64, blockTime int64) {
	var (
		err               error
		config            *config.ChainConfig
		chainConfigRecord *dbcommon.ChainConfig
	)
	sdkClientPool := GetSdkClientPool()
	sdkClient := sdkClientPool.SdkClients[chainId]
	if sdkClient == nil {
		err = errors.New("ClientIsNil")
		return
	}
	config, err = sdkClient.ChainClient.GetChainConfig()
	if err != nil {
		return
	}

	configString, err := json.Marshal(config)

	chainConfigRecord = &dbcommon.ChainConfig{
		ChainId:     chainId,
		BlockHeight: blockHeight,
		BlockTime:   blockTime,
		Config:      string(configString),
	}

	err = chain.CreateChainConfigRecord(chainConfigRecord)
	if err != nil {
		return
	}

	return
}
