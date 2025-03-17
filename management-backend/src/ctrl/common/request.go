/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package common

const (
	OffsetDefault = 0
	OffsetMin     = 0
	LimitDefault  = 10
	LimitMax      = 100
)

type RequestBody interface {
	// IsLegal 是否合法
	IsLegal() bool
}

type RangeBody struct {
	PageNum  int64
	PageSize int
}

func (rangeBody *RangeBody) IsLegal() bool {
	if rangeBody.PageSize > LimitMax || rangeBody.PageNum < OffsetMin {
		return false
	}
	return true
}
