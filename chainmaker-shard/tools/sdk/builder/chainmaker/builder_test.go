/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package chainmaker

import (
	"testing"

	"chainmaker.org/chainmaker/pb-go/common"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"chainmaker.org/chainmaker-cross/mock"
	"chainmaker.org/chainmaker-cross/sdk/builder"
	conf "chainmaker.org/chainmaker-cross/sdk/config"
)

func TestChainMakerBuilder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mock.NewMockSDKInterface(ctrl)
	m.EXPECT().GetTxRequest(gomock.Eq("TransactionStable"), gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&common.TxRequest{
			//Header: &common.TxHeader{
			//	ChainId: "chain1",
			//	Sender: &accesscontrol.SerializedMember{
			//		OrgId: "wx-org1.chainmaker.org",
			//	},
			//},
			//Payload: []byte{10, 17, 84, 114, 97, 110, 115, 97, 99, 116, 105, 111, 110, 83, 116, 97, 98, 108, 101, 18, 7, 69, 120, 101, 99, 117, 116, 101},
		}, nil)
	txBuilder := txRequestBuilder{
		chainMakerSDK: m,
	}
	//txBuilder, err := NewTxRequestBuilder("/Users/leon/go/src/chainmaker.org/chainmaker-cross-chain/config/chainmaker/chainmaker_sdk1.yml")
	//require.Nil(t, err)
	reqBytes, err := txBuilder.Build(&builder.TxRequestBuildParam{
		CrossID: uuid.New().String(),
		Contract: &builder.Contract{
			Name:   "TransactionStable",
			Method: "Execute",
			Params: builder.NewParamsWithMap(map[string]string{}),
		},
	})
	require.Nil(t, err)
	require.Condition(t, func() bool {
		t.Log(reqBytes)
		return len(reqBytes) > 0
	})
}

func TestChainMakerParamBuilder(t *testing.T) {
	configM, err := conf.InitConfigByFilepath("../../config/template/cross_chain_sdk.yml")
	require.Nil(t, err)
	pbr := NewTxContractParamBuilder(configM.GetCrossChainConfByChainID("chain1"))
	buildParam := &builder.CrossTxBuildParam{
		CrossID:                  uuid.New().String(),
		Index:                    0,
		ExecuteBusinessContract:  builder.NewContract("ContractName", "execute", builder.NewParamsWithMap(map[string]string{"hello": "world"})),
		RollbackBusinessContract: builder.NewContract("ContractName", "rollback", builder.NewParamsWithMap(map[string]string{})),
	}
	options := builder.NewCrossBuildOptions()
	options.ProofKey = "123456"
	f := builder.WithUseProofKey(true)
	f(options)
	params, _ := pbr.BuildExecuteParam(buildParam, options.ParamOptions...)
	require.Condition(t, func() bool {
		t.Log(params)
		return params.Len() > 0
	})
	params, _ = pbr.BuildCommitParam(buildParam)
	require.Condition(t, func() bool {
		t.Log(params)
		return params.Len() > 0
	})
	params, _ = pbr.BuildRollbackParam(buildParam)
	require.Condition(t, func() bool {
		t.Log(params)
		return params.Len() > 0
	})
}
