/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"chainmaker.org/chainmaker/common/v2/crypto/hash"
	"encoding/hex"
)

const (
	SHA256 = "SHA256"
)

func Sha256(content []byte) ([]byte, error) {
	return hash.GetByStrType(SHA256, content)
}

func Sha256HexString(content []byte) string {
	hash, err := Sha256(content)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(hash)
}
