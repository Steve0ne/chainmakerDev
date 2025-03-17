/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

const (
	KeyContractName = "contractName"
	KeyMethod       = "method"
	KeyParams       = "params"

	EmptyCrossID = ""
	EmptyParam	 = ""
)

const (
	SUCCESS200 int32 = 200
	SUCCESS int32 = 0
	ERROR   int32 = 1
)

func ToArgs(method string, params []string) [][]byte {
	var bytes [][]byte
	bytes = append(bytes, []byte(method))
	for _, param := range params {
		bytes = append(bytes, []byte(param))
	}
	return bytes
}