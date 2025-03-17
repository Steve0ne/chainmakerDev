/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package web_listener

import (
	"testing"

	"chainmaker.org/chainmaker-cross/conf"
	"github.com/stretchr/testify/require"
)

func TestListenStart(t *testing.T) {
	// overwrite web config
	webConf := conf.WebConfig{
		Address: "127.0.0.1",
		Port:    12307,
		//EnableTLS: true,
		//Security: &conf.TransportSecurity{
		//	CAFile:         "/Users/leon/Cert/ca.crt",
		//	EnableCertAuth: true,
		//	CertFile:       "/Users/leon/Cert/server.crt",
		//	KeyFile:        "/Users/leon/Cert/server.key",
		//},
	}
	conf.Config.ListenerConfig = &conf.ListenerConfig{WebConfig: &webConf}

	// new listener
	webListener := NewWebListener()

	// start server
	err := webListener.ListenStart()
	require.NoError(t, err)
	// hold service
	//select {}
}

func TestListenStop(t *testing.T) {
	var err error
	// overwrite web config
	webConf := conf.WebConfig{
		Address: "127.0.0.1",
		Port:    12307,
	}
	conf.Config.ListenerConfig = &conf.ListenerConfig{WebConfig: &webConf}

	// new listener
	webListener := NewWebListener()

	// start server
	err = webListener.ListenStart()
	require.NoError(t, err)

	// hold service
	err = webListener.Stop()
	require.NoError(t, err)
}
