/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"context"
	"fmt"
	"time"

	"chainmaker.org/chainmaker-cross/net/net_http"
)

const (
	urlCrossEvent       = "/cross?method=InvokeCrossEvent"
	urlQueryEventResult = "/cross?method=GetCrossEvent"
)

var (
	ErrCrossTxOverrun            = fmt.Errorf("number of cross transaction gather than %d", CrossTxsLimit)
	ErrCrossTxMismatch           = fmt.Errorf("number of cross transaction is not %d", CrossTxsLimit)
	ErrCrossEventCrossIDNotMatch = fmt.Errorf("this cross event sent CrossID does not match the received CrossID ")

	defaultEventSendOptions = eventSendOptions{
		Timeout:      10 * time.Second,
		SyncStrategy: defaultSyncStrategy,
	}
	defaultSyncStrategy = SyncStrategy{
		MaxRetries:    5,
		DelaySyncTime: 5 * time.Second,
		Interval:      100 * time.Millisecond,
	}
)

type eventSendOptions struct {
	//the context for http send
	Ctx context.Context
	//the http request timeout
	Timeout time.Duration
	//strategy for obtaining cross-chain event result synchronously
	SyncStrategy SyncStrategy

	HttpSendOption []net_http.SendOption
}

type SyncStrategy struct {
	//maximum number of attempts if MaxRetries is zero no limit until get the corresponding result
	MaxRetries int
	//the delay time to sync the result for the first time
	DelaySyncTime time.Duration
	//time interval between sync result calls
	Interval time.Duration
}

type EventSendOption func(*eventSendOptions)

func WithTimeoutOpt(timeout time.Duration) EventSendOption {
	return func(options *eventSendOptions) {
		options.Timeout = timeout
		options.HttpSendOption = append(options.HttpSendOption, net_http.WithTimeout(timeout))
	}
}

func WithContextOpt(ctx context.Context) EventSendOption {
	return func(options *eventSendOptions) {
		options.Ctx = ctx
	}
}

func WithSyncStrategyOpt(s SyncStrategy) EventSendOption {
	return func(options *eventSendOptions) {
		options.SyncStrategy = s
	}
}

func (opt *eventSendOptions) HttpSendOptions() []net_http.SendOption {
	return opt.HttpSendOption
}

func NewSendOptions(opts ...EventSendOption) *eventSendOptions {
	sendOpts := defaultEventSendOptions
	sendOpts.Ctx = context.Background()
	for _, o := range opts {
		o(&sendOpts)
	}
	return &sendOpts
}

//SyncBackOff retry the synchronization request according to the parameters s
//invoke is the request logic, if invoke returns true，it means that the correct data is requested, otherwise it will retry
func SyncBackOff(ctx context.Context, invoke func() bool, s SyncStrategy) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	attempts := 0
	if s.DelaySyncTime > 0 {
		time.Sleep(s.DelaySyncTime)
	}

	for {
		if s.MaxRetries > 0 && attempts >= s.MaxRetries {
			return nil
		}
		if ok := invoke(); ok {
			return nil
		}
		select {
		case <-time.After(s.Interval):
		case <-ctx.Done():
			return ctx.Err()
		}
		attempts++
	}
	return nil
}
