/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net

type ConnectionProvider string

type PeerProvider string

const (
	// connection type
	LibP2PConnection ConnectionProvider = "libp2p"
	HttpConnection   ConnectionProvider = "http"
	// peer type
	LibP2PPeer      PeerProvider = "libp2p"
	LibP2PDummyPeer PeerProvider = "libp2p_dummy"
)
