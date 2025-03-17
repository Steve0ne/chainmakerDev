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

func TestGetProofResponseContexts(t *testing.T) {
	prcs := GetProofResponseContexts()
	require.Equal(t, prcs, &ProofResponseContexts{
		contexts: make(map[string]*ProofResponseContext),
	})

	// test register
	pr := NewProofResponse("crossID", "chainID", ExecuteOpFunc)
	prc := NewProofResponseContext(pr)
	prcs.Register(prc)
	require.NotNil(t, prcs.contexts)

	// test done
	prcs.Done(pr.Key, pr.GetChainID(), pr.GetTxKey(), pr.GetBlockHeight(), pr.GetIndex(), pr.GetContract(), pr.GetExtra())
	prcs.DoneByProofResp(pr)
	prcs.DoneError(pr.Key, pr.Msg)

	// test remove
	prcs.Remove(pr.Key)
	require.Equal(t, prcs.contexts, map[string]*ProofResponseContext{})
}

func TestNewTransactionEventContext(t *testing.T) {
	tec := NewTransactionEventContext("key", nil)
	var x *eventproto.TransactionEvent

	require.Equal(t, tec.GetType(), eventproto.TransactionCtxEventType)
	require.Equal(t, tec.GetEvent(), x)
	require.Equal(t, tec.GetKey(), "key")
}
