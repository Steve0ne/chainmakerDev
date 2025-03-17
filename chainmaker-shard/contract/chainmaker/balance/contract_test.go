/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshal(t *testing.T) {
	// set marshalled string
	setParamHex := "01000000010000000100000076010000000100000031"
	// param map
	params := make(map[string]string)
	params["v"] = "1"
	// marshal
	cdcs := ParamsMapToEasyCodecItem(params)
	paramsBytes := EasyMarshal(cdcs)
	hexStr := hex.EncodeToString(paramsBytes)
	// compare
	require.Equal(t, hexStr, setParamHex)
}
