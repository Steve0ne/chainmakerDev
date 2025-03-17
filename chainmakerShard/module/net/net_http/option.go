/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package net_http

type RetryStrategy struct {
	//maximum number of attempts if MaxRetries is -1 no limit
	MaxRetries int
	//the delay time to begin retrying(ms)
	DelayTime int
	//time interval between two transmissions(ms)
	Interval int
}

type connOptions struct {
	RetryStrategy RetryStrategy
}

type ConnOption interface {
	apply(*connOptions)
}

type funcConnOption struct {
	f func(*connOptions)
}

func (fdo *funcConnOption) apply(opts *connOptions) {
	fdo.f(opts)
}

func newFuncConnOption(f func(*connOptions)) *funcConnOption {
	return &funcConnOption{
		f: f,
	}
}

func WithRetryStrategy(strategy RetryStrategy) ConnOption {
	return newFuncConnOption(func(opts *connOptions) {
		opts.RetryStrategy = strategy
	})
}
