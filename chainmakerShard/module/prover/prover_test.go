/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package prover

import (
	"fmt"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/prover/impl"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestProver(t *testing.T) {
	// set config
	conf.Config.ProverConfigs = []*conf.ProverConfig{
		{
			Provider:   "trust",
			ChainIDs:   []string{"ChainID"},
			ConfigPath: "",
		},
		//{
		//	"spv",
		//	[]string{"ChainID"},
		//	p,
		//},
	}
	// test init prover
	InitProvers()

	// test get dispatcher
	pd := GetProverDispatcher()

	// test register
	pd.Register(impl.NewTrustProver([]string{"ChainID"}))

	// test to proof
	p1, err := pd.ToProof("ChainID", "txKey", 0, 0, nil, []byte{})
	require.NoError(t, err)
	require.NotNil(t, p1)
	p2, err := pd.ToProof("ChainIDNotExist", "txKey", 0, 0, nil, []byte{})
	require.NoError(t, err)
	require.NotNil(t, p2)

	// test prove
	ok, err := pd.Prove(p1)
	require.Nil(t, err)
	require.Equal(t, ok, true)
	ok, err = pd.Prove(p2)
	require.Equal(t, err, fmt.Errorf("can not find prover for chainID [ChainIDNotExist]"))
	require.Equal(t, ok, false)

	// test get prover
	prover, ok := pd.GetProver("ChainID")
	require.Equal(t, ok, true)
	require.NotNil(t, prover)
}
