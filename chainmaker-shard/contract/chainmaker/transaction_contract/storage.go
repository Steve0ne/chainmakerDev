/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import "strconv"

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
	ContractName string
	Method       string
	Params       map[string][]byte
}

func callParamsToMap(cp *CallContractParams) map[string]string {
	m := make(map[string]string)
	// check params
	if cp.ContractName == "" || cp.Method == "" || cp.Params == nil {
		return nil
	}
	// marshal params
	paramsItems := ParamsMapToEasyCodecItem(cp.Params)
	paramsStr := string(EasyMarshal(paramsItems))
	// load fields
	m[KeyParams] = paramsStr
	m[KeyContractName] = cp.ContractName
	m[KeyMethod] = cp.Method

	return m
}

func callParamsFromMap(m map[string][]byte) *CallContractParams {
	var cp CallContractParams
	// load params
	if paramsStr, ok := m[KeyParams]; ok {
		paramsItems := EasyUnmarshal(paramsStr)
		params := EasyCodecItemToParamsMap(paramsItems)
		cp.Params = params
	} else {
		return nil
	}
	// load contractName
	if contractName, ok := m[KeyContractName]; ok {
		cp.ContractName = string(contractName)
	} else {
		return nil
	}
	// load method
	if method, ok := m[KeyMethod]; ok {
		cp.Method = string(method)
	} else {
		return nil
	}

	return &cp
}

type Response struct {
	Code int
	Result string
}

func ResponseToJsonString(e *Response) string {
	code := strconv.Itoa(e.Code)
	m := map[string][]byte{
		"Code": []byte(code),
		"Result": []byte(e.Result),
	}
	items := ParamsMapToEasyCodecItem(m)
	return EasyCodecItemToJsonStr(items)
}

// put execute data
func putExecute(crossID string, data map[string]string) ResultCode {
	bytes := ParamsMapToBytes(data)
	return PutStateByte(crossID, FieldExecute, bytes)
}

// get execute data
func getExecute(crossID string) (map[string][]byte, ResultCode) {
	if result, resultCode := GetStateByte(crossID, FieldExecute); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to call getExecute, crossID: " + crossID)
		return nil, ERROR
	} else {
		items := EasyUnmarshal(result)
		params := EasyCodecItemToParamsMap(items)
		// 返回结果
		return params, resultCode
	}
}

// put rollback data
func putRollback(crossID string, data map[string]string) ResultCode {
	bytes := ParamsMapToBytes(data)
	return PutStateByte(crossID, FieldRollback, bytes)
}

// get rollback data
func getRollback(crossID string) (map[string][]byte, ResultCode) {
	if result, resultCode := GetStateByte(crossID, FieldRollback); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to call getRollback, crossID: " + crossID)
		return nil, ERROR
	} else {
		items := EasyUnmarshal(result)
		params := EasyCodecItemToParamsMap(items)
		// 返回结果
		return params, resultCode
	}
}

// put cross state
func putState(crossID string, state State) ResultCode {
	return PutStateByte(crossID, FieldState, []byte(state))
}

// put cross state
func getState(crossID string) State {
	if result, resultCode := GetStateByte(crossID, FieldState); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to call getState, crossID: " + crossID)
		return StateUnknown
	} else {
		// 返回结果
		return State(result)
	}
}

// put cross proof
func putProof(proofKey, txProof string) ResultCode {
	return PutStateByte(proofKey, FieldProof, []byte(txProof))
}

// get cross proof
func getProof(proofKey string) (string, ResultCode) {
	if result, resultCode := GetStateByte(proofKey, FieldProof); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to call getState, getProof: " + proofKey)
		return "", ERROR
	} else {
		// 返回结果
		return string(result), SUCCESS
	}
}

// check cross exist
func isCrossIDExist(crossID string) bool {
	if v, resultCode := GetStateByte(crossID, FieldState); resultCode != SUCCESS {
		return false
	} else {
		if v == nil || len(v) == 0 {
			return false
		} else {
			ErrorResult("duplicated crossID, get value: " + string(v))
			return true
		}
	}
}
