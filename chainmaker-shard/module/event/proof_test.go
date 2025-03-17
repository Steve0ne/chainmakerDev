/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"testing"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"github.com/stretchr/testify/require"
)

func TestNewProof(t *testing.T) {
	p := NewProof("chainID", "txKey", 1, 0, nil, []byte{})
	var x *eventproto.ContractInfo
	require.Equal(t, p.GetType(), eventproto.TxProofType)
	require.Equal(t, p.GetChainID(), "chainID")
	require.Equal(t, p.GetTxKey(), "txKey")
	require.Equal(t, p.GetBlockHeight(), int64(1))
	require.Equal(t, p.GetIndex(), int32(0))
	require.Equal(t, p.GetContract(), x)
	require.Equal(t, p.GetExtra(), []byte{})
}
