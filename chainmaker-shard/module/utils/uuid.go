/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"strings"

	"github.com/google/uuid"
)

// NewUUID return random string
func NewUUID() string {
	return GetUUID()
}

// NewRandomKey return random key
func NewRandomKey() string {
	return GetUUID()
}

// getStandardUUID
func getStandardUUID() string {
	return uuid.New().String()
}

// GetUUID return uuid
func GetUUID() string {
	return strings.Replace(getStandardUUID(), "-", "", -1)
}
