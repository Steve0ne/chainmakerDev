/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package net_http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	"github.com/gin-gonic/gin"
)

//var listen = "127.0.0.1:8080"

type T_event struct {
	Name string
}

func (e *T_event) GetType() eventproto.EventType {
	return eventproto.TransactionCtxEventType
}

func TestConnection(t *testing.T) {
	//GenTestCert("./")
	t.Log("hello world")
	closeCh := make(chan struct{})
	ginService()
	time.Sleep(2 * time.Second)
	certRootPath := "/Users/leon/Workstation/go/src/hello/cmd/cert/ca/"
	conn, err := NewConnection(&conf.HttpRouterConfig{
		Address: "127.0.0.1:8887",
		HttpTransport: conf.HttpTransport{
			MaxConnection:    1000,
			IdleIConnTimeout: 5,
			EnableTLS:        true,
			Security: &conf.TransportSecurity{
				CAFile:         certRootPath + "ca.pem" + ",/Users/leon/Workstation/go/src/hello/cmd/cert/ca2/ca.pem",
				EnableCertAuth: false,
				CertFile:       certRootPath + "client.pem",
				KeyFile:        certRootPath + "client.key",
			},
		},
	}, WithRetryStrategy(RetryStrategy{
		MaxRetries: 10,
		Interval:   1000,
	}))
	require.Nil(t, err)
	go func() {
		ch, _ := conn.ReadData()
		for msg := range ch {
			resp, ok := msg.(*Message)
			if !ok {
				fmt.Println("it's not a message")
			}
			eve := T_event{}
			if err := resp.Unmarshal(&eve); err != nil {
				t.Error("response unmarshal to obj fail", err)
			}
			//payload, _ := utils.Base64DecodeToBytes(string(msg.GetPayload()))
			//fmt.Printf("event type is %d, marshal type is %d\n", int(payload[0]), int(payload[1]))
			fmt.Printf("the payload of response is %+v\n", eve)
			close(closeCh)
		}
	}()
	req, err := NewRequest(&T_event{Name: "Leon"}, "cross?method=transaction", event.BinaryMarshalType)
	require.NoError(t, err)
	err = conn.WriteData(req)
	require.NoError(t, err)
	select {
	case <-closeCh:
		return

	}
}

func ginService() {
	certRootPath := "/Users/leon/Workstation/go/src/hello/cmd/cert/ca2/"
	webConf := conf.WebConfig{
		Address:   "127.0.0.1",
		Port:      8887,
		EnableTLS: true,
		Security: &conf.TransportSecurity{
			CAFile:         certRootPath + "ca.pem" + ",/Users/leon/Workstation/go/src/hello/cmd/cert/ca/ca.pem",
			EnableCertAuth: true,
			CertFile:       certRootPath + "server.pem",
			KeyFile:        certRootPath + "server.key",
		},
	}
	conf.Config.ListenerConfig = &conf.ListenerConfig{WebConfig: &webConf}
	gin.SetMode(gin.DebugMode)
	ginRouter := gin.Default()
	ginRouter.Handle(http.MethodPost, "/cross", func(ctx *gin.Context) {
		if param, ok := ctx.GetQuery("method"); ok {
			fmt.Println("param is ", param)
			req := &Request{}
			if err := ctx.ShouldBindJSON(req); err == nil {
				fmt.Printf("req is %+v\n", req)
				eve := &T_event{}
				if err := req.Unmarshal(eve); err != nil {
					fmt.Println("request unmarshal to obj fail", err)
				}
				fmt.Printf("request message is %+v\n", eve)
				data, _ := NewMessage(&T_event{Name: "Jankin"}, 0)
				ctx.JSON(http.StatusOK, Response{
					Code:    0,
					Message: "Hello World",
					Data:    data,
				})
				return
			}
		}
	})
	srv := &http.Server{
		Addr:    conf.Config.ListenerConfig.WebConfig.ToUrl(),
		Handler: ginRouter,
	}
	if webConf.EnableTLS {
		err := listenTLS(srv, &webConf)
		if err != nil {
			panic(err)
		}
		return
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("web service error", err)
		}
	}()
}

func listenTLS(srv *http.Server, config *conf.WebConfig) error {
	if config.Security == nil {
		return errors.New("missing tls security configuration")
	}
	security := config.Security
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
		srv.TLSConfig = &tls.Config{
			ClientCAs:  pool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
	}
	go func() {
		// service connections
		if err := srv.ListenAndServeTLS(security.CertFile, security.KeyFile); err != nil && err != http.ErrServerClosed {
			log.Fatal("Web Server TLS Listen:", err)
		}
	}()
	return nil
}
