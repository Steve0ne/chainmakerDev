/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package builder

type TxRequestBuildParam struct {
	CrossID string
	//ContractName string
	//Method       *Method
	Contract *Contract
	TxID     string
}

//TxRequestBuilder specifies the interface to build parameters (in) into different transaction request formats according to different chain platforms
type TxRequestBuilder interface {
	Build(in *TxRequestBuildParam) ([]byte, error)
	//CustomizedSignMethod() ([]byte, error)
}

//TxContractParamBuilder specifies the interface to build transaction contract parameters
//different chain platforms maybe has different ways to build
type TxContractParamBuilder interface {
	BuildExecuteParam(*CrossTxBuildParam, ...ParamsBuildOption) (*Params, error)
	BuildCommitParam(*CrossTxBuildParam) (*Params, error)
	BuildRollbackParam(*CrossTxBuildParam) (*Params, error)
}
