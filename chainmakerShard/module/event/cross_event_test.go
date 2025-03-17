/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package event

import (
	"testing"
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"github.com/stretchr/testify/require"
)

func TestCrossEvent(t *testing.T) {
	// test new cross event

	ce := NewCrossEvent(nil)
	require.NotNil(t, ce)
	ct := NewCrossTx("chainID", 0, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))
	ece := NewEmptyCrossEvent()
	require.NotNil(t, ece)
	ece.TxEvents.Events = append(ece.TxEvents.Events, ct)

	// test set methods
	ece.SetExtra([]byte("test"))
	ece.SetVersion("test")
	ece.SetTimestamp(time.Now().Unix())
	ece.SetCrossID("test")

	// test get methods
	require.Equal(t, ece.GetType(), eventproto.CrossEventType)
	require.Equal(t, ece.GetPkgTxEvents(), ece.TxEvents)
	require.Equal(t, ece.GetChainIDs(), []string{"chainID"})
	require.Equal(t, ece.GetCrossID(), "test")
	require.NotNil(t, ece.GetPkgTxEvents())

	// test check method
	require.Equal(t, ece.IsValid(), true)

	// test abnormal
	// test crossEvent invalid
	ct2 := NewCrossTx("chainID", 2, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))
	ece.TxEvents.Events = append(ece.TxEvents.Events, ct2)
	require.Equal(t, ece.IsValid(), false)
}

func TestCrossSearch(t *testing.T) {
	cse := eventproto.CrossSearchEvent{
		CrossId: "test",
	}

	require.Equal(t, cse.GetType(), eventproto.CrossEventSearchType)
	require.Equal(t, cse.GetCrossID(), "test")
}
