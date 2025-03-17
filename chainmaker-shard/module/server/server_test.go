/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package server

import (
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	"fmt"
	"time"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/utils"
)

const (
	chain1 = "chain1"
	chain2 = "chain2"
)

func sendCrossEventTicker(count int, eventCh chan event.Event) {
	go func(count int, eventCh chan event.Event) {
		// 定时来处理
		for i := 0; i < count; i++ {
			fmt.Println("I will create new cross event and send it")
			crossEvent := innerNewCrossEvent()
			eventCh <- crossEvent
			time.Sleep(60 * time.Second)
		}
	}(count, eventCh)
}

func innerNewCrossEvent() *eventproto.CrossEvent {
	crossTx0 := event.NewCrossTx(chain1, 0, []byte("execute"), []byte("commit"), []byte("rollback"))
	crossTx1 := event.NewCrossTx(chain2, 1, []byte("execute"), []byte("commit"), []byte("rollback"))
	crossTxs := make([]*eventproto.CrossTx, 0)
	crossTxs = append(crossTxs, crossTx0, crossTx1)

	crossEvent := &eventproto.CrossEvent{
		CrossId:   utils.NewUUID(),
		TxEvents:  event.NewCrossTxs(crossTxs),
		Version:   "v1.0.0",
		Timestamp: time.Now().Unix(),
		Extra:     []byte("chainmaker"),
	}
	return crossEvent
}
