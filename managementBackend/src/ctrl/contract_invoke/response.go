/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package contract_invoke

import (
	"github.com/emirpasic/gods/lists/arraylist"

	dbcommon "management_backend/src/db/common"
)

type InvokeRecordListView struct {
	Id           int64
	UserName     string
	OrgName      string
	ContractName string
	TxStatus     int
	Status       int
	TxId         string
	CreateTime   int64
}

func NewInvokeRecordListView(invokeRecords []*dbcommon.InvokeRecords) []interface{} {
	invokeRecordsViews := arraylist.New()
	for _, invokeRecord := range invokeRecords {
		invokeRecordView := InvokeRecordListView{
			Id:           invokeRecord.Id,
			UserName:     invokeRecord.UserName,
			OrgName:      invokeRecord.OrgName,
			ContractName: invokeRecord.ContractName,
			TxStatus:     invokeRecord.TxStatus,
			Status:       invokeRecord.Status,
			TxId:         invokeRecord.TxId,
			CreateTime:   invokeRecord.CreatedAt.Unix(),
		}
		invokeRecordsViews.Add(invokeRecordView)
	}

	return invokeRecordsViews.Values()
}

type InvokeContractListView struct {
	ContractName string
	ContractId   int64
}

func NewInvokeContractListView(contracts []*dbcommon.Contract) []interface{} {
	contractViews := arraylist.New()
	for _, contract := range contracts {
		contractView := InvokeContractListView{
			ContractName: contract.Name,
			ContractId:   contract.Id,
		}
		contractViews.Add(contractView)
	}

	return contractViews.Values()
}
