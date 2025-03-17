/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package conf

import "time"

const (
	TxMsgResultMaxWaitTimeout = time.Second * 30 // 等待结果时间，默认半分钟

	LogWritePeriod = time.Second * 10 // 日志打印周期
)
