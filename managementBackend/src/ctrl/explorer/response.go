/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package explorer

import (
	"encoding/json"

	"management_backend/src/db/chain"
	"management_backend/src/db/chain_participant"
	dbcommon "management_backend/src/db/common"
)

type BlockView struct {
	Id           int64
	ChainId      string
	BlockHash    string
	PreBlockHash string
	BlockHeight  uint64
	Timestamp    int64
	DagHash      string
	RwSetHash    string
	TxRootHash   string
	TxNum        int
	OrgId        string
	OrgName      string
	NodeName     string
}

// NewBlockView 创建BlockView对象
func NewBlockView(block *dbcommon.Block) *BlockView {

	orgName, _ := chain_participant.GetOrgNameByOrgId(block.OrgId)

	blockView := BlockView{
		Id:           block.Id,
		ChainId:      block.ChainId,
		BlockHash:    block.BlockHash,
		PreBlockHash: block.PreBlockHash,
		BlockHeight:  block.BlockHeight,
		Timestamp:    block.Timestamp,
		DagHash:      block.DagHash,
		RwSetHash:    block.RwSetHash,
		TxRootHash:   block.TxRootHash,
		TxNum:        block.TxCount,
		OrgId:        block.OrgId,
		OrgName:      orgName,
		NodeName:     block.ProposerId,
	}
	return &blockView
}

type ContractView struct {
	Id           int64
	ContractName string
	Creator      string
	TxId         string
	CreateTime   int64
}

// NewContractView 创建ContractView
func NewContractView(contract *dbcommon.Contract) *ContractView {
	return &ContractView{
		Id:           contract.Id,
		ContractName: contract.Name,
		Creator:      contract.ContractOperator,
		TxId:         contract.TxId,
		CreateTime:   contract.Timestamp,
	}
}

type TransactionView struct {
	Id                 int64
	ChainId            string
	BlockHeight        uint64
	BlockHash          string
	TxId               string
	OrgId              string
	OrgName            string
	UserName           string
	Gas                uint64
	TxType             string
	Timestamp          int64
	TxStatusCode       string
	ContractName       string
	ContractMethod     string
	ContractParameters []Parameter
	ContractResult     []byte
	OrgList            []string
}

func NewTransactionView(transaction *dbcommon.Transaction) *TransactionView {
	orgName, _ := chain_participant.GetOrgNameByOrgId(transaction.OrgId)
	return &TransactionView{
		Id:                 transaction.Id,
		ChainId:            transaction.ChainId,
		TxId:               transaction.TxId,
		OrgId:              transaction.OrgId,
		OrgName:            orgName,
		UserName:           transaction.Sender,
		Gas:                transaction.Gas,
		BlockHeight:        transaction.BlockHeight,
		BlockHash:          GetBlockHashByBlockHeight(transaction.ChainId, transaction.BlockHeight),
		TxType:             transaction.TxType,
		Timestamp:          transaction.Timestamp,
		TxStatusCode:       transaction.TxStatusCode,
		ContractName:       transaction.ContractName,
		ContractMethod:     transaction.ContractMethod,
		ContractParameters: FormatContractParameters(transaction.ContractParameters),
		ContractResult:     transaction.ContractResult,
		OrgList:            chain_participant.GetConsensusNodeNameList(transaction.ChainId),
	}
}

func GetBlockHashByBlockHeight(chainId string, blockHeight uint64) string {
	block, err := chain.GetBlockByBlockHeight(chainId, blockHeight)
	if err != nil {
		return ""
	}
	return block.BlockHash
}

type Parameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func FormatContractParameters(params string) (formatParams []Parameter) {
	err := json.Unmarshal([]byte(params), &formatParams)
	if err != nil {
		formatParams = []Parameter{}
		return
	}
	return
}

type HomePageSearchView struct {
	Type int
	Id   int64
}

func NewHomePageSearchView(searchType int, id int64) *HomePageSearchView {
	homePageSearchView := HomePageSearchView{
		Type: searchType,
		Id:   id,
	}
	return &homePageSearchView
}
