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
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/http2"

	"chainmaker.org/chainmaker-cross/conf"
)

var (
	transportMu             sync.Mutex
	transports              = make(map[transportKey]*http.Transport)
	defaultMaxConns         = 100
	defaultIdleIConnTimeout = 90
)

type transportKey struct {
	MaxConnection   int
	IdleConnTimeout int
	security        conf.TransportSecurity
	useH2           bool
}

func GetTransport(configIn *conf.HttpTransport) (tr *http.Transport, err error) {
	config := *configIn
	if config.MaxConnection == 0 {
		config.MaxConnection = defaultMaxConns
	}
	if config.IdleIConnTimeout == 0 {
		config.IdleIConnTimeout = defaultIdleIConnTimeout
	}

	transportMu.Lock()
	defer transportMu.Unlock()
	key := genTransportKey(&config)
	tr, ok := transports[key]
	if ok {
		return
	}
	tr, err = genTransport(&config)
	if err != nil {
		return
	}
	transports[key] = tr
	return
}

func genTransport(config *conf.HttpTransport) (tr *http.Transport, err error) {
	var tlsConfig *tls.Config
	if config.EnableTLS {
		tlsConfig, err = genClientTlsConfig(config.Security)
		if err != nil {
			return
		}
	}

	tr = &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    config.MaxConnection,
		IdleConnTimeout: time.Duration(config.IdleIConnTimeout) * time.Second,
	}
	if config.EnableH2 {
		err = http2.ConfigureTransport(tr)
		if err != nil {
			return nil, err
		}
	}
	return
}

func genClientTlsConfig(security *conf.TransportSecurity) (*tls.Config, error) {
	if security == nil {
		return nil, errors.New("tls config is nil")
	}

	caFiles := strings.Split(security.CAFile, ",")
	if security.EnableCertAuth && len(caFiles) == 0 {
		return nil, errors.New("there is no ca files in security config")
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: !security.EnableCertAuth,
	}
	//如果配置了客户端证书/私钥,则放入请求中，用于服务端验证
	if security.CertFile != "" || security.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(security.CertFile, security.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	//进行ca证书解析，放入证书池
	pool := x509.NewCertPool()
	for _, caFile := range caFiles {
		caCrt, err := ioutil.ReadFile(caFile)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("reading ca[%s]", caFile))
		}
		pool.AppendCertsFromPEM(caCrt)
	}
	tlsConfig.RootCAs = pool

	return tlsConfig, nil
}

func genTransportKey(config *conf.HttpTransport) transportKey {
	security := conf.TransportSecurity{}
	if config.Security != nil {
		security = *config.Security
	}
	return transportKey{
		MaxConnection:   config.MaxConnection,
		IdleConnTimeout: config.IdleIConnTimeout,
		security:        security,
		useH2:           config.EnableH2,
	}
}
