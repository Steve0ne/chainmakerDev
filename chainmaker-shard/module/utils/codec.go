/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package utils

import "encoding/base64"

// Base64EncodeToString encode byte array to base64 string
func Base64EncodeToString(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64DecodeToBytes decode base64 string to byte array
func Base64DecodeToBytes(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
