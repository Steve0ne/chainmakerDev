/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net

import (
	"bufio"
)

// Read loop read data from connection
type Read func(*bufio.ReadWriter) (<-chan Message, error)

// Write write bytes to connection
type Write func(*bufio.ReadWriter) error

// Connection is connection handler
type Connection interface {
	//GetProvider return the provider type of the connection,such as http,libp2p etc
	GetProvider() ConnectionProvider

	// PeerID return peer id of the network
	PeerID() string

	// ReadData read data from connection
	ReadData() (chan Message, error)

	// WriteData write data to connection
	WriteData(Message) error

	// Close close the connection
	Close() error
}
