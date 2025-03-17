/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"github.com/spf13/viper"
)

//InitConfigByFilepath init the config of cross sdk
func InitConfigByFilepath(ymlFile string) (*ConfigMap, error) {
	cmViper := viper.New()
	cmViper.SetConfigFile(ymlFile)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}
	configL := &Config{}
	if err := cmViper.Unmarshal(configL); err != nil {
		return nil, err
	}
	configM := configL.ConvertToMap()

	return configM, nil
}

func LoadConfig(ymlFile string) (*Config, error) {
	cmViper := viper.New()
	cmViper.SetConfigFile(ymlFile)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}
	configL := &Config{}
	if err := cmViper.Unmarshal(configL); err != nil {
		return nil, err
	}
	configL.Init()
	return configL, nil
}
