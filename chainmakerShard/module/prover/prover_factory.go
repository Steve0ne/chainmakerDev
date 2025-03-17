/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package prover

import (
	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/prover/impl"
)

type Provider string

const (
	TrustProvider Provider = "trust"
	SpvProvider   Provider = "spv"
)

// InitProvers init all the provers
func InitProvers() *ProverDispatcher {
	for _, proverCfg := range conf.Config.ProverConfigs {
		var prov Prover
		switch Provider(proverCfg.Provider) {
		case TrustProvider:
			prov = impl.NewTrustProver(proverCfg.GetChainIDs())
		case SpvProvider:
			prov = impl.NewSpvProver(conf.FinalCfgPath(proverCfg.ConfigPath), proverCfg.GetChainIDs()) // TODO Unit Test
		default:
			prov = impl.NewTrustProver(proverCfg.GetChainIDs())
		}
		dispatcher.Register(prov)
	}
	return dispatcher
}
