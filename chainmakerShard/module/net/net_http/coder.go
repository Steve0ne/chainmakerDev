/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package net_http

import (
	"encoding/json"
	"errors"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/event/coder"
	"chainmaker.org/chainmaker-cross/utils"
)

type EventCoder interface {
	Marshal(in event.Event) ([]byte, error)
	Unmarshal(data []byte, out event.Event) error
}

type ErrorCoder struct {
}

func (e ErrorCoder) Marshal(event.Event) ([]byte, error) {
	return nil, errors.New("no supported this type of coder")
}

func (e ErrorCoder) Unmarshal([]byte, event.Event) error {
	return errors.New("no supported this type of coder")
}

type JsonCoder struct {
}

func (j JsonCoder) Marshal(in event.Event) ([]byte, error) {
	return json.Marshal(in)
}

func (j JsonCoder) Unmarshal(data []byte, out event.Event) error {
	return json.Unmarshal(data, out)
}

type BinaryCoder struct {
}

func (b BinaryCoder) Marshal(in event.Event) ([]byte, error) {
	binary, err := coder.JsonBinaryMarshal(in.GetType(), in)
	if err != nil {
		return nil, err
	}
	return []byte(utils.Base64EncodeToString(binary)), nil
}

func (b BinaryCoder) Unmarshal(data []byte, out event.Event) error {
	resp, err := utils.Base64DecodeToBytes(string(data))
	if err != nil {
		return err
	}
	return coder.JsonBinaryUnmarshal(resp, byte(out.GetType()), out)
}

//or use slice?
//type EventCoders map[event.MarshalType]EventCoder

type EventCoders []EventCoder

var eventCoders EventCoders

func init() {
	eventCoders = EventCoders{
		event.JsonMarshalType:   JsonCoder{},
		event.BinaryMarshalType: BinaryCoder{},
	}
}

func (ecs EventCoders) GetCoder(typ event.MarshalType) EventCoder {
	if ecr := ecs[typ]; ecr != nil {
		return ecr
	}
	return ErrorCoder{}
}
