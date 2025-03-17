/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package net_libp2p

import (
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrivateKeyFromPEM(t *testing.T) {
	var err error
	// 创建临时文件目录
	p, err := ioutil.TempDir(os.TempDir(), "mock")
	require.NoError(t, err)
	// 删除目录
	defer os.RemoveAll(p)
	path := p + "/PEM"
	// 写入私钥
	err = ioutil.WriteFile(path, []byte(privatePEMStr), 0600)
	require.NoError(t, err)
	// 读取私钥
	privateKey, err := prepareKey(path)
	require.NoError(t, err)
	// check private key bytes
	bz, err := privateKey.Bytes()
	res := hex.EncodeToString(bz)
	require.Equal(t, res, privateKeyHex)
}

func TestPrivateKeyFromPEMFilePath(t *testing.T) {
	//var err error
	//// 读取私钥
	//privateKey, err := prepareKey("/root/chainmaker-cross-chain/BJ2020.key")
	//require.NoError(t, err)
	//// check private key bytes
	//bz, err := privateKey.Bytes()
	//_ = hex.EncodeToString(bz)
	//require.NoError(t, err)
}
