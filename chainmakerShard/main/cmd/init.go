/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"fmt"
	"path/filepath"

	"chainmaker.org/chainmaker-cross/conf"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	flagNameOfConfigFilepath          = "conf"
	flagNameShortHandOfConfigFilepath = "c"
	flagNameOfBinaryDirPath           = "dir"
	flagNameShortHandOfBinaryDirPath  = "d"
)

// initLocalConfig init local config
func initLocalConfig(cmd *cobra.Command) {
	if err := conf.InitLocalConfig(cmd); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// initFlagSet init flag set
func initFlagSet() *pflag.FlagSet {
	dir := filepath.Dir(".")
	flags := &pflag.FlagSet{}
	flags.StringVarP(&conf.ConfigFilepath, flagNameOfConfigFilepath, flagNameShortHandOfConfigFilepath, conf.ConfigFilepath, "specify config file path, if not set, default use ./cross_chain.yml")
	flags.StringVarP(&conf.BinaryAbsDirPath, flagNameOfBinaryDirPath, flagNameShortHandOfBinaryDirPath, dir, "specify binary dir path, if not set, default use filepath.Dir(.)")
	return flags
}

// attachFlags
func attachFlags(cmd *cobra.Command, flagNames []string) {
	flags := initFlagSet()
	cmdFlags := cmd.Flags()
	for _, flagName := range flagNames {
		if flag := flags.Lookup(flagName); flag != nil {
			cmdFlags.AddFlag(flag)
		}
	}
}
