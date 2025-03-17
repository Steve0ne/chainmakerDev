/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	"chainmaker.org/chainmaker-cross/sdk/builder"

	"encoding/json"

	"chainmaker.org/chainmaker-cross/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	mainCmd := &cobra.Command{Use: "cross-chain-cli"}
	mainCmd.AddCommand(DeliverEventCMD())
	mainCmd.AddCommand(ShowCrossResultCMD())

	err := mainCmd.Execute()
	if err != nil {
		fmt.Printf("cross-chain-cli error, %v\n", err)
	}
	return
}

// DeliverEventCMD deliver event command
func DeliverEventCMD() *cobra.Command {
	deliverCmd := &cobra.Command{
		Use:   "deliver",
		Short: "Deliver CrossEvent",
		Long:  "Deliver CrossEvent To Proxy",
		RunE:  DeliverEventRun,
	}
	attachFlags(deliverCmd, []string{flagNameOfConfigFilepath, flagNameOfParams, flagNameOfUrl})
	deliverCmd.Flags().String(flagNameOfParams, "", "the parameters for cross tx")
	return deliverCmd
}

func DeliverEventRun(cmd *cobra.Command, _ []string) error {
	cmViper := viper.New()
	paramPath, err := cmd.Flags().GetString(flagNameOfParams)
	if err != nil {
		return fmt.Errorf("missing flag --params")
	}
	cmViper.SetConfigFile(paramPath)
	if err := cmViper.ReadInConfig(); err != nil {
		return fmt.Errorf("read params error: %v", err)
	}
	crossParams := &CrossTxParams{}
	if err := cmViper.Unmarshal(crossParams); err != nil {
		return fmt.Errorf("unmarshal params error: %v", err)
	}
	crossSDK, err := sdk.NewCrossSDK(sdk.WithConfigFile(ConfigFilepath))
	if err != nil {
		return err
	}
	params := make([]*sdk.CrossTxBuildCtx, len(crossParams.Params))
	for i, crossParam := range crossParams.Params {
		params[i] = sdk.NewCrossTxBuildCtx(
			crossParam.ChainID, crossParam.Index,
			builder.NewContract(crossParam.ContractName, crossParam.ExecuteMethod, crossParam.ExecuteParams.ToBuilderParams()),
			builder.NewContract(crossParam.ContractName, crossParam.RollbackMethod, crossParam.RollbackParams.ToBuilderParams()))
	}
	crossEvent, err := crossSDK.GenCrossEvent(params...)
	if err != nil {
		return err
	}
	resp, err := crossSDK.SendCrossEvent(crossEvent, DefaultURL, false)
	if err != nil {
		return fmt.Errorf("send cross tx error: [%v]", err)
	}
	if resp.CrossId == "" {
		return fmt.Errorf("deliver tx error, remote server may stoped")
	}
	fmt.Println(resp.CrossId)
	return nil
}

// ShowCrossResultCMD show cross result command
func ShowCrossResultCMD() *cobra.Command {
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "Show CrossEvent Result",
		Long:  "Show CrossEvent Result By Proxy",
		RunE: func(cmd *cobra.Command, _ []string) error {
			crossID, err := cmd.Flags().GetString(flagNameOfCrossID)
			if err != nil {
				return fmt.Errorf("missing flag --crossID")
			}
			// 使用 crossID 查询跨链执行结果
			resp, err := sdk.NewCrossSearchEvent(crossID).Query(DefaultURL)
			if err != nil {
				return fmt.Errorf("show cross result error: [%s]", err.Error())
			}
			respStr, err := json.Marshal(resp)
			if err != nil {
				return nil
			}
			fmt.Printf("%s", respStr)
			return nil
		},
	}
	attachFlags(showCmd, []string{flagNameOfCrossID, flagNameOfUrl})
	showCmd.Flags().String(flagNameOfCrossID, "", "the cross id for event")
	return showCmd
}

func initFlagSet() *pflag.FlagSet {
	flags := &pflag.FlagSet{}
	flags.StringVarP(&ConfigFilepath, flagNameOfConfigFilepath, flagNameShortHandOfConfigFilepath, ConfigFilepath, "specify config file path, if not set, default use ./cross_chain_sdk.yml")
	flags.StringVarP(&DefaultURL, flagNameOfUrl, flagNameShortHandOfUrl, DefaultURL, "specify default url, if not set, default use http://localhost:8080")
	return flags
}

func attachFlags(cmd *cobra.Command, flagNames []string) {
	flags := initFlagSet()
	cmdFlags := cmd.Flags()
	for _, flagName := range flagNames {
		if flag := flags.Lookup(flagName); flag != nil {
			cmdFlags.AddFlag(flag)
		}
	}
}
