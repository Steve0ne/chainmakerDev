/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"reflect"
	"unsafe"
)

func GetPtrUnExportField(s interface{}, filed string) reflect.Value {
	v := reflect.ValueOf(s).Elem().FieldByName(filed)
	// 必须要调用 Elem()
	return reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
}
