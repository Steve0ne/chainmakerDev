/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"
)

const (
	RandomRange = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
)

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(data string) []byte {
	decodeBytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}
	return decodeBytes
}

func RandomString(len int) string {
	var container string
	b := bytes.NewBufferString(RandomRange)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(RandomRange[randomInt.Int64()])
	}
	return container
}

func CurrentMillSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func CurrentSeconds() int64 {
	return time.Now().UnixNano() / 1e9
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func ConvertToPercent(percent float64) string {
	doubleDecimal := fmt.Sprintf("%.2f", 100*percent)
	// Strip trailing zeroes
	for doubleDecimal[len(doubleDecimal)-1] == '0' {
		doubleDecimal = doubleDecimal[:len(doubleDecimal)-1]
	}
	// Strip the decimal point if it's trailing.
	if doubleDecimal[len(doubleDecimal)-1] == '.' {
		doubleDecimal = doubleDecimal[:len(doubleDecimal)-1]
	}
	return fmt.Sprintf("%s%%", doubleDecimal)
}

func GetHostFromAddress(addr string) string {
	s := strings.Split(addr, ":")
	return s[0]
}
