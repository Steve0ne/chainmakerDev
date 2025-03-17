/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	conf "chainmaker.org/chainmaker-cross/sdk/config"
)

type SDKInterface interface {
	GetConfig() *conf.Config
	GenCrossEvent(params ...*CrossTxBuildCtx) (*CrossEventContext, error)
	SendCrossEvent(event *CrossEventContext, url string, syncResult bool, opts ...EventSendOption) (*eventproto.CrossResponse, error)
	QueryCrossResult(crossID string, url string, opts ...EventSendOption) (*eventproto.CrossResponse, error)
}
