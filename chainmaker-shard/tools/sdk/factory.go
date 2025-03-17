/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package sdk

import (
	"errors"
	"fmt"

	"chainmaker.org/chainmaker-cross/sdk/builder/chainmaker"
	"chainmaker.org/chainmaker-cross/sdk/builder/fabric"
	conf "chainmaker.org/chainmaker-cross/sdk/config"

	"chainmaker.org/chainmaker-cross/sdk/builder"
)

// ChainType chain type
type ChainType string

const (
	ChainTypeChainMaker ChainType = "chainmaker"
	ChainTypeFabric 	ChainType = "fabric"
)

var (
	ErrCrossChainConfigNil = errors.New("the cross chain config is nil, please check it")
)

func TxRequestBuilderFactory(typ ChainType, conf *conf.CrossChainConf) (builder.TxRequestBuilder, error) {
	switch typ {
	case ChainTypeChainMaker:
		return chainmaker.NewTxRequestBuilder(conf)
	case ChainTypeFabric:
		return fabric.NewTxRequestBuilder(conf)
	}
	return nil, fmt.Errorf("TxRequestBuilder of the ChainType [%v] Unsupported", typ)
}

func TxParamBuilderFactory(typ ChainType, conf *conf.CrossChainConf) (builder.TxContractParamBuilder, error) {
	switch typ {
	case ChainTypeChainMaker:
		return chainmaker.NewTxContractParamBuilder(conf), nil
	case ChainTypeFabric:
		return fabric.NewTxContractParamBuilder(conf), nil
	}
	return nil, errors.New(fmt.Sprintf("TxParamBuilder of the ChainType [%v] Unsupported", typ))
}

func NewCrossTxBuilder(chainType ChainType, chainConf *conf.CrossChainConf) (*builder.CrossTxBuilder, error) {
	if chainConf == nil {
		return nil, ErrCrossChainConfigNil
	}
	requestBuilder, err := TxRequestBuilderFactory(chainType, chainConf)
	if err != nil {
		return nil, err
	}
	txParamBuilder, err := TxParamBuilderFactory(chainType, chainConf)
	if err != nil {
		return nil, err
	}
	return &builder.CrossTxBuilder{
		ChainID:              chainConf.ChainID,
		SdkTxBuilder:         requestBuilder,
		Config:               chainConf,
		ContractParamBuilder: txParamBuilder,
	}, nil
}
