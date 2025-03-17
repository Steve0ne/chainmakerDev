/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net

// Peer host utility
type Peer interface {

	//ID return peer id
	ID() string

	// Listen listen and return data to channel
	Listen() (chan Message, error)

	// Write write message to connection
	Write(Message) error

	// Stop close the node
	Stop() error
}
