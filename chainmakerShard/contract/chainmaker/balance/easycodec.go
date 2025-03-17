/*
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"sort"
	"strconv"
)

type EasyKeyType int32
type EasyValueType int32

const (
	EasyKeyType_SYSTEM EasyKeyType = 0
	EasyKeyType_USER   EasyKeyType = 1

	EasyValueType_INT32  EasyValueType = 0
	EasyValueType_STRING EasyValueType = 1
	EasyValueType_BYTES  EasyValueType = 2
)

type EasyCodecItem struct {
	KeyType EasyKeyType
	Key     string

	ValueType EasyValueType
	Value     interface{}
}

//ParamsMapToEasyCodecItem Params map converter
func ParamsMapToEasyCodecItem(params map[string]string) []*EasyCodecItem {
	keys := make([]string, 0)
	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	items := make([]*EasyCodecItem, 0)

	for _, key := range keys {

		var easyCodecItem EasyCodecItem

		easyCodecItem.KeyType = EasyKeyType_USER

		easyCodecItem.Key = key

		easyCodecItem.ValueType = EasyValueType_STRING

		easyCodecItem.Value = params[key]

		items = append(items, &easyCodecItem)

	}

	return items
}

//EasyCodecItemToParamsMap easyCodecItem converter
func EasyCodecItemToParamsMap(items []*EasyCodecItem) map[string]string {
	params := make(map[string]string)
	for _, item := range items {
		params[item.Key] = item.Value.(string)
	}
	return params
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

//EasyCodecItemToJsonStr easyCodecItem converter
func EasyCodecItemToJsonStr(items []*EasyCodecItem) string {
	str := "{"
	for _, item := range items {
		var val string
		switch item.ValueType {
		case EasyValueType_INT32:
			val = strconv.Itoa(int(item.Value.(int32)))
		case EasyValueType_STRING:
			val = item.Value.(string)
		case EasyValueType_BYTES:
			val = string(item.Value.([]byte))
		}
		key := item.Key
		str = str + "\"" + key + "\":\"" + val + "\","
	}
	if len(str) > 1 {
		str = str[:len(str)-1]
	}
	str = str + "}"
	return str
}

//GetValue get value from item
func (e *EasyCodecItem) GetValue(key string, keyType EasyKeyType) (interface{}, bool) {
	if e.KeyType == keyType && e.Key == key {
		return e.Value, true
	}
	return "", false
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

//SetValue set value from item
func (e *EasyCodecItem) SetValue(key string, keyType EasyKeyType, value interface{}) bool {
	if e.KeyType == keyType && e.Key == key {
		e.Value = value
		return true
	}
	return false
}

//SetValueFromItems get value from items
func SetValueFromItems(items []*EasyCodecItem, key string, keyType EasyKeyType, value interface{}) (ok bool) {
	for i := 0; i < len(items); i++ {
		if ok = items[i].SetValue(key, keyType, value); ok {
			return true
		}
	}
	return false
}

//EasyMarshal serialize item into binary
func EasyMarshal(items []*EasyCodecItem) []byte {
	buf := new(bytes.Buffer)
	uint32DataBytes := make([]byte, 4)

	binaryUint32Marshal(buf, uint32(len(items)), uint32DataBytes)

	for _, item := range items {

		if item.KeyType != EasyKeyType_SYSTEM && item.KeyType != EasyKeyType_USER {
			continue
		}

		binaryUint32Marshal(buf, uint32(item.KeyType), uint32DataBytes)
		binaryUint32Marshal(buf, uint32(len(item.Key)), uint32DataBytes)
		buf.Write([]byte(item.Key))

		switch item.ValueType {

		case EasyValueType_INT32:

			binaryUint32Marshal(buf, uint32(item.ValueType), uint32DataBytes)
			binaryUint32Marshal(buf, uint32(4), uint32DataBytes)
			binaryUint32Marshal(buf, uint32(item.Value.(int32)), uint32DataBytes)

		case EasyValueType_STRING:

			binaryUint32Marshal(buf, uint32(item.ValueType), uint32DataBytes)
			binaryUint32Marshal(buf, uint32(len(item.Value.(string))), uint32DataBytes)
			buf.Write([]byte(item.Value.(string)))

		case EasyValueType_BYTES:

			binaryUint32Marshal(buf, uint32(item.ValueType), uint32DataBytes)
			binaryUint32Marshal(buf, uint32(len(item.Value.([]byte))), uint32DataBytes)
			buf.Write(item.Value.([]byte))

		}
	}

	return buf.Bytes()
}

//EasyUnmarshal Deserialized from binary to item
func EasyUnmarshal(data []byte) []*EasyCodecItem {
	buf := bytes.NewBuffer(data)
	uint32DataBytes := make([]byte, 4)
	count := binaryUint32Unmarshal(buf, uint32DataBytes)

	var (
		items         []*EasyCodecItem
		easyKeyType   EasyKeyType
		keyLength     int32
		keyContent    []byte
		easyValueType EasyValueType
		valueLength   int32
	)

	for i := 0; i < int(count); i++ {
		// Key Part
		easyKeyType = EasyKeyType(binaryUint32Unmarshal(buf, uint32DataBytes))

		keyLength = int32(binaryUint32Unmarshal(buf, uint32DataBytes))
		keyContent = make([]byte, keyLength)
		buf.Read(keyContent)

		// Value Part
		easyValueType = EasyValueType((binaryUint32Unmarshal(buf, uint32DataBytes)))

		valueLength = int32(binaryUint32Unmarshal(buf, uint32DataBytes))

		var easyCodecItem EasyCodecItem

		switch easyValueType {

		case EasyValueType_INT32:

			valueContent := int32(binaryUint32Unmarshal(buf, uint32DataBytes))
			easyCodecItem.Value = valueContent

		case EasyValueType_STRING:

			valueContent := make([]byte, valueLength)
			buf.Read(valueContent)
			easyCodecItem.Value = string(valueContent)

		case EasyValueType_BYTES:

			valueContent := make([]byte, valueLength)
			buf.Read(valueContent)
			easyCodecItem.Value = valueContent

		}

		easyCodecItem.KeyType = easyKeyType
		easyCodecItem.Key = string(keyContent)
		easyCodecItem.ValueType = easyValueType

		items = append(items, &easyCodecItem)
	}

	return items
}

func binaryUint32Marshal(buf *bytes.Buffer, data uint32, dataBytes []byte) {
	_ = dataBytes[3]
	dataBytes[0] = byte(data)
	dataBytes[1] = byte(data >> 8)
	dataBytes[2] = byte(data >> 16)
	dataBytes[3] = byte(data >> 24)
	buf.Write(dataBytes)
}

func binaryUint32Unmarshal(buf *bytes.Buffer, bs []byte) uint32 {
	buf.Read(bs)
	_ = bs[3]
	return uint32(bs[0]) | uint32(bs[1])<<8 | uint32(bs[2])<<16 | uint32(bs[3])<<24
}
