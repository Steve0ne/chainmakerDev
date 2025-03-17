/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package grpc_listener

import (
	"net"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GrpcListener - grpc server
type GrpcListener struct {
	server   *grpc.Server       // server
	listener net.Listener       // listener
	log      *zap.SugaredLogger // log
}

// NewGrpcListener create new grpc server listener
func NewGrpcListener() *GrpcListener {
	log := logger.GetLogger(logger.ModuleGrpcListener)
	grpcConf := conf.Config.ListenerConfig.GrpcConfig

	lis, err := net.Listen(grpcConf.Network, grpcConf.Address)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	return &GrpcListener{
		server:   s,
		listener: lis,
		log:      log,
	}
}

// ListenStart start grpc listener
func (gl *GrpcListener) ListenStart() error {
	var err error
	// TODO register methods
	go func() {
		if err = gl.server.Serve(gl.listener); err != nil {
			gl.log.Errorf("failed to start grpc serve: %v", err)
		}
	}()

	return err
}

// Stop stop listener server
func (gl *GrpcListener) Stop() error {
	gl.server.Stop()
	gl.log.Info("Module grpc-listener stopped")
	return nil
}
