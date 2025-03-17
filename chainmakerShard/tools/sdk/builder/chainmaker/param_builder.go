/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package chainmaker

import (
	"chainmaker.org/chainmaker-cross/sdk/builder"
	conf "chainmaker.org/chainmaker-cross/sdk/config"
	"chainmaker.org/chainmaker/common/serialize"
)

type txContractParamBuilder struct {
	Config *conf.CrossChainConf
}

func (pb *txContractParamBuilder) BuildExecuteParam(in *builder.CrossTxBuildParam, opts ...builder.ParamsBuildOption) (*builder.Params, error) {
	in = pb.refactorParam(in, opts...)
	eParams := map[string]string{
		pb.Config.BusinessCrossIDKey:      in.CrossID,
		pb.Config.BusinessContractNameKey: in.ExecuteBusinessContract.Name,
		pb.Config.BusinessMethodKey:       in.ExecuteBusinessContract.Method,
		pb.Config.BusinessParamsKey:       string(serialize.EasyMarshal(serialize.ParamsMapToEasyCodecItem(in.ExecuteBusinessContract.Params.GetKVBytesMap()))),
	}
	rParams := map[string]string{
		pb.Config.BusinessCrossIDKey:      in.CrossID,
		pb.Config.BusinessContractNameKey: in.RollbackBusinessContract.Name,
		pb.Config.BusinessMethodKey:       in.RollbackBusinessContract.Method,
		pb.Config.BusinessParamsKey:       string(serialize.EasyMarshal(serialize.ParamsMapToEasyCodecItem(in.RollbackBusinessContract.Params.GetKVBytesMap()))),
	}

	m := map[string]string{
		pb.Config.BusinessCrossIDKey:         in.CrossID,
		pb.Config.TransactionExecuteDataKey:  string(serialize.EasyMarshal(serialize.ParamsMapToEasyCodecItem(stringMap2BytesMap(eParams)))),
		pb.Config.TransactionRollbackDataKey: string(serialize.EasyMarshal(serialize.ParamsMapToEasyCodecItem(stringMap2BytesMap(rParams)))),
	}
	return builder.NewParamsWithMap(m), nil
}

func (pb *txContractParamBuilder) BuildCommitParam(in *builder.CrossTxBuildParam) (*builder.Params, error) {
	return builder.NewParamsWithMap(map[string]string{
		pb.Config.BusinessCrossIDKey: in.CrossID,
	}), nil
}

func (pb *txContractParamBuilder) BuildRollbackParam(in *builder.CrossTxBuildParam) (*builder.Params, error) {
	return builder.NewParamsWithMap(map[string]string{
		pb.Config.BusinessCrossIDKey: in.CrossID,
	}), nil
}

func (pb *txContractParamBuilder) refactorParam(in *builder.CrossTxBuildParam, opts ...builder.ParamsBuildOption) *builder.CrossTxBuildParam {
	if len(opts) == 0 {
		return in
	}
	options := builder.NewParamsBuildOptions(opts...)
	out := *in
	if options.UseProofKey {
		m := in.ExecuteBusinessContract.Params.GetKVMap()
		m[pb.Config.BusinessProofKey] = options.ProofKey
		out.ExecuteBusinessContract.Params = builder.NewParamsWithMap(m)
	}
	return &out
}

func stringMap2BytesMap(m map[string]string) map[string][]byte {
	bm := make(map[string][]byte, len(m))
	for k, v := range m {
		bm[k] = []byte(v)
	}
	return bm
}

func NewTxContractParamBuilder(conf *conf.CrossChainConf) *txContractParamBuilder {
	return &txContractParamBuilder{
		Config: conf,
	}
}
