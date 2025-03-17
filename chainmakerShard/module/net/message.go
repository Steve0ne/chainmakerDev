/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net

// Message is interface of message
type Message interface {

	// GetNodeID return id in the network
	GetNodeID() string

	// GetPayload return data which will be transfer
	GetPayload() []byte
}
