/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type State string

const (
	// unknown state
	StateUnknown State = "Unknown"
	// execute state
	ExecuteSuccess State = "ExecuteSuccess"
	ExecuteFail    State = "ExecuteFail"
	// commit state
	CommitSuccess State = "CommitSuccess"
	CommitFail    State = "CommitFail"
	// rollback state
	RollbackSuccess State = "RollbackSuccess"
	RollbackFail    State = "RollbackFail"
	RollbackIgnore  State = "RollbackIgnore" // 由于没有 crossID 的回滚，忽略该步骤
)

const (
	FieldExecute  = "Execute"
	FieldRollback = "Rollback"
	FieldState    = "State"
	FieldProof    = "Proof"
)

// necessary params to call a contract
type CallContractParams struct {
	ContractName string 	`json:"contract_name"`
	Method       string		`json:"method"`
	Params       []string	`json:"params,omitempty" metadata:",optional"`
}

type Response struct {
	Code int
	Result string
}

// put execute data
func putExecute(ctx contractapi.TransactionContextInterface, crossID string, data []byte) error {
	return putStateByte(ctx, crossID, FieldExecute, data)
}

// get execute data
func getExecute(ctx contractapi.TransactionContextInterface, crossID string) (*CallContractParams, error) {
	if result, err := getStateByte(ctx, crossID, FieldExecute); err != nil {
		// 返回结果
		return nil, fmt.Errorf("failed to call getExecute, crossID: " + crossID)
	} else {
		m := &CallContractParams{}
		err := json.Unmarshal(result, &m)
		if err != nil {
			return nil, err
		}
		// 返回结果
		return m, nil
	}
}

// put rollback data
func putRollback(ctx contractapi.TransactionContextInterface, crossID string, data []byte) error {
	return putStateByte(ctx, crossID, FieldRollback, data)
}

// get rollback data
func getRollback(ctx contractapi.TransactionContextInterface, crossID string) (*CallContractParams, error) {
	if result, err := getStateByte(ctx, crossID, FieldRollback); err != nil {
		// 返回结果
		return nil, fmt.Errorf("failed to call getRollback, crossID: " + crossID)
	} else {
		m := &CallContractParams{}
		err := json.Unmarshal(result, &m)
		if err != nil {
			return nil, err
		}
		// 返回结果
		return m, nil
	}
}

// put cross state
func putState(ctx contractapi.TransactionContextInterface, crossID string, state State) error {
	return putStateByte(ctx, crossID, FieldState, []byte(state))
}

// put cross state
func getState(ctx contractapi.TransactionContextInterface, crossID string) (State, error) {
	if result, err := getStateByte(ctx, crossID, FieldState); err != nil {
		// 返回结果
		return StateUnknown, err
	} else {
		// 返回结果
		return State(result), nil
	}
}

// put cross proof
func putProof(ctx contractapi.TransactionContextInterface, proofKey, txProof string) error {
	return putStateByte(ctx, proofKey, FieldProof, []byte(txProof))
}

// get cross proof
func getProof(ctx contractapi.TransactionContextInterface, proofKey string) (string, error) {
	if result, err := getStateByte(ctx, proofKey, FieldProof); err != nil {
		// 返回结果
		return "", err
	} else {
		// 返回结果
		return string(result), nil
	}
}

// check cross exist
func isCrossIDExist(ctx contractapi.TransactionContextInterface, crossID string) bool {
	if v, err := getStateByte(ctx, crossID, FieldState); err != nil {
		return false
	} else {
		if v == nil || len(v) == 0 {
			return false
		} else {
			return true
		}
	}
}

func putStateByte(ctx contractapi.TransactionContextInterface, key, field string, bytes []byte) error {
	return ctx.GetStub().PutState(key+field, bytes)
}

func getStateByte(ctx contractapi.TransactionContextInterface, key, field string) ([]byte, error) {
	return ctx.GetStub().GetState(key+field)
}
