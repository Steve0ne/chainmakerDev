/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"sort"
	"testing"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestInitLog(t *testing.T) {
	InitLog(nil)
	var x *zap.SugaredLogger
	require.Equal(t, log, x)
}

func TestCrossTx(t *testing.T) {
	ct := NewCrossTx("chainID", 0, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))

	// test get method
	require.Equal(t, ct.GetType(), eventproto.CrossTxType)
	require.Equal(t, ct.GetChainID(), "chainID")
	require.Equal(t, string(ct.GetExecutePayload()), "executePayload")
	require.Equal(t, string(ct.GetCommitPayload()), "commitPayload")
	require.Equal(t, string(ct.GetRollbackPayload()), "rollbackPayload")
}

func TestCrossTxs(t *testing.T) {
	// create ct
	ct1 := NewCrossTx("chainID", 0, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))
	ct2 := NewCrossTx("chainID", 1, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))
	ct3 := NewCrossTx("chainID", 2, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))

	// test get cross txs
	cts := NewCrossTxs([]*eventproto.CrossTx{ct1, ct3, ct2})
	require.Equal(t, cts.GetCrossTxs(), cts.Events)
	// test sort
	sort.Sort(cts)
	require.Equal(t, cts, NewCrossTxs([]*eventproto.CrossTx{ct1, ct2, ct3}))
}
