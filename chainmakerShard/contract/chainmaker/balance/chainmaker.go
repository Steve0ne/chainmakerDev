/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

// sdk for user

import (
	"strconv"
	"unsafe"
)

type ResultCode int

const (
	// special parameters passed to contract
	ContractParamCreatorOrgId = "__creator_org_id__"
	ContractParamCreatorRole  = "__creator_role__"
	ContractParamCreatorPk    = "__creator_pk__"
	ContractParamSenderOrgId  = "__sender_org_id__"
	ContractParamSenderRole   = "__sender_role__"
	ContractParamSenderPk     = "__sender_pk__"
	ContractParamBlockHeight  = "__block_height__"
	ContractParamTxId         = "__tx_id__"
	ContractParamContextPtr   = "__context_ptr__"

	// method name used by smart contract sdk
	ContractMethodGetStateLen     = "GetStateLen"
	ContractMethodGetState        = "GetState"
	ContractMethodPutState        = "PutState"
	ContractMethodDeleteState     = "DeleteState"
	ContractMethodSuccessResult   = "SuccessResult"
	ContractMethodErrorResult     = "ErrorResult"
	ContractMethodCallContract    = "CallContract"
	ContractMethodCallContractLen = "CallContractLen"
	ContractMethodEmitEvent       = "EmitEvent"

	SUCCESS ResultCode = 0
	ERROR   ResultCode = 1
)

// sysCall provides data interaction with the chain. sysCallReq common param, request var param
//export sys_call
//func sysCall(requestHeader string, requestBody string) int32

//export log_message
//func logMessage(msg string)

var argsBytes []byte
var argsMap []*EasyCodecItem
var argsFlag bool

//export runtime_type
func runtimeType() int32 {
	var ContractRuntimeGoSdkType int32 = 4
	argsFlag = false
	return ContractRuntimeGoSdkType
}

//export deallocate
func deallocate(size int32) {
	argsBytes = make([]byte, size)
	argsMap = make([]*EasyCodecItem, 0)
	argsFlag = false
}

//export allocate
func allocate(size int32) uintptr {
	argsBytes = make([]byte, size)
	argsMap = make([]*EasyCodecItem, 0)
	argsFlag = false

	return uintptr(unsafe.Pointer(&argsBytes[0]))
}

func getRequestHeader(method string) string {
	var items []*EasyCodecItem
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_SYSTEM,
		Key:       "ctx_ptr",
		ValueType: EasyValueType_INT32,
		Value:     getCtxPtr(),
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_SYSTEM,
		Key:       "version",
		ValueType: EasyValueType_STRING,
		Value:     "v0.7.2",
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_SYSTEM,
		Key:       "method",
		ValueType: EasyValueType_STRING,
		Value:     method,
	})
	return string(EasyMarshal(items))
}

// LogMessage
func LogMessage(msg string) {
	//logMessage(msg)
}

// GetState get state from chain
func GetState(key string, field string) (string, ResultCode) {
	result, code := GetStateByte(key, field)
	if code != SUCCESS {
		return "", code
	}
	return string(result), code
}

// GetState get state from chain
func GetStateByte(key string, field string) ([]byte, ResultCode) {
	// prepare param
	var valueLen int32 = 0
	valuePtr := int(uintptr(unsafe.Pointer(&valueLen)))
	var items []*EasyCodecItem
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "key",
		ValueType: EasyValueType_STRING,
		Value:     key,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "field",
		ValueType: EasyValueType_STRING,
		Value:     field,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "value_ptr",
		ValueType: EasyValueType_STRING,
		Value:     strconv.Itoa(valuePtr),
	})
	//b := EasyMarshal(items)
	//req := string(b)
	// send req get len
	//code := sysCall(getRequestHeader(ContractMethodGetStateLen), req)
	//if code != int32(SUCCESS) {
	//	return nil, ERROR
	//}
	if valueLen == 0 { // nothing found
		return nil, SUCCESS
	}
	// prepare param
	valueByte := make([]byte, valueLen)
	valuePtr = int(uintptr(unsafe.Pointer(&valueByte[0])))
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "key",
		ValueType: EasyValueType_STRING,
		Value:     key,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "field",
		ValueType: EasyValueType_STRING,
		Value:     field,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "value_ptr",
		ValueType: EasyValueType_STRING,
		Value:     strconv.Itoa(valuePtr),
	})
	//b = EasyMarshal(items)
	//req = string(b)
	// send req get value
	//code2 := sysCall(getRequestHeader(ContractMethodGetState), req)
	//if code2 != int32(SUCCESS) {
	//	return nil, ERROR
	//}
	return valueByte, SUCCESS
}

// GetStateFromKey get state from chain
func GetStateFromKey(key string) ([]byte, ResultCode) {
	return GetStateByte(key, "")
}

//EmitEvent emit Event to chain
func EmitEvent(topic string, data ...string) ResultCode {
	// prepare param
	var items []*EasyCodecItem
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "topic\\0",
		ValueType: EasyValueType_STRING,
		Value:     topic,
	})
	var values []*EasyCodecItem
	values = make([]*EasyCodecItem, 0)
	for _, value := range data {
		values = append(values, &EasyCodecItem{
			KeyType:   EasyKeyType_USER,
			Key:       "",
			ValueType: EasyValueType_STRING,
			Value:     value,
		})
	}
	databytes := EasyMarshal(values)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "data\\0",
		ValueType: EasyValueType_STRING,
		Value:     string(databytes),
	})
	b := EasyMarshal(items)
	jsonParam := string(b)
	LogMessage("goContractEventData :" + jsonParam)

	// send req put value
	//code := sysCall(getRequestHeader(ContractMethodEmitEvent), jsonParam)
	//if code != int32(SUCCESS) {
	//	return ERROR
	//}
	return SUCCESS
}

// PutState put state to chain
func PutState(key string, field string, value string) ResultCode {
	// prepare param
	var items []*EasyCodecItem
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "key",
		ValueType: EasyValueType_STRING,
		Value:     key,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "field",
		ValueType: EasyValueType_STRING,
		Value:     field,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "value",
		ValueType: EasyValueType_STRING,
		Value:     value,
	})

	//b := EasyMarshal(items)
	//jsonParam := string(b)
	// send req put value
	//code := sysCall(getRequestHeader(ContractMethodPutState), jsonParam)
	//if code != int32(SUCCESS) {
	//	return ERROR
	//}
	return SUCCESS
}

// PutState put state to chain
func PutStateByte(key string, field string, value []byte) ResultCode {
	return PutState(key, field, string(value))
}

// PutStateFromKey put state to chain
func PutStateFromKey(key string, value string) ResultCode {
	return PutState(key, "", value)
}

// PutStateFromKey put state to chain
func PutStateFromKeyByte(key string, value []byte) ResultCode {
	return PutStateByte(key, "", value)
}

// DeleteState delete state to chain
func DeleteState(key string, field string) ResultCode {
	// prepare param
	var items []*EasyCodecItem
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "key",
		ValueType: EasyValueType_STRING,
		Value:     key,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "field",
		ValueType: EasyValueType_STRING,
		Value:     field,
	})
	//b := EasyMarshal(items)
	//jsonParam := string(b)
	// send req put value
	//code := sysCall(getRequestHeader(ContractMethodDeleteState), jsonParam)
	//if code != int32(SUCCESS) {
	//	return ERROR
	//}
	return SUCCESS
}

// CallContract call other contract from chain
func CallContract(contractName string, method string, param map[string]string) ([]byte, ResultCode) {
	// prepare
	var valueLen int32 = 0
	valuePtr := int(uintptr(unsafe.Pointer(&valueLen)))

	// prepare param
	var items []*EasyCodecItem
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "value_ptr",
		ValueType: EasyValueType_STRING,
		Value:     strconv.Itoa(valuePtr),
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "contract_name",
		ValueType: EasyValueType_STRING,
		Value:     contractName,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "method",
		ValueType: EasyValueType_STRING,
		Value:     method,
	})
	params := ParamsMapToEasyCodecItem(param)
	paramBytes := EasyMarshal(params)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "param",
		ValueType: EasyValueType_BYTES,
		Value:     paramBytes,
	})
	//b := EasyMarshal(items)
	//jsonParam := string(b)
	// send req get call len
	//code := sysCall(getRequestHeader(ContractMethodCallContractLen), jsonParam)
	//if code != int32(SUCCESS) {
	//	return nil, ERROR
	//}
	if valueLen == 0 { // nothing found
		return nil, SUCCESS
	}

	// prepare param
	valueByte := make([]byte, valueLen)
	valuePtr = int(uintptr(unsafe.Pointer(&valueByte[0])))
	items = make([]*EasyCodecItem, 0)
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "value_ptr",
		ValueType: EasyValueType_STRING,
		Value:     strconv.Itoa(valuePtr),
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "contract_name",
		ValueType: EasyValueType_STRING,
		Value:     contractName,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "method",
		ValueType: EasyValueType_STRING,
		Value:     method,
	})
	items = append(items, &EasyCodecItem{
		KeyType:   EasyKeyType_USER,
		Key:       "param",
		ValueType: EasyValueType_BYTES,
		Value:     paramBytes,
	})
	//b = EasyMarshal(items)
	//jsonParam = string(b)
	// send req get value
	//code2 := sysCall(getRequestHeader(ContractMethodCallContract), jsonParam)
	//if code2 != int32(SUCCESS) {
	//	return nil, ERROR
	//}
	return valueByte, SUCCESS
}

func DeleteStateFromKey(key string) ResultCode {
	return DeleteState(key, "")
}

// SuccessResult record success data
func SuccessResult(msg string) {
	//sysCall(getRequestHeader(ContractMethodSuccessResult), msg)
}

// SuccessResult record success data
func SuccessResultByte(msg []byte) {
	//sysCall(getRequestHeader(ContractMethodSuccessResult), string(msg))
}

// ErrorResult record error msg
func ErrorResult(msg string) {
	//sysCall(getRequestHeader(ContractMethodErrorResult), string(msg))
}

func GetCreatorOrgId() (string, ResultCode) {
	return Arg(ContractParamCreatorOrgId)
}
func GetCreatorRole() (string, ResultCode) {
	return Arg(ContractParamCreatorRole)
}
func GetCreatorPk() (string, ResultCode) {
	return Arg(ContractParamCreatorPk)
}
func GetSenderOrgId() (string, ResultCode) {
	return Arg(ContractParamSenderOrgId)
}
func GetSenderRole() (string, ResultCode) {
	return Arg(ContractParamSenderRole)
}
func GetSendePk() (string, ResultCode) {
	return Arg(ContractParamSenderPk)
}
func GetBlockHeight() (string, ResultCode) {
	return Arg(ContractParamBlockHeight)
}
func GetTxId() (string, ResultCode) {
	return Arg(ContractParamTxId)
}
func getCtxPtr() int32 {
	if str, resultCode := Arg(ContractParamContextPtr); resultCode != SUCCESS {
		LogMessage("failed to get ctx ptr")
		return 0
	} else {
		ptr, err := strconv.Atoi(str) //string转int32
		if err != nil {
			LogMessage("get ptr err: " + err.Error())
		}
		return int32(ptr)
	}
}

// cbor
//func getArgsMap() error {
//	if !argsFlag {
//		ctx = vmcbor.RuntimeContext{}
//		ctx.Unmarshal(argsBytes)
//		args := ctx.Call.Args
//		for i := range args {
//			argsMap[args[i].field] = args[i].value
//		}
//		argsFlag = true
//	}
//	return nil
//}

func getArgsMap() error {
	if !argsFlag {
		argsMap = EasyUnmarshal(argsBytes)
		argsFlag = true
	}
	return nil
}

func Arg(key string) (string, ResultCode) {
	err := getArgsMap()
	if err != nil {
		LogMessage("get arg error:" + err.Error())
		return "", ERROR
	}
	for _, v := range argsMap {
		if v.Key == key {
			return v.Value.(string), SUCCESS
		}
	}
	return "", ERROR
}

func Args() []*EasyCodecItem {
	err := getArgsMap()
	if err != nil {
		LogMessage("get Args error:" + err.Error())
	}
	return argsMap
}
