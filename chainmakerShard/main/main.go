/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"

	"chainmaker.org/chainmaker-cross/main/cmd"
	"github.com/spf13/cobra"
)

func main() {
	mainCmd := &cobra.Command{Use: "start"}
	mainCmd.AddCommand(cmd.StartCMD())

	err := mainCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("cross proxy start error, %v", err)
	}
}
