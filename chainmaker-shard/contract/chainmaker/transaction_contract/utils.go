/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

const (
	KeyCrossID      = "crossID"
	KeyExecuteData  = "executeData"
	KeyRollbackData = "rollbackData"

	KeyContractName = "contractName"
	KeyMethod       = "method"
	KeyParams       = "params"

	EmptyCrossID = ""
)

// UnpackUploadParams - check and parse transaction called params
// transaction contract call params format:
// map[crossID] 		= "crossID"
// map[executeData] 	= "executeData" 	// json string
// map[rollbackData] 	= "rollbackData" 	// json string
func UnpackUploadParams(args []*EasyCodecItem) (string, *CallContractParams, *CallContractParams) {
	var crossID = EmptyCrossID
	var eParams, rParams *CallContractParams
	// get crossID
	if v, ok := GetValueFromItems(args, KeyCrossID, EasyKeyType_USER); ok {
		ID, _ := v.([]byte)
		crossID = string(ID)
	}
	// check and parse data
	// parse eParams
	if v, ok := GetValueFromItems(args, KeyExecuteData, EasyKeyType_USER); ok {
		if eParamsBytes, convertOK := v.([]byte); convertOK {
			eParamsItems := EasyUnmarshal(eParamsBytes)
			eMap := EasyCodecItemToParamsMap(eParamsItems)
			eParams = callParamsFromMap(eMap)
		}
	}
	// parse eParams
	if v, ok := GetValueFromItems(args, KeyRollbackData, EasyKeyType_USER); ok {
		if rParamsBytes, convertOK := v.([]byte); convertOK {
			rParamsItems := EasyUnmarshal(rParamsBytes)
			rMap := EasyCodecItemToParamsMap(rParamsItems)
			rParams = callParamsFromMap(rMap)
		}
	}
	return crossID, eParams, rParams
}

func ParamsMapToBytes(params map[string]string) []byte {
	items := make([]*EasyCodecItem, 0)
	for key, value := range params {
		items = append(items, &EasyCodecItem{
			KeyType:   EasyKeyType_USER,
			Key:       key,
			ValueType: EasyValueType_STRING,
			Value:     value,
		})
	}
	itemsBytes := EasyMarshal(items)
	return itemsBytes
}

//GetValueFromItems get value from items
func GetValueFromItems(items []*EasyCodecItem, key string, keyType EasyKeyType) (value interface{}, ok bool) {
	for _, v := range items {
		if value, ok = v.GetValue(key, keyType); ok {
			break
		}
	}
	return
}