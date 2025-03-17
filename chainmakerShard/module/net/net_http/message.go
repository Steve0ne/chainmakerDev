/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package net_http

import (
	"errors"
	"time"

	"chainmaker.org/chainmaker-cross/event"
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	"github.com/google/uuid"
)

type Request struct {
	*Message
	URL       string `json:"-"`
	ID        string `json:"id,omitempty"`        // 消息的UUID
	Timestamp int64  `json:"timestamp,omitempty"` // 时间戳
}

type Message struct {
	Type       eventproto.EventType `json:"event_type"`
	EncodeType event.MarshalType    `json:"encode_type"`
	Payload    []byte               `json:"payload"`
}

func (m *Message) GetNodeID() string {
	return ""
}

func (m *Message) GetPayload() []byte {
	return m.Payload
}

func (m *Message) Marshal(in event.Event) ([]byte, error) {
	return eventCoders.GetCoder(m.EncodeType).Marshal(in)
}

func (m *Message) Unmarshal(out event.Event) error {
	if out.GetType() != m.Type {
		return errors.New("Mismatched data types")
	}
	return eventCoders.GetCoder(m.EncodeType).Unmarshal(m.Payload, out)
}

func NewRequest(eve event.Event, url string, encodeType event.MarshalType) (req *Request, err error) {
	req = &Request{
		URL:       url,
		Timestamp: time.Now().Unix(),
		ID:        uuid.New().String(),
	}
	if req.Message, err = NewMessage(eve, encodeType); err != nil {
		return nil, err
	}
	return
}

func NewMessage(eve event.Event, encodeType event.MarshalType) (msg *Message, err error) {
	msg = &Message{
		Type:       eve.GetType(),
		EncodeType: encodeType,
	}
	msg.Payload, err = msg.Marshal(eve)
	return
}

type Data = *Message

type Response struct {
	Code    int    `json:"code,omitempty"` //usually 0 means successful processing
	Message string `json:"message,omitempty"`
	Data    `json:"data"`
}
