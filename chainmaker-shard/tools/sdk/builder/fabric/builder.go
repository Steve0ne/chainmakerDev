/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"fmt"

	"chainmaker.org/chainmaker-cross/adapter/fabric"
	"chainmaker.org/chainmaker-cross/sdk/builder"
	conf "chainmaker.org/chainmaker-cross/sdk/config"
	"chainmaker.org/chainmaker/common/json"
	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/txn"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type txRequestBuilder struct {
	fabricSDK *fabsdk.FabricSDK
	config    *conf.CrossChainConf
}

func (cb *txRequestBuilder) Build(in *builder.TxRequestBuildParam) ([]byte, error) {
	return createRawTransaction(cb, in.Contract.Name, in.Contract.Method, in.Contract.Params.Values())
}

func createRawTransaction(cb *txRequestBuilder, contractName, method string, args []string) ([]byte, error) {
	// get config
	//orgI, err := cb.config.GetExtraParamsByKey("org_name")
	//if err != nil {
	//	return nil, err
	//}
	//org, ok := orgI.(string)
	//if !ok {
	//	return nil ,fmt.Errorf("convert extra params user to string error")
	//}
	//
	//userI, err := cb.config.GetExtraParamsByKey("org_user")
	//if err != nil {
	//	return nil, err
	//}
	//user, ok := userI.(string)
	//if !ok {
	//	return nil ,fmt.Errorf("convert extra params user to string error")
	//}

	ccp := cb.fabricSDK.ChannelContext(cb.config.ChainID, fabsdk.WithUser("user"), fabsdk.WithOrg("org"))
	ctx, err := ccp()
	if err != nil {
		return nil, err
	}

	request := fab.ChaincodeInvokeRequest{
		ChaincodeID:  contractName,
		Fcn:          method,
		Args:         toArgs(args),
		TransientMap: nil,
		IsInit:       false,
	}

	reqBz, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	txnHeader, err := txn.NewHeader(ctx, cb.config.ChainID)
	if err != nil {
		return nil, err
	}

	proposal, err := txn.CreateChaincodeInvokeProposal(txnHeader, request)
	proposalBytes, err := proto.Marshal(proposal)
	if err != nil {
		return nil, err
	}
	signingMgr := ctx.SigningManager()
	if signingMgr == nil {
		return nil, fmt.Errorf("signing manager is nil")
	}

	signature, err := signingMgr.Sign(proposalBytes, ctx.PrivateKey())
	if err != nil {
		return nil, err
	}

	header := fabric.NewTxHeader(txnHeader.ChannelID(), string(txnHeader.TransactionID()))
	req := fabric.NewTxRequest(header, reqBz, proposalBytes, signature)
	bz, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func toArgs(params []string) [][]byte {
	var bytes [][]byte
	for _, param := range params {
		bytes = append(bytes, []byte(param))
	}
	return bytes
}

func NewTxRequestBuilder(conf *conf.CrossChainConf) (*txRequestBuilder, error) {
	provider := fromLocalConfig(conf)
	fabClient, err := fabsdk.New(
		provider,
	)
	if err != nil {
		return nil, err
	}
	return &txRequestBuilder{
		fabricSDK: fabClient,
		config:    conf,
	}, nil
}

func fromLocalConfig(conf *conf.CrossChainConf) core.ConfigProvider {
	return func() ([]core.ConfigBackend, error) {
		backend, err := newBackend()
		if err != nil {
			return nil, err
		}

		if conf.ChainConfigTemplatePath == "" {
			return nil, fmt.Errorf("filename is required")
		}

		// create new viper
		backend.configViper.SetConfigFile(conf.ChainConfigTemplatePath)

		// If a config file is found, read it in.
		err = backend.configViper.MergeInConfig()
		if err != nil {
			return nil, err
		}

		// set sign key and cert
		backend.configViper.Set("organizations", backend.configViper.Get("organizations"))
		backend.configViper.Set("organizations.org.mspid", conf.OrgID)
		backend.configViper.Set("organizations.org.users.user.key.path", conf.SignKeyPath)
		backend.configViper.Set("organizations.org.users.user.cert.path", conf.SignCrtPath)

		//backend.configViper.Set("channels." + conf.OrgID, signCertPath)

		return []core.ConfigBackend{backend}, nil
	}
}
