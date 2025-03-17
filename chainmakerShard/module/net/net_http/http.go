/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package net_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

//type EncodeType int
//
//const (
//	EncodeTypeJSON EncodeType = iota
//)

const (
	DefaultContentType = "application/json"
)

var (
	defaultSendOption = sendOptions{
		Timeout: 10 * time.Second,
		//Encode:     EncodeTypeJSON,
	}
)

type sendOptions struct {
	Timeout   time.Duration
	Transport http.RoundTripper
	//Encode     EncodeType
}

type SendOption func(*sendOptions)

func WithTimeout(timeout time.Duration) SendOption {
	return func(opts *sendOptions) {
		opts.Timeout = timeout
	}
}

func WithRoundTripper(rt http.RoundTripper) SendOption {
	return func(opts *sendOptions) {
		opts.Transport = rt
	}
}

type HttpRequest struct {
	URL     string
	Method  string
	Content interface{}
	options sendOptions
}

type HttpResponse struct {
	*http.Response
	content []byte
}

func (r *HttpRequest) Send(opts ...SendOption) (*HttpResponse, error) {
	for _, o := range opts {
		o(&r.options)
	}
	client := http.Client{
		Timeout:   r.options.Timeout,
		Transport: r.options.Transport,
	}
	var body io.Reader
	if r.Content != nil {
		jsonStr, err := json.Marshal(r.Content)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonStr)
	}

	req, err := http.NewRequest(r.Method, r.URL, body)
	req.Header.Set("Content-Type", DefaultContentType)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return &HttpResponse{
		Response: resp,
	}, nil
}

func (resp *HttpResponse) UnmarshalToObj(out interface{}) (err error) {
	if resp.content == nil {
		if resp.Body == nil {
			return errors.New("response has no body")
		}
		defer resp.Body.Close()
		resp.content, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
	}
	return json.Unmarshal(resp.content, out)
}

func (resp *HttpResponse) RawResponse() *http.Response {
	return resp.Response
}

func NewHttpRequest(url string, method string, data interface{}) *HttpRequest {
	return &HttpRequest{
		URL:     url,
		Method:  method,
		Content: data,
		options: defaultSendOption,
	}
}
