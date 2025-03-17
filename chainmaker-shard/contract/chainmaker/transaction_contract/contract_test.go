/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshalNestMap(t *testing.T) {
	// set hex string
	setHexString := "0100000001000000010000006b020000001600000001000000010000000100000076020000000100000031"
	// sub map
	params := make(map[string][]byte)
	params["v"] = []byte("1")
	// marshal sub map
	cdcs := ParamsMapToEasyCodecItem(params)
	paramsBytes := EasyMarshal(cdcs)
	// nest map
	Params := make(map[string][]byte)
	Params["k"] = paramsBytes
	// marshal nest map
	cdcParams := ParamsMapToEasyCodecItem(Params)
	ParamsBytes := EasyMarshal(cdcParams)
	hexStr := hex.EncodeToString(ParamsBytes)
	fmt.Println(hexStr)
	require.Equal(t, hexStr, setHexString)
}

func TestUnmarshallParams(t *testing.T) {
	// sub params
	p := map[string][]byte{"v": {'1'}}
	cdcs := ParamsMapToEasyCodecItem(p)
	paramsBytes := EasyMarshal(cdcs)
	// pack executeParams
	executeParams := map[string][]byte{
		"contractName": []byte("balance"),
		"method":       []byte("Plus"),
		"params":       paramsBytes,
	}
	executeData := EasyMarshal(ParamsMapToEasyCodecItem(executeParams))
	// pack rollbackParams
	rollbackParams := map[string][]byte{
		"contractName": []byte("balance"),
		"method":       []byte("Reset"),
		"params":       nil,
	}
	rollbackData := EasyMarshal(ParamsMapToEasyCodecItem(rollbackParams))
	// create params
	params := map[string][]byte{
		"crossID":      []byte("test1"),
		"executeData":  executeData,
		"rollbackData": rollbackData,
	}
	crossID, x, y := UnpackUploadParams(ParamsMapToEasyCodecItem(params))
	require.NotEmpty(t, crossID)
	require.NotEmpty(t, x)
	require.NotEmpty(t, y)
}
