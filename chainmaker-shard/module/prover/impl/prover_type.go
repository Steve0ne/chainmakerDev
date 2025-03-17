/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package impl

type ProverType int

// 跨链代理转发消息的证名类型
const (
	// 0验证证明
	TrustProverType ProverType = iota
	// spv验证证明
	SpvProverType
)
