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

func TestNewCommitTransactionEvent(t *testing.T) {
	ce := NewCommitTransactionEvent("crossID", "chainID", []byte{})
	require.NotNil(t, ce)

	require.Equal(t, ce.GetType(), eventproto.TransactionEventType)
	require.Equal(t, ce.GetChainID(), "chainID")
	require.Equal(t, ce.GetPayload(), []byte{})
	require.Equal(t, ce.GetOpFunc(), CommitOpFunc)
}

func TestNewExecuteTransactionEvent(t *testing.T) {
	ee := NewExecuteTransactionEvent("crossID", "chainID", []byte{}, nil)
	require.NotNil(t, ee)

	require.Equal(t, ee.GetType(), eventproto.TransactionEventType)
	require.Equal(t, ee.GetChainID(), "chainID")
	require.Equal(t, ee.GetPayload(), []byte{})
	require.Equal(t, ee.GetOpFunc(), ExecuteOpFunc)
}

func TestNewRollbackTransactionEvent(t *testing.T) {
	re := NewRollbackTransactionEvent("crossID", "chainID", []byte{})
	require.NotNil(t, re)

	require.Equal(t, re.GetType(), eventproto.TransactionEventType)
	require.Equal(t, re.GetChainID(), "chainID")
	require.Equal(t, re.GetPayload(), []byte{})
	require.Equal(t, re.GetOpFunc(), RollbackOpFunc)
}
