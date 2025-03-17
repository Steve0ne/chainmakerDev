/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package fabric

import (
	"encoding/json"

	"chainmaker.org/chainmaker-cross/sdk/builder"
	conf "chainmaker.org/chainmaker-cross/sdk/config"
)

type txContractParamBuilder struct {
	Config *conf.CrossChainConf
}

// necessary params to call a contract
type CallContractParams struct {
	ContractName string   `json:"contract_name"`
	Method       string   `json:"method"`
	Params       []string `json:"params,omitempty" metadata:",optional"`
}

// ExecuteCallRequest chainmaker execute payload
type ExecuteCallRequest struct {
	crossID            string
	executeCallParams  CallContractParams
	rollbackCallParams CallContractParams
}

func (pb *txContractParamBuilder) BuildExecuteParam(in *builder.CrossTxBuildParam, opts ...builder.ParamsBuildOption) (*builder.Params, error) {
	in = pb.refactorParam(in, opts...)
	eParams := ContractToCallContractParams(in.ExecuteBusinessContract)
	eParamsBz, err := json.Marshal(eParams)
	if err != nil {
		return nil, err
	}
	rParams := ContractToCallContractParams(in.RollbackBusinessContract)
	rParamsBz, err := json.Marshal(rParams)
	if err != nil {
		return nil, err
	}

	return builder.NewParamsNoKeys(in.CrossID, string(eParamsBz), string(rParamsBz)), nil
}

func (pb *txContractParamBuilder) BuildCommitParam(in *builder.CrossTxBuildParam) (*builder.Params, error) {
	return builder.NewParamsNoKeys(in.CrossID), nil
}

func (pb *txContractParamBuilder) BuildRollbackParam(in *builder.CrossTxBuildParam) (*builder.Params, error) {
	return builder.NewParamsNoKeys(in.CrossID), nil
}

func (pb *txContractParamBuilder) refactorParam(in *builder.CrossTxBuildParam, opts ...builder.ParamsBuildOption) *builder.CrossTxBuildParam {
	if len(opts) == 0 {
		return in
	}
	options := builder.NewParamsBuildOptions(opts...)
	out := *in
	if options.UseProofKey { //如果业务合约需要使用proofkey，那么在fabric中，强制将proofkey作为第一个参数
		m := make([]string, in.ExecuteBusinessContract.Params.Len()+1)
		m[0] = options.ProofKey
		copy(m[1:], in.ExecuteBusinessContract.Params.Values())
		out.ExecuteBusinessContract.Params = builder.NewParamsNoKeys(m...)
	}
	return &out
}

func NewTxContractParamBuilder(conf *conf.CrossChainConf) *txContractParamBuilder {
	return &txContractParamBuilder{
		Config: conf,
	}
}

func ContractToCallContractParams(c *builder.Contract) *CallContractParams {
	return &CallContractParams{
		ContractName: c.Name,
		Method:       c.Method,
		Params:       c.Params.Values(),
	}
}
