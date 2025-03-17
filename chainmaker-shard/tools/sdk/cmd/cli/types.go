/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"sort"

	"chainmaker.org/chainmaker-cross/sdk/builder"
)

var (
	flagNameOfConfigFilepath          = "conf"
	flagNameShortHandOfConfigFilepath = "c"
	flagNameOfUrl                     = "url"
	flagNameShortHandOfUrl            = "u"
	flagNameOfCrossID                 = "crossID"
	flagNameOfParams                  = "params"
)

var (
	ConfigFilepath = "./cross_chain_sdk.yml"
	DefaultURL     = "http://localhost:8080"
)

type crossTxParam struct {
	ChainID        string `mapstructure:"chain_id"`
	Provider       string `mapstructure:"provider"`
	ContractName   string `mapstructure:"contract_name"`
	ExecuteMethod  string `mapstructure:"execute_method"`
	ExecuteParams  Params `mapstructure:"execute_params"`
	RollbackMethod string `mapstructure:"rollback_method"`
	RollbackParams Params `mapstructure:"rollback_params"`
	Index          int32  `mapstructure:"index"`
	ChainType      string `mapstructure:"chain_type"`
}

type Params map[string]string

func (p Params) ToBuilderParams() *builder.Params {
	//if len(p) == 0 {
	//	return nil
	//}
	keys := make([]string, 0, len(p))
	for k, _ := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	kvs := make([]*builder.KV, 0, len(p))
	for _, k := range keys {
		if v, ok := p[k]; ok {
			kvs = append(kvs, builder.NewKV(k, v))
		}
	}
	//fmt.Printf("p is %+v\n", p)
	return builder.NewParams(kvs...)
}

type CrossTxParams struct {
	Params []*crossTxParam `mapstructure:"cross_tx_params"`
}
