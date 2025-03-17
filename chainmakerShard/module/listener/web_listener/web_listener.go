/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package web_listener

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/listener/web_listener/methods"
	"chainmaker.org/chainmaker-cross/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// WebListener is listener of web for CrossSDK
type WebListener struct {
	server *http.Server       // http server
	logger *zap.SugaredLogger // log
	config *conf.WebConfig
}

// NewWebListener create new instance of web listener
func NewWebListener() *WebListener {
	// 启动Web服务(默认Debug级别)
	gin.SetMode(gin.ReleaseMode)
	// 生成route
	ginRouter := gin.Default()
	// 初始化路由配置
	initRouter(ginRouter)
	// new server
	srv := &http.Server{
		Addr:    conf.Config.ListenerConfig.WebConfig.ToUrl(),
		Handler: ginRouter,
		//TLSConfig:
	}

	return &WebListener{
		server: srv,
		logger: logger.GetLogger(logger.ModuleWebListener),
		config: conf.Config.ListenerConfig.WebConfig,
	}
}

// ListenStart web listener server start
func (wl *WebListener) ListenStart() error {
	// 设置logger
	methods.InitHandlers(wl.logger)
	if wl.config.EnableTLS {
		err := wl.listenTLSStart()
		if err != nil {
			wl.logger.Error("Web Server ListenTLS:", err)
		}
		return err
	}
	// 启动Http服务
	go func() {
		// service connections
		if err := wl.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			wl.logger.Error("Web Server Listen:", err)
		}
	}()
	return nil
}

func (wl *WebListener) listenTLSStart() error {
	if wl.config.Security == nil {
		return errors.New("missing tls security configuration")
	}
	security := wl.config.Security
	if security.CertFile == "" || security.KeyFile == "" {
		return errors.New("missing server's cert/key configuration")
	}
	if security.EnableCertAuth {
		caFiles := strings.Split(security.CAFile, ",")
		if len(caFiles) == 0 {
			return errors.New("missing ca files in security config")
		}

		pool := x509.NewCertPool()
		for _, caFile := range caFiles {
			caCrt, err := ioutil.ReadFile(caFile)
			if err != nil {
				return errors.WithMessage(err, fmt.Sprintf("reading ca[%s]", caFile))
			}
			pool.AppendCertsFromPEM(caCrt)
		}
		wl.server.TLSConfig = &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
	}
	go func() {
		// service connections
		if err := wl.server.ListenAndServeTLS(security.CertFile, security.KeyFile); err != nil && err != http.ErrServerClosed {
			wl.logger.Error("Web Server TLS Listen:", err)
		}
	}()
	return nil
}

func (wl *WebListener) checkConfig() error {
	if wl.config.EnableTLS {

	}
	return nil
}

func initRouter(router *gin.Engine) {
	group := router.Group("/")
	initControllers(group) // 定义接口
}

// initControllers 初始化Controller配置
func initControllers(routeGroup *gin.RouterGroup) {
	routeGroup.POST(methods.CrossTag, methods.Dispatch)
	//routeGroup.GET(methods.CrossTag, func(ctx *gin.Context) {
	//	ctx.JSON(http.StatusOK, "hello world!!!")
	//})
}

// Stop web listener server stop
func (wl *WebListener) Stop() error {
	// delay
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := wl.server.Shutdown(ctx); err != nil {
		wl.logger.Error("Web Server Shutdown:", err)
	}
	wl.logger.Info("Module web-listener stopped")
	return nil
}
