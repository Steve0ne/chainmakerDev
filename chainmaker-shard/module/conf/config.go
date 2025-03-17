/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"path/filepath"

	"chainmaker.org/chainmaker-cross/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	flagSets         = make([]*pflag.FlagSet, 0) // flag set
	ConfigFilepath   = "./cross_chain.yml"       // common config path
	Config           = &LocalConf{}              // local config instance for global
	BinaryAbsDirPath = ""                        // default release path
)

// InitLocalConfig init local config
func InitLocalConfig(cmd *cobra.Command) error {
	// 1. init config
	config, err := initLocal(cmd)
	if err != nil {
		return err
	}
	// 处理 log config
	logModuleConfigs := config.LogConfig
	for i := 0; i < len(logModuleConfigs); i++ {
		logModuleConfig := logModuleConfigs[i]
		logModuleConfig.FilePath = FinalCfgPath(logModuleConfig.FilePath)
	}
	// 2. set log config
	logger.InitLogConfig(config.LogConfig)
	// 3. set global config and export
	Config = config
	return nil
}

// InitLocalConfigByFilepath init local config by yml file
func InitLocalConfigByFilepath(ymlFile string) (*LocalConf, error) {
	cmViper := viper.New()
	cmViper.SetConfigFile(ymlFile)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}
	config := &LocalConf{}
	if err := cmViper.Unmarshal(config); err != nil {
		return nil, err
	}
	Config = config
	ConfigFilepath = ymlFile

	return config, nil
}

func initLocal(cmd *cobra.Command) (*LocalConf, error) {
	cmViper := viper.New()

	// 1. load the path of the config files
	ymlFile := ConfigFilepath
	if !filepath.IsAbs(ymlFile) {
		// 获取绝对路径
		ymlFile = FinalCfgPath(ymlFile)
		ConfigFilepath = ymlFile
	}

	// 2. load the config file
	cmViper.SetConfigFile(ymlFile)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}

	for _, command := range cmd.Commands() {
		flagSets = append(flagSets, command.PersistentFlags())
		err := cmViper.BindPFlags(command.PersistentFlags())
		if err != nil {
			return nil, err
		}
	}

	// 3. create new CMConfig instance
	config := &LocalConf{}
	if err := cmViper.Unmarshal(config); err != nil {
		return nil, err
	}

	return config, nil
}

// FinalCfgPath check config and return absolute path
func FinalCfgPath(innerCfgPath string) string {
	var finalCfgPath string
	if filepath.IsAbs(innerCfgPath) {
		finalCfgPath = innerCfgPath
	} else {
		finalCfgPath = filepath.Join(BinaryAbsDirPath, innerCfgPath)
	}
	return finalCfgPath
}
