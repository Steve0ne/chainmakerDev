/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWebConfig_ToUrl(t *testing.T) {
	wc := WebConfig{
		Address: "127.0.0.1",
		Port:    8080,
	}
	url := wc.ToUrl()
	require.Equal(t, url, "127.0.0.1:8080")
}

func TestLibP2PChannelConfig_GetDelimit(t *testing.T) {
	libp2pCC := &LibP2PChannelConfig{
		Address:     "/ip4/0.0.0.0/tcp/19527",
		PrivKeyFile: "config/ecprikey.key",
		ProtocolID:  "/listener",
		Delimit:     "\n",
	}
	deli := libp2pCC.GetDelimit()
	require.Equal(t, deli, uint8(0xa))

	// test panic method
	libp2pCC = &LibP2PChannelConfig{
		Address:     "/ip4/0.0.0.0/tcp/19527",
		PrivKeyFile: "config/ecprikey.key",
		ProtocolID:  "/listener",
		Delimit:     "\\n",
	}
	defer func() {
		if err := recover(); err != nil {
			require.Equal(t, err, "delimit config more than one rune")
		}
	}()
	deli = libp2pCC.GetDelimit()
}

func TestLibP2PRouterConfig_GetDelimit(t *testing.T) {
	libp2pRC := &LibP2PRouterConfig{
		Address:           "/ip4/IP/tcp/PORT/PEER_ID",
		ProtocolID:        "/listener",
		Delimit:           "\n",
		ReconnectLimit:    1000,
		ReconnectInterval: 5000,
	}
	deli := libp2pRC.GetDelimit()
	require.Equal(t, deli, uint8(0xa))

	// test check method
	libp2pRC = &LibP2PRouterConfig{
		Address:           "/ip4/IP/tcp/PORT/PEER_ID",
		ProtocolID:        "/listener",
		Delimit:           "\\n",
		ReconnectLimit:    1000,
		ReconnectInterval: 5000,
	}
	defer func() {
		if err := recover(); err != nil {
			require.Equal(t, err, "delimit config more than one rune")
		}
	}()
	deli = libp2pRC.GetDelimit()
}

func TestProverConfig_GetChainIDs(t *testing.T) {
	pc := ProverConfig{
		Provider:   "spv",
		ConfigPath: "config/chainmaker/spv.yml",
		ChainIDs:   []string{"chain1", "chain2"},
	}
	ids := pc.GetChainIDs()
	require.Equal(t, ids, []string{"chain1", "chain2"})
}

func TestRouterConfig_GetChainIDs(t *testing.T) {
	libp2pRC := &LibP2PRouterConfig{
		Address:           "/ip4/IP/tcp/PORT/PEER_ID",
		ProtocolID:        "/listener",
		Delimit:           "\n",
		ReconnectLimit:    1000,
		ReconnectInterval: 5000,
	}

	rc := RouterConfig{
		Provider:     "libp2p",
		ChainIDs:     []string{"chain1", "chain2"},
		LibP2PRouter: libp2pRC,
	}

	ids := rc.GetChainIDs()
	require.Equal(t, ids, []string{"chain1", "chain2"})
}
