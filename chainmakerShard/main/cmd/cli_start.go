/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/server"
	"github.com/google/martian/log"
	"github.com/spf13/cobra"
)

// StartCMD start by command for init params
func StartCMD() *cobra.Command {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Startup Cross Server",
		Long:  "Startup ChainMaker Cross Chain Server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			initLocalConfig(cmd)
			start()
			fmt.Println("cross server stopped")
			return nil
		},
	}
	attachFlags(startCmd, []string{flagNameOfConfigFilepath, flagNameOfBinaryDirPath})
	return startCmd
}

// start this is real start function
func start() {
	// get log
	cliLog := logger.GetLogger(logger.ModuleCli)
	// init server
	proxyServer := server.NewServer()
	if err := proxyServer.Start(); err != nil {
		log.Errorf("server start failed, %s", err.Error())
		return
	}

	// new an error channel to receive errors
	errorC := make(chan error, 1)

	// handle exit signal in separate go routines
	go handleExitSignal(errorC)

	// listen error signal in main function
	err := <-errorC
	if err != nil {
		cliLog.Error("server encounters error ", err)
	}
	err = proxyServer.Stop()
	if err != nil {
		cliLog.Error("Stop err: ", err)
	}
	cliLog.Info("All is stopped!")
}

// handleExitSignal listen exit signal for process stop
func handleExitSignal(exitC chan<- error) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, os.Interrupt, syscall.SIGINT)
	defer signal.Stop(signalChan)

	for sig := range signalChan {
		log.Infof("received signal: %d (%s)", sig, sig)
		exitC <- nil
	}
}
