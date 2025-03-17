/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package builder

import (
	"fmt"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	conf "chainmaker.org/chainmaker-cross/sdk/config"
)

type CrossTxBuilder struct {
	//ChainID identification chain
	ChainID string
	//SdkTxBuilder serializes business transaction request
	SdkTxBuilder TxRequestBuilder
	//Config chain config
	Config *conf.CrossChainConf
	//transaction contract parameters builder
	ContractParamBuilder TxContractParamBuilder //事务合约参数构建
}

func (cb *CrossTxBuilder) buildParams(param *CrossTxBuildParam, opts ...ParamsBuildOption) (*requestParam, error) {
	var (
		params = &requestParam{}
		err    error
	)
	params.ExecuteParam, err = cb.ContractParamBuilder.BuildExecuteParam(param, opts...)
	if err != nil {
		return nil, err
	}
	params.CommitParam, err = cb.ContractParamBuilder.BuildCommitParam(param)
	if err != nil {
		return nil, err
	}
	params.RollbackParam, err = cb.ContractParamBuilder.BuildRollbackParam(param)
	if err != nil {
		return nil, err
	}
	return params, nil
}

//Build generates a CrossTx with the parameters
func (cb *CrossTxBuilder) Build(param *CrossTxBuildParam, opts ...CrossBuildOption) (*eventproto.CrossTx, error) {
	proofKey := genProofKey(param)
	options := &crossBuildOptions{
		ProofKey: proofKey,
	}
	for _, o := range opts {
		o(options)
	}
	params, err := cb.buildParams(param, options.ParamOptions...)
	if err != nil {
		return nil, err
	}
	executePayload, err := cb.SdkTxBuilder.Build(&TxRequestBuildParam{
		Contract: &Contract{
			Name:   cb.Config.TransactionContractName,
			Method: cb.Config.TransactionExecuteMethod,
			Params: params.ExecuteParam,
		},
	})
	if err != nil {
		return nil, err
	}
	commitPayload, err := cb.SdkTxBuilder.Build(&TxRequestBuildParam{
		Contract: &Contract{
			Name:   cb.Config.TransactionContractName,
			Method: cb.Config.TransactionCommitMethod,
			Params: params.CommitParam,
		},
	})
	if err != nil {
		return nil, err
	}
	rollbackPayload, err := cb.SdkTxBuilder.Build(&TxRequestBuildParam{
		Contract: &Contract{
			Name:   cb.Config.TransactionContractName,
			Method: cb.Config.TransactionRollbackMethod,
			Params: params.RollbackParam,
		},
	})
	if err != nil {
		return nil, err
	}
	return &eventproto.CrossTx{
		ChainId:         cb.ChainID,
		Index:           param.Index,
		ExecutePayload:  executePayload,
		CommitPayload:   commitPayload,
		RollbackPayload: rollbackPayload,
		ProofKey:        proofKey,
	}, nil
}

func genProofKey(param *CrossTxBuildParam) string {
	return fmt.Sprintf("%s_%04d", param.CrossID, param.Index)
}
