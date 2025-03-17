/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package coder

import (
	"testing"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/event"
	"github.com/stretchr/testify/require"
)

func TestGetEventCoderTools(t *testing.T) {
	var coder event.EventCoder
	var ok bool
	var err error
	var bz []byte
	// test init coder
	ect := GetEventCoderTools()
	require.NotNil(t, ect)

	// test CrossEventType coder
	coder, ok = ect.GetDefaultCoder(eventproto.CrossEventType)
	require.Equal(t, ok, true)
	CE := event.NewEmptyCrossEvent()
	bz, err = coder.MarshalToBinary(CE)
	require.NoError(t, err)
	ce, err := coder.UnmarshalFromBinary(bz)
	require.NoError(t, err)
	require.Equal(t, CE, ce)

	// test TransactionEventType coder
	coder, ok = ect.GetDefaultCoder(eventproto.TransactionEventType)
	require.Equal(t, ok, true)
	ETE := event.NewExecuteTransactionEvent("crossID", "chainID", []byte("payload"), "", nil)
	bz, err = coder.MarshalToBinary(ETE)
	require.NoError(t, err)
	ete, err := coder.UnmarshalFromBinary(bz)
	require.NoError(t, err)
	require.Equal(t, ETE, ete)

	CTE := event.NewCommitTransactionEvent("crossID", "chainID", []byte("payload"))
	bz, err = coder.MarshalToBinary(CTE)
	require.NoError(t, err)
	cte, err := coder.UnmarshalFromBinary(bz)
	require.NoError(t, err)
	require.Equal(t, CTE, cte)

	RTE := event.NewRollbackTransactionEvent("crossID", "chainID", []byte("payload"))
	bz, err = coder.MarshalToBinary(RTE)
	require.NoError(t, err)
	rte, err := coder.UnmarshalFromBinary(bz)
	require.NoError(t, err)
	require.Equal(t, RTE, rte)

	// test CrossTxType coder
	coder, ok = ect.GetDefaultCoder(eventproto.CrossTxType)
	require.Equal(t, ok, true)
	CT := event.NewCrossTx("chainID", 0, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))
	bz, err = coder.MarshalToBinary(CT)
	require.NoError(t, err)
	ct, err := coder.UnmarshalFromBinary(bz)
	require.NoError(t, err)
	require.Equal(t, CT, ct)

	// test CrossRespEventType coder
	coder, ok = ect.GetDefaultCoder(eventproto.CrossRespEventType)
	require.Equal(t, ok, true)
	CR := event.NewCrossResponse("crossID", 0, "msg")
	bz, err = coder.MarshalToBinary(CR)
	require.NoError(t, err)
	cr, err := coder.UnmarshalFromBinary(bz)
	require.NoError(t, err)
	require.Equal(t, CR, cr)

	// test ProofRespEventType coder
	coder, ok = ect.GetDefaultCoder(eventproto.ProofRespEventType)
	require.Equal(t, ok, true)
	PR := event.NewProofResponse("crossID", "chainID", 0)
	bz, err = coder.MarshalToBinary(PR)
	require.NoError(t, err)
	pr, err := coder.UnmarshalFromBinary(bz)
	require.NoError(t, err)
	require.Equal(t, PR.GetType(), pr.GetType())

	// test TransactionCtxEventType coder
	coder, ok = ect.GetDefaultCoder(eventproto.TransactionCtxEventType)
	require.Equal(t, ok, true)
	TEC := event.NewTransactionEventContext("key", nil)
	bz, err = coder.MarshalToBinary(TEC)
	require.NoError(t, err)
	tec, err := coder.UnmarshalFromBinary(bz)
	require.NoError(t, err)
	require.Equal(t, TEC, tec)
}
