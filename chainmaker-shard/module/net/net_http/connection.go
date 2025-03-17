/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package net_http

import (
	"context"
	"io"
	"io/ioutil"
	"strings"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/net"
	"go.uber.org/zap"

	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	DefaultTimeout  int = 10000
	ReadDataBufSize     = 1024
)

type Connection struct {
	log            *zap.SugaredLogger
	address        string
	requestTimeout time.Duration
	transport      http.RoundTripper
	ch             chan net.Message
	opts           connOptions
}

func (c *Connection) GetProvider() net.ConnectionProvider {
	return net.HttpConnection
}

func (c *Connection) WriteData(msg net.Message) error {
	req, ok := msg.(*Request)
	if !ok {
		return errors.New("msg is not http formatted message")
	}
	if len(req.URL) == 0 {
		c.log.Infof("the URL of the http request [%v] is invalid", req)
		return errors.New("invalid url")
	}
	if !strings.HasPrefix(req.URL, "/") {
		req.URL = "/" + req.URL
	}

	resp, err := c.send(req)
	if err != nil {
		c.log.Errorf("URL [%s] send request error [%s]", c.address+req.URL, err.Error())
		return err
	}
	return c.deliverMessage(resp.RawResponse())
}

func (c *Connection) send(req *Request) (*HttpResponse, error) {
	var (
		resp *HttpResponse
		err  error
	)
	c.log.Infof("send request to URL [%s]", c.address+req.URL)
	httpReq := NewHttpRequest(c.address+req.URL, http.MethodPost, req)
	if c.opts.RetryStrategy.MaxRetries != 0 {
		backOff(context.Background(), func() bool {
			resp, err = httpReq.Send(WithTimeout(c.requestTimeout), WithRoundTripper(c.transport))
			if err != nil {
				c.log.Infof("retry send request to URL [%s], error: [%v]", c.address+req.URL, err)
				return false
			}
			return true
		}, c.opts.RetryStrategy)
	} else {
		resp, err = httpReq.Send(WithTimeout(c.requestTimeout), WithRoundTripper(c.transport))
	}
	return resp, err
}

func backOff(ctx context.Context, invoke func() bool, s RetryStrategy) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if s.DelayTime > 0 {
		time.Sleep(time.Duration(s.DelayTime) * time.Millisecond)
	}
	attempts := 0
	for {
		if s.MaxRetries > 0 && attempts >= s.MaxRetries {
			return nil
		}
		if ok := invoke(); ok {
			return nil
		}
		select {
		case <-time.After(time.Duration(s.Interval) * time.Millisecond):
		case <-ctx.Done():
			return ctx.Err()
		}
		attempts++
	}
	return nil
}

func (c *Connection) deliverMessage(resp *http.Response) error {
	if resp.StatusCode >= http.StatusInternalServerError {
		c.log.Infof("the response from address [%s] is error [%+v]", c.address, resp.Status)
		return errors.New(resp.Status)
	}
	if resp.Body != nil {
		rsp := &Response{}
		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()
		if err := decoder.Decode(rsp); err != nil && err != io.EOF {
			b, _ := ioutil.ReadAll(decoder.Buffered())
			c.log.Infof("decoder the response error [%v], raw content is [%s]", err, string(b))
			return err
		}
		if rsp.Code != 0 {
			c.log.Infof("the response from address [%s] is fail [%+v]", c.address, rsp)
			return nil
		}
		select {
		case c.ch <- rsp.Data:
		default:
			go func() {
				c.ch <- rsp.Data
			}()
		}
	}
	return nil
}

func (c *Connection) ReadData() (chan net.Message, error) {
	return c.ch, nil
}

func (c *Connection) PeerID() string {
	return ""
}

func (c *Connection) Close() error {
	return nil
}

func NewConnection(config *conf.HttpRouterConfig, opt ...ConnOption) (*Connection, error) {
	opts := connOptions{}
	for _, o := range opt {
		o.apply(&opts)
	}
	addr := config.Address
	if !strings.HasPrefix(addr, "http://") && !strings.HasPrefix(addr, "https://") {
		if config.EnableTLS {
			addr = "https://" + addr
		} else {
			addr = "http://" + addr
		}
	}

	rto := config.RequestTimeout
	if rto <= 0 {
		rto = DefaultTimeout
	}
	tr, err := GetTransport(&config.HttpTransport)
	if err != nil {
		return nil, err
	}
	return &Connection{
		address:        addr,
		transport:      tr,
		requestTimeout: time.Duration(rto) * time.Millisecond,
		ch:             make(chan net.Message, ReadDataBufSize),
		log:            logger.GetLogger(logger.ModuleNet),
		opts:           opts,
	}, nil
}
