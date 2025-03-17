/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package event

// AddTxResponse add tx response
func (c *CrossResponse) AddTxResponse(response *CrossTxResponse) {
	c.TxResponses = append(c.TxResponses, response)
}

// GetType return type of this event
func (c *CrossResponse) GetType() EventType {
	return CrossRespEventType
}

func (tr *TxResponse) GetChainID() string {
	return tr.GetChainId()
}
