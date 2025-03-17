/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package listener

import (
	"io/ioutil"
	"os"
	"testing"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/handler"
	"github.com/stretchr/testify/require"
)

const privatePEMStr = "-----BEGIN EC PRIVATE KEY-----\n" +
	"MHcCAQEEIKwz2DjJZWtxnHNtMroZEjr28dqihVAPcyrmQp/5XeGRoAoGCCqGSM49\n" +
	"AwEHoUQDQgAERrAT6tdtd95gd6zh2/qics3s/h7eUtAGLdPhuGstl1j9iqSXI/0i\n" +
	"9qyw+HHd2IkkLdQX6ivmaovGS0vbSDpRCA==\n" +
	"-----END EC PRIVATE KEY-----\n"

func TestListenerManger(t *testing.T) {
	// prepare environment
	path := preparePEM(t)
	defer os.RemoveAll(path)
	// setup config
	conf.Config.ListenerConfig = &conf.ListenerConfig{
		WebConfig: &conf.WebConfig{
			Address: "120.0.0.1",
			Port:    8080},
		ChannelConfig: &conf.ChannelConfig{
			Provider: "libp2p",
			LibP2PChannel: &conf.LibP2PChannelConfig{
				Address:     "/ip4/0.0.0.0/tcp/19527",
				PrivKeyFile: path,
				ProtocolID:  "/listener",
				Delimit:     "\n"},
		},
	}
	// new listener manager
	lm := InitListener()
	handler.InitEventHandlers(nil, nil)
	// Test init listeners
	lm.InitListeners()
	// Test start
	err := lm.Start()
	require.NoError(t, err)
	// Test stop
	err = lm.Stop()
	require.NoError(t, err)
}

func preparePEM(t *testing.T) string {
	var err error
	// 创建临时文件目录
	p, err := ioutil.TempDir(os.TempDir(), "mock")
	require.NoError(t, err)
	path := p + "/PEM"
	// 写入私钥
	err = ioutil.WriteFile(path, []byte(privatePEMStr), 0600)
	require.NoError(t, err)
	return path
}
