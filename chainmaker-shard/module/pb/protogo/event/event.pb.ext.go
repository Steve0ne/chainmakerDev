/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/
package event

// Type is the base type for all event type by each event implements
type EventType byte

// event types define
const (
	CrossEventType       EventType = iota // cross-chain event type, send from Client to Proxy
	CrossEventSearchType                  // cross-chain event type, send from Client to Proxy
	TransactionEventType                  // transaction event type, transaction messages connect each Proxy
	CrossTxType                           // para-chain event type, an interface, implied by each para-chain itself
	CrossRespEventType
	ProofRespEventType
	TransactionCtxEventType
	TxProofType
)

// SetExtra set extra
func (c *CrossEvent) SetExtra(extra []byte) {
	c.Extra = extra
}

// SetVersion set version
func (c *CrossEvent) SetVersion(version string) {
	c.Version = version
}

// SetTimestamp set timestamp
func (c *CrossEvent) SetTimestamp(timestamp int64) {
	c.Timestamp = timestamp
}

// SetCrossID set cross-id
func (c *CrossEvent) SetCrossID(crossID string) {
	c.CrossId = crossID
}

// GetType return type of this event
func (c *CrossEvent) GetType() EventType {
	return CrossEventType
}

// GetPkgTxEvents return set of tx-event
func (c *CrossEvent) GetPkgTxEvents() *CrossTxs {
	return c.TxEvents
}

// IsValid check txs
func (c *CrossEvent) IsValid() bool {
	for i, tx := range c.TxEvents.GetCrossTxs() {
		if int(tx.Index) != i {
			return false
		}
	}
	return true
}

// GetChainIDs return array of chain-id
func (c *CrossEvent) GetChainIDs() []string {
	chainIDs := make([]string, 0)
	if c.TxEvents != nil {
		for _, crossTx := range c.TxEvents.GetCrossTxs() {
			innerCrossTx := crossTx
			chainIDs = append(chainIDs, innerCrossTx.GetChainID())
		}
	}
	return chainIDs
}

// return cross-id
func (c *CrossEvent) GetCrossID() string {
	return c.CrossId
}

// GetType return type of this event
func (c *CrossSearchEvent) GetType() EventType {
	return CrossEventSearchType
}

// return cross-id
func (c *CrossSearchEvent) GetCrossID() string {
	return c.CrossId
}

// GetType return type of event
func (p *CrossTx) GetType() EventType {
	return CrossTxType
}

// GetChainID return chain-id
func (p *CrossTx) GetChainID() string {
	return p.GetChainId()
}

// Len return length of array
func (ps *CrossTxs) Len() int {
	return len(ps.Events)
}

func (ps *CrossTxs) Append(txs ...*CrossTx) {
	ps.Events = append(ps.Events, txs...)
}

// Less compare object
func (ps *CrossTxs) Less(i, j int) bool {
	return ps.Events[i].Index < ps.Events[j].Index
}

// Swap swap object
func (ps *CrossTxs) Swap(i, j int) {
	ps.Events[i], ps.Events[j] = ps.Events[j], ps.Events[i]
}

// GetCrossTxs return cross txs
func (ps *CrossTxs) GetCrossTxs() []*CrossTx {
	return ps.Events
}

// GetType return type of this event
func (p *Proof) GetType() EventType {
	return TxProofType
}

// GetChainID return chain-id
func (p *Proof) GetChainID() string {
	return p.GetChainId()
}

type VerifiedProof struct {
	TxProof        *Proof
	VerifiedResult bool
	ProverType     string
	Identity       string // 预留，暂不实现
}

func NewVerifiedProof(txProof *Proof, verifiedResult bool, proverType, identity string) *VerifiedProof {
	return &VerifiedProof{
		TxProof:        txProof,
		VerifiedResult: verifiedResult,
		ProverType:     proverType,
		Identity:       identity,
	}
}

func (c *ContractInfo) AddParameter(parameter *ContractParameter) {
	c.Parameters = append(c.Parameters, parameter)
}

// AddParameters add set of parameter to contractInfo
func (c *ContractInfo) AddParameters(parameter []*ContractParameter) {
	c.Parameters = append(c.Parameters, parameter...)
}

func (te *TransactionEvent) GetType() EventType {
	return TransactionEventType
}

// GetCrossID return cross-id
func (te *TransactionEvent) GetCrossID() string {
	return te.GetCrossId()
}

// GetChainID return chain-id
func (te *TransactionEvent) GetChainID() string {
	return te.GetChainId()
}

func (te *TransactionEvent) NeedProve() bool {
	return te.TxProof != nil
}

func (te *TransactionEvent) SetProofKey(key string) {
	if te != nil {
		te.ProofKey = key
	}
}
