/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package common

import (
	"time"
)

const (
	VOTING    = 0
	NO_VOTING = 1
)

const EVM = 5

type TotalNum struct {
	Count int64 `gorm:"column:count"`
}

type CommonIntField struct {
	Id        int64     `gorm:"column:id;AUTO_INCREMENT;PRIMARY_KEY" json:"id"`
	CreatedAt time.Time `gorm:"column:create_at" json:"createAt"`
	UpdatedAt time.Time `gorm:"column:update_at" json:"updateAt"`
}
