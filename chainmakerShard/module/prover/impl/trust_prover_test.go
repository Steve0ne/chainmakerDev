/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package impl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrustProver(t *testing.T) {
	// new prover
	tp := NewTrustProver([]string{"chainID1", "chainID2"})
	require.NotNil(t, tp)

	// test infos to proof
	p, err := tp.ToProof("ChainID", "txKey", 0, 0, nil, []byte{})
	require.NoError(t, err)
	require.NotNil(t, p)

	// test proof
	ok, err := tp.Prove(nil)
	require.NoError(t, err)
	require.Equal(t, ok, true)

	// test get type
	tpType := tp.GetType()
	require.Equal(t, tpType, TrustProverType)

	// test get chain ids
	ids := tp.GetChainIDs()
	require.NotNil(t, ids)
}
