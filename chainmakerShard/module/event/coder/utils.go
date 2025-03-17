/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package coder

import (
	"errors"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
	jsoniter "github.com/json-iterator/go"
)

const (
	EventTyIndex = iota
	MarshalTyIndex
	MinLength = 2
)

var (
	dataNotRightErr = errors.New("this data is not right")
	jsonCoder       = jsoniter.ConfigCompatibleWithStandardLibrary
)

// JsonMarshal marshal object to byte array
func JsonMarshal(v interface{}) ([]byte, error) {
	bytes, err := jsonCoder.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// JsonBinaryMarshal marshal object to byte array
func JsonBinaryMarshal(eventTy eventproto.EventType, v interface{}) ([]byte, error) {
	// 对象转换为json字符串
	bytes, err := JsonMarshal(v)
	if err != nil {
		return nil, err
	}
	// 在字节数组前面增加对应的类型
	typeBytes := []byte{byte(eventTy), byte(event.BinaryMarshalType)}
	totalBytes := append(typeBytes, bytes...)
	return totalBytes, nil
}

// JsonBinaryUnmarshal unmarshal to event by byte array and event-type
func JsonBinaryUnmarshal(data []byte, eveTyByte byte, eve event.Event) error {
	// 首先判断前两个字节
	if len(data) < MinLength {
		return dataNotRightErr
	}
	dataEveTyByte := data[EventTyIndex]
	if dataEveTyByte != eveTyByte {
		return dataNotRightErr
	}
	marshalTypeByte := data[MarshalTyIndex]
	if marshalTypeByte != byte(event.BinaryMarshalType) {
		return dataNotRightErr
	}
	jsonBytes := data[MinLength:]
	return jsonCoder.Unmarshal(jsonBytes, eve)
}
