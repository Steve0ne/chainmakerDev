/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package fabric

// TxHeader
type TxHeader struct {
	// blockchain identifier
	ChainId string `json:"chain_id,omitempty"`
	// transaction id set by sender, should be unique
	TxId string `json:"tx_id,omitempty"`
}

// NewTxHeader create TxHeader
func NewTxHeader(chainID, txID string) *TxHeader {
	return &TxHeader{
		ChainId: chainID,
		TxId:    txID,
	}
}

// TxRequest
type TxRequest struct {
	Header *TxHeader `json:"header,omitempty"`
	// request params bytes
	Request []byte `json:"request,omitempty"`
	// payload of the request, can be unmarshalled according to tx_type in header
	Payload []byte `json:"payload,omitempty"`
	// signature of [header bytes || payload bytes]
	Signature []byte `json:"signature,omitempty"`
}

// NewTxRequest create TxRequest
func NewTxRequest(header *TxHeader, request []byte, payload []byte, signature []byte) *TxRequest {
	return &TxRequest{
		Header:    header,
		Request:   request,
		Payload:   payload,
		Signature: signature,
	}
}
