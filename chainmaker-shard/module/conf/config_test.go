/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"chainmaker.org/chainmaker-cross/logger"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestInitLocalConfigByFilepath(t *testing.T) {
	var storageTest = StorageConfig{
		Provider: "leveldb",
		LevelDB: &LevelDBConfig{
			StorePath:       "./testdata",
			WriteBufferSize: 4,
			BloomFilterBits: 10,
		},
	}
	dir, err := os.Getwd()
	require.Nil(t, err)
	ymlFile := getParentDirectory(getParentDirectory(dir)) + string(os.PathSeparator) + "config" +
		string(os.PathSeparator) + "cross_chain.yml"
	localConf, err := InitLocalConfigByFilepath(ymlFile)
	require.Nil(t, err)
	// 检查存储配置
	storageConfig := localConf.StorageConfig
	err = checkStorageConfig(&storageTest, storageConfig)
	require.Nil(t, err)
}

func TestInitLocalConfig(t *testing.T) {
	dir, err := os.Getwd()
	require.Nil(t, err)
	// test absolute path
	ConfigFilepath = getParentDirectory(getParentDirectory(dir)) + string(os.PathSeparator) + "config" +
		string(os.PathSeparator) + "cross_chain.yml"
	err = InitLocalConfig(&cobra.Command{})
	require.Nil(t, err)
	// test path
	ConfigFilepath = "./cross_chain.yml"
	err = InitLocalConfig(&cobra.Command{})
	require.NotNil(t, err)
}

func TestStringToByte(t *testing.T) {
	line := "\n"
	toByte := convertToByte(line)
	if toByte != 10 {
		t.Error("line's answer is error")
	}
}

func TestBasePath(t *testing.T) {
	cfgPath := "../config/cross_chain.yml"
	dir := filepath.Dir(cfgPath)
	absPath := filepath.Join(dir, "adapters/chainmaker_sdk.yml")
	fmt.Println(absPath)
}

func TestILog(t *testing.T) {
	absPath, err := filepath.Abs("../../config/cross_chain.yml")
	localConf, err := InitLocalConfigByFilepath(absPath)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	logger.InitLogConfig(localConf.LogConfig)
	log := logger.GetLogger("[Test]")
	log.Info("this is INFO")
	log.Infof("this is INFOF -> %s", "ChainCross")
	log.Error("this is ERROR")
	log.Error(errors.New("this is new ERROR"))
	log.Error("this is new new ERROR, ", errors.New("ChainCross"))
	log.Errorf("this is ERRORF -> %v", errors.New("ChainCross"))
}

func TestFinalCfgPath(t *testing.T) {
	// test path
	BinaryAbsDirPath = "/path"
	p := FinalCfgPath("test")
	require.Equal(t, p, "/path/test")

	// test absolute path
	p = FinalCfgPath("/test")
	require.Equal(t, p, "/test")
}

func convertToByte(content string) byte {
	bs := []byte(content)
	return bs[0]
}

func checkStorageConfig(testData, confData *StorageConfig) error {
	if testData.Provider != confData.Provider {
		return errors.New("db's provider is not equal")
	}
	if testData.LevelDB == nil || confData.LevelDB == nil {
		return errors.New("data is nil")
	}
	return nil
}

func getParentDirectory(directory string) string {
	return substr(directory, 0, strings.LastIndex(directory, "/"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
