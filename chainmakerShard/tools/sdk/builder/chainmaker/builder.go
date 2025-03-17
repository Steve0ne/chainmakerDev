/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package chainmaker

import (
	"io/ioutil"

	conf "chainmaker.org/chainmaker-cross/sdk/config"

	"chainmaker.org/chainmaker-cross/sdk/builder"
	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmsdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/golang/protobuf/proto"
)

type txRequestBuilder struct {
	chainMakerSDK cmsdk.SDKInterface
}

func (cb *txRequestBuilder) Build(in *builder.TxRequestBuildParam) ([]byte, error) {
	req, err := cb.chainMakerSDK.GetTxRequest(in.Contract.Name, in.Contract.Method, in.TxID, paramsToKeyValuePair(in.Contract.Params))
	if err != nil {
		return nil, err
	}
	return proto.Marshal(req)
}

func paramsToKeyValuePair(p *builder.Params) []*common.KeyValuePair {
	m := p.GetKVBytesMap()
	kvp := make([]*common.KeyValuePair, len(m))
	i := 0
	for k, v := range m {
		kvp[i] = &common.KeyValuePair{
			k, v,
		}
		i++
	}
	return kvp
}

func NewTxRequestBuilder(conf *conf.CrossChainConf) (*txRequestBuilder, error) {
	SignKeyBytes, err := loadBytes(conf.SignKeyPath)
	if err != nil {
		return nil, err
	}
	SignCrtBytes, err := loadBytes(conf.SignCrtPath)
	if err != nil {
		return nil, err
	}

	sdkClient, err := cmsdk.NewChainClient(
		cmsdk.WithConfPath(conf.ChainConfigTemplatePath),
		cmsdk.WithChainClientChainId(conf.ChainID),
		cmsdk.WithChainClientOrgId(conf.OrgID),
		cmsdk.WithUserKeyBytes(SignKeyBytes),
		cmsdk.WithUserCrtBytes(SignCrtBytes),
		cmsdk.WithUserSignKeyBytes(SignKeyBytes),
		cmsdk.WithUserSignCrtBytes(SignCrtBytes),
	)
	if err != nil {
		return nil, err
	}
	return &txRequestBuilder{
		chainMakerSDK: sdkClient,
	}, nil
}

func loadBytes(path string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
