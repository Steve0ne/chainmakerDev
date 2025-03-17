/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"chainmaker.org/chainmaker-cross/conf"

	"chainmaker.org/chainmaker-cross/net/net_http"

	"chainmaker.org/chainmaker-cross/event"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
)

const (
	CrossTxsLimit = 2
)

//CrossEventContext represents a context of CrossEvent
type CrossEventContext struct {
	event    *eventproto.CrossEvent
	sendTime int64
	result   *eventproto.CrossResponse
	config   *conf.HttpTransport
}

//CrossEventSendResp a send CrossEvent response
type CrossEventSendResp struct {
	CrossID string
}

type crossSearchEvent struct {
	eventproto.CrossSearchEvent
}

//NewCrossEventCtx create a CrossEventContext
func NewCrossEventCtx() *CrossEventContext {
	return &CrossEventContext{
		event: event.NewEmptyCrossEvent(),
	}
}

//func (cc *CrossEventContext) SetHttpConfig(config *conf.HttpTransport) {
//	cc.config = config
//}

//GetCrossID get the cross ID of CrossEvent
func (cc *CrossEventContext) GetCrossID() string {
	//if cc == nil || cc.event == nil {
	//	return ""
	//}
	return cc.event.GetCrossID()
}

//BuildEvent construct the CrossEvent through parameters txs
//txs is a variable parameter, note: the current limit for cross-chain transactions is two chains
func (cc *CrossEventContext) BuildEvent(txs ...*eventproto.CrossTx) error {
	if cc.event.TxEvents.Len()+len(txs) > CrossTxsLimit {
		return ErrCrossTxOverrun
	}
	cc.event.TxEvents.Append(txs...)
	return nil
}

func (cc *CrossEventContext) sendCheck() error {
	if cc.event.TxEvents.Len() != CrossTxsLimit {
		return ErrCrossTxMismatch
	}
	return nil
}

func (cc *CrossEventContext) getEventSendOptions(url string, opts ...EventSendOption) (*eventSendOptions, error) {
	eventSendOpts := NewSendOptions(opts...)
	if strings.HasPrefix(url, "https://") && cc.config != nil {
		tr, err := net_http.GetTransport(cc.config)
		if err != nil {
			return nil, err
		}
		eventSendOpts.HttpSendOption = append(eventSendOpts.HttpSendOption, net_http.WithRoundTripper(tr))
	}
	return eventSendOpts, nil
}

/*
send a cross-chain transaction event by http
if syncResult is true, this method will return the processing result synchronously
otherwise you will get the result by the GetResult method
opts can be used to control this sending behavior, such as http request timeout,
*/
func (cc *CrossEventContext) Send(url string, syncResult bool, opts ...EventSendOption) (*eventproto.CrossResponse, error) {
	//if err := cc.sendCheck(); err != nil {
	//	return nil, err
	//}
	//cc.sendTime = time.Now().Unix()
	eventSendOpts, err := cc.getEventSendOptions(url, opts...)
	if err != nil {
		return nil, err
	}
	return cc.Send2(url, syncResult, *eventSendOpts)
}

func (cc *CrossEventContext) Send2(url string, syncResult bool, eventSendOpts eventSendOptions) (*eventproto.CrossResponse, error) {
	if err := cc.sendCheck(); err != nil {
		return nil, err
	}
	cc.sendTime = time.Now().Unix()
	req := net_http.NewHttpRequest(url+urlCrossEvent, http.MethodPost, cc.event)
	httpResp, err := req.Send(eventSendOpts.HttpSendOptions()...)
	if err != nil {
		return nil, err
	}
	sendResp := &CrossEventSendResp{}
	if err := httpResp.UnmarshalToObj(sendResp); err != nil {
		return nil, err
	}
	if sendResp.CrossID != cc.event.CrossId {
		return nil, ErrCrossEventCrossIDNotMatch
	}

	if syncResult {
		//ctx, _ := context.WithTimeout(eventSendOpts.Ctx, 10*time.Second)
		resp, err := cc.syncEventResult(eventSendOpts.Ctx, url, eventSendOpts)
		if err != nil {
			err = fmt.Errorf("Sync Cross Event Result Error: [%v]", err)
			return nil, err
		}
		cc.setResult(resp)
		return resp, nil
	}
	return &eventproto.CrossResponse{
		CrossId: sendResp.CrossID,
		Code:    event.UnknownResp,
	}, nil
}

func (cc *CrossEventContext) syncEventResult(ctx context.Context, url string, options eventSendOptions) (resp *eventproto.CrossResponse, err error) {
	searchEvent := NewCrossSearchEvent(cc.event.CrossId)
	netSendOpts := options.HttpSendOptions()
	syncErr := SyncBackOff(ctx, func() bool {
		resp, err = searchEvent.queryWithNetSendOption(url, netSendOpts...)
		if err != nil {
			return true
		}
		if resp.Code == event.UnknownResp {
			return false
		}
		return true
	}, options.SyncStrategy)
	if syncErr != nil {
		return nil, syncErr
	}
	return
}

func (cc *CrossEventContext) setResult(resp *eventproto.CrossResponse) {
	if resp.Code == event.SuccessResp || resp.Code == event.FailureResp {
		cc.result = resp
	}
}

//GetResult obtain the result of the specified crossID event
func (cc *CrossEventContext) GetResult(url string, opts ...EventSendOption) (*eventproto.CrossResponse, error) {
	if cc.result != nil {
		return cc.result, nil
	}
	resp, err := NewCrossSearchEvent(cc.event.CrossId).Query(url, opts...)
	if err != nil {
		return nil, err
	}
	cc.setResult(resp)
	return resp, nil
}

//Query the result of the specified crossID event
func (q *crossSearchEvent) Query(url string, opts ...EventSendOption) (*eventproto.CrossResponse, error) {
	eventSendOpts := NewSendOptions(opts...)
	return q.queryWithNetSendOption(url, eventSendOpts.HttpSendOptions()...)
}

func (q *crossSearchEvent) QueryWithOptions(url string, eventSendOpts eventSendOptions) (*eventproto.CrossResponse, error) {
	return q.queryWithNetSendOption(url, eventSendOpts.HttpSendOptions()...)
}

func (q *crossSearchEvent) queryWithNetSendOption(url string, opt ...net_http.SendOption) (*eventproto.CrossResponse, error) {
	req := net_http.NewHttpRequest(url+urlQueryEventResult, http.MethodPost, q.CrossSearchEvent)
	resp, err := req.Send(opt...)
	if err != nil {
		return nil, err
	}
	eventResult := &eventproto.CrossResponse{}
	err = resp.UnmarshalToObj(eventResult)
	if err != nil {
		return nil, err
	}
	return eventResult, nil
}

//NewCrossSearchEvent create CrossSearchEvent with the specified crossID
func NewCrossSearchEvent(crossID string) *crossSearchEvent {
	return &crossSearchEvent{CrossSearchEvent: eventproto.CrossSearchEvent{
		CrossId: crossID,
	}}
}
