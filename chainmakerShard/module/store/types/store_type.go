/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package types

// State event status
type State byte

// 跨链事物的状态流转
const (
	StateUnknown State = iota
	StateInit
	StateReceived
	StateProofConvertFailed
	StateProofFailed
	StateProofSuccess
	StateExecuteSuccess
	StateExecuteFailed
	StateRollbackSuccess
	StateRollbackFailed
	StateCommitSuccess
	StateCommitFailed
	StateSuccess
	StateFailed
)

// StateDBProvider state db type, contains leveldb temporary
type StateDBProvider string

const (
	LevelDB StateDBProvider = "leveldb"
	Memory  StateDBProvider = "memory"
)
