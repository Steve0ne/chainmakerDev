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

func TestNewCrossResponse(t *testing.T) {
	// test default resp
	dcr := DefaultCrossResponse()
	require.Equal(t, *dcr, eventproto.CrossResponse{
		TxResponses: make([]*eventproto.CrossTxResponse, 0),
	})

	// test cross resp
	cr := NewCrossResponse("crossID", 0, "msg")
	cr.AddTxResponse(&eventproto.CrossTxResponse{})
	require.Equal(t, cr.TxResponses, []*eventproto.CrossTxResponse{NewCrossTxResponse("", "", 0, 0, nil)})
	require.Equal(t, cr.GetType(), eventproto.CrossRespEventType)
}

func TestNewTxResponse(t *testing.T) {
	txResp := NewTxResponse("chainID", "txKey", 1, 0, nil, nil)
	require.NotNil(t, txResp)
}

func TestNewCommonTxResponse(t *testing.T) {
	txResp := NewTxResponse("chainID", "txKey", 1, 0, nil, nil)
	cTxResp := NewCommonTxResponse(txResp, 1, "msg")
	require.NotNil(t, cTxResp)

	require.Equal(t, cTxResp.IsSuccess(), false)
}

func TestNewContract(t *testing.T) {
	c := NewContract("name", "version", "method", []byte{})
	require.NotNil(t, c)

	cp := NewContractParameter("key", "value")

	c.AddParameter(cp)
	require.Equal(t, c.Parameters, []*eventproto.ContractParameter{cp})
	c.AddParameters([]*eventproto.ContractParameter{cp})
	require.Equal(t, c.Parameters, []*eventproto.ContractParameter{cp, cp})

	v := NewContractValue("value")
	require.Equal(t, v, &eventproto.ContractParameter{
		Key:   "",
		Value: "value",
	})
}

func TestNewProofResponseByProof(t *testing.T) {
	p := NewProof("chainID", "txKey", 1, 0, nil, []byte{})
	pr := NewProofResponseByProof("crossID", "chainID", "msg", 0, ExecuteOpFunc, p)

	pr.SetKey("test")
	require.Equal(t, pr.Key, "test")
	require.Equal(t, pr.GetKey(), "test")
	require.Equal(t, pr.GetType(), eventproto.ProofRespEventType)
	require.Equal(t, pr.IsSuccess(), true)

	pr.Done("chainID", "txKey", 1, 0, nil, []byte{})
	pr.DoneError("msg")

	//pr.Wait(time.Second * 1)
}
