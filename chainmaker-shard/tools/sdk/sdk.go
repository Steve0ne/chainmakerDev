/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"fmt"
	"strings"

	"chainmaker.org/chainmaker-cross/net/net_http"
	"github.com/pkg/errors"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	"chainmaker.org/chainmaker-cross/sdk/builder"
	conf "chainmaker.org/chainmaker-cross/sdk/config"
)

//CrossSDK define the Cross SDK struct
type CrossSDK struct {
	builders map[string]*builder.CrossTxBuilder
	opts     crossSdkOptions
}

//CrossTxBuildCtx uses to build cross tx
type CrossTxBuildCtx struct {
	chainID    string
	buildParam *builder.CrossTxBuildParam
	buildOpts  []builder.CrossBuildOption
}

type crossSdkOptions struct {
	configFile       string
	configDecorators []conf.CfgDecorator
	config           *conf.Config
}

//SDKOption a interface applied to crossSdkOptions
type SDKOption interface {
	apply(*crossSdkOptions)
}

type funcSDKOption struct {
	do func(*crossSdkOptions)
}

func (f *funcSDKOption) apply(opt *crossSdkOptions) {
	f.do(opt)
}

func newFuncSDKOption(f func(options *crossSdkOptions)) *funcSDKOption {
	return &funcSDKOption{
		do: f,
	}
}

//WithConfigDecorators set opts config decorators
func WithConfigDecorators(decors ...conf.CfgDecorator) SDKOption {
	return newFuncSDKOption(func(options *crossSdkOptions) {
		options.configDecorators = append(options.configDecorators, decors...)
	})
}

//WithConfigFile set opts config file
func WithConfigFile(file string) SDKOption {
	return newFuncSDKOption(func(options *crossSdkOptions) {
		options.configFile = file
	})
}

//WithConfig set opts config
func WithConfig(cfg *conf.Config) SDKOption {
	return newFuncSDKOption(func(options *crossSdkOptions) {
		options.config = cfg
	})
}

func NewCrossTxBuildCtx(chainID string, index int32, execute *builder.Contract, rollback *builder.Contract, opts ...builder.CrossBuildOption) *CrossTxBuildCtx {
	return &CrossTxBuildCtx{
		chainID:    chainID,
		buildParam: builder.NewCrossTxBuildParam("", index, execute, rollback),
		buildOpts:  opts,
	}
}

//NewCrossSDK create a Cross SDK instance
func NewCrossSDK(opt ...SDKOption) (*CrossSDK, error) {
	opts := crossSdkOptions{}
	for _, o := range opt {
		o.apply(&opts)
	}
	var (
		cfg = opts.config
		err error
	)
	if cfg == nil {
		cfg, err = conf.LoadConfig(opts.configFile)
		if err != nil {
			return nil, err
		}
	}

	cfg, err = cfg.Overload(opts.configDecorators...)
	if err != nil {
		return nil, errors.WithMessage(err, "config overload fail:")
	}
	opts.config = cfg
	sdk := &CrossSDK{
		builders: make(map[string]*builder.CrossTxBuilder, len(opts.config.ConfigLists)),
		opts:     opts,
	}
	err = sdk.init()
	if err != nil {
		return nil, err
	}
	return sdk, nil
}

func (s *CrossSDK) init() error {
	chainList := s.opts.config.ConfigLists
	if s.builders == nil {
		s.builders = make(map[string]*builder.CrossTxBuilder, len(chainList))
	}

	for _, c := range chainList {
		b, err := NewCrossTxBuilder(ChainType(c.ChainType), c)
		if err != nil {
			return err
		}
		s.builders[c.ChainID] = b
	}
	return nil
}

func (s *CrossSDK) GetConfig() *conf.Config {
	return s.opts.config
}

func (s *CrossSDK) getCrossTxBuilder(chainID string) (b *builder.CrossTxBuilder, ok bool) {
	b, ok = s.builders[chainID]
	if ok && b != nil {
		return
	}
	return nil, false
}

func (s *CrossSDK) GenCrossEvent(params ...*CrossTxBuildCtx) (*CrossEventContext, error) {
	crossEvent := NewCrossEventCtx()
	crossTxs := make([]*eventproto.CrossTx, 0, len(params))
	for _, param := range params {
		b, ok := s.getCrossTxBuilder(param.chainID)
		if ok {
			if param.buildParam == nil {
				return nil, errors.New("CrossTxParam is invalid")
			}
			param.buildParam.SetCrossID(crossEvent.GetCrossID())
			crossTx, err := b.Build(param.buildParam, param.buildOpts...)
			if err != nil {
				return nil, err
			}
			crossTxs = append(crossTxs, crossTx)
		} else {
			return nil, fmt.Errorf("chainID [%s] builder is not exist", param.chainID)
		}
	}
	err := crossEvent.BuildEvent(crossTxs...)
	if err != nil {
		return nil, err
	}
	return crossEvent, err
}

func (s *CrossSDK) SendCrossEvent(event *CrossEventContext, url string, syncResult bool, opts ...EventSendOption) (*eventproto.CrossResponse, error) {
	if event == nil {
		return nil, errors.New("crossEvent to be sent is invalid")
	}
	eventSendOpts, err := s.getEventSendOptions(url, opts...)
	if err != nil {
		return nil, err
	}
	return event.Send2(url, syncResult, *eventSendOpts)
}

func (s *CrossSDK) getEventSendOptions(url string, opts ...EventSendOption) (*eventSendOptions, error) {
	eventSendOpts := NewSendOptions(opts...)
	if strings.HasPrefix(url, "https://") && s.opts.config.Http != nil {
		tr, err := net_http.GetTransport(s.opts.config.Http)
		if err != nil {
			return nil, err
		}
		eventSendOpts.HttpSendOption = append(eventSendOpts.HttpSendOption, net_http.WithRoundTripper(tr))
	}
	return eventSendOpts, nil
}

func (s *CrossSDK) QueryCrossResult(crossID string, url string, opts ...EventSendOption) (*eventproto.CrossResponse, error) {
	searchEvent := NewCrossSearchEvent(crossID)
	eventSendOpts, err := s.getEventSendOptions(url, opts...)
	if err != nil {
		return nil, err
	}
	return searchEvent.QueryWithOptions(url, *eventSendOpts)
}
