/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package coder

import (
	"fmt"
	"testing"

	"chainmaker.org/chainmaker-cross/event"
)

func TestJsonBinaryMarshal(t *testing.T) {
	crossTx0 := event.NewCrossTx("chain1", 0, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))
	crossTx1 := event.NewCrossTx("chain2", 1, []byte("executePayload"), []byte("commitPayload"), []byte("rollbackPayload"))
	crossTxs := make([]*event.CrossTx, 0)
	crossTxs = append(crossTxs, crossTx0, crossTx1)
	crossEvent := event.NewCrossEvent(crossTxs)
	bytes, err := JsonBinaryMarshal(crossEvent.GetType(), crossEvent)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("marshal's result's length = %v \n", len(bytes))
	// 进行unmarshal
	var newCrossEvent = &event.CrossEvent{}
	err = JsonBinaryUnmarshal(bytes, byte(crossEvent.GetType()), newCrossEvent)
	if err != nil {
		t.Error(err)
	}
	// check
	if newCrossEvent.CrossID != crossEvent.CrossID {
		t.Error("crossID is not equal")
	}
	if newCrossEvent.Timestamp != crossEvent.Timestamp {
		t.Error("timestamp is not equal")
	}
}
