/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package builder

import (
	"strconv"
)

//parameters for build cross-chain transactions
type CrossTxBuildParam struct {
	//CrossID representing a cross-chain transaction
	CrossID string
	//transaction index,don't repeat
	Index int32
	//business contract information to be executed
	ExecuteBusinessContract *Contract
	//business contract information to be rolled back
	RollbackBusinessContract *Contract
}

func (cbp *CrossTxBuildParam) SetCrossID(crossID string) {
	if cbp != nil {
		cbp.CrossID = crossID
	}
}

func (c *CrossTxBuildParam) Fix() {
	if c.ExecuteBusinessContract.Name == "" && c.RollbackBusinessContract.Name != "" {
		c.ExecuteBusinessContract.Name = c.RollbackBusinessContract.Name
		return
	}
	if c.RollbackBusinessContract.Name == "" && c.ExecuteBusinessContract.Name != "" {
		c.RollbackBusinessContract.Name = c.ExecuteBusinessContract.Name
		return
	}
}

type requestParam struct {
	ExecuteParam  *Params
	CommitParam   *Params
	RollbackParam *Params
}

//Contract represents a contract information
type Contract struct {
	//Contract Name
	Name string
	//contract method to be called
	Method string
	//call parameters of contract methods
	Params *Params
}

func NewCrossTxBuildParam(crossID string, index int32, execute *Contract, rollback *Contract) *CrossTxBuildParam {
	return &CrossTxBuildParam{
		CrossID:                  crossID,
		Index:                    index,
		ExecuteBusinessContract:  execute,
		RollbackBusinessContract: rollback,
	}
}

func NewContract(name, method string, params *Params) *Contract {
	return &Contract{
		Name:   name,
		Method: method,
		Params: params,
	}
}

type Params struct {
	kvs []*KV
}
type KV struct {
	key   string
	value string
}

func (p *Params) Len() int {
	if p == nil {
		return 0
	}
	return len(p.kvs)
}

func (p *Params) GetKVMap() map[string]string {
	if p == nil {
		return nil
	}
	m := make(map[string]string, len(p.kvs))
	for _, kv := range p.kvs {
		m[kv.key] = kv.value
	}
	return m
}

func (p *Params) GetKVBytesMap() map[string][]byte {
	if p == nil {
		return nil
	}
	m := make(map[string][]byte, len(p.kvs))
	for _, kv := range p.kvs {
		m[kv.key] = []byte(kv.value)
	}
	return m
}

func (p *Params) Values() []string {
	if p == nil {
		return nil
	}
	s := make([]string, len(p.kvs))
	for i, kv := range p.kvs {
		s[i] = kv.value
	}
	return s
}

func NewParams(kvs ...*KV) *Params {
	p := &Params{
		kvs: kvs,
	}
	return p
}

//NewParamsNoKeys creates a new Params with all values
func NewParamsNoKeys(values ...string) *Params {
	kvs := make([]*KV, len(values))
	for i, v := range values {
		kvs[i] = &KV{key: strconv.Itoa(i), value: v}
	}
	return NewParams(kvs...)
}

//NewParamsNoKeys creates a new Params with a map[string]string
func NewParamsWithMap(m map[string]string) *Params {
	kvs := make([]*KV, len(m))
	i := 0
	for k, v := range m {
		kvs[i] = &KV{key: k, value: v}
		i++
	}
	return NewParams(kvs...)
}

func NewKV(k string, v string) *KV {
	return &KV{key: k, value: v}
}
