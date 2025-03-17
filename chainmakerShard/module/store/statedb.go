/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	storetypes "chainmaker.org/chainmaker-cross/store/types"
)

// StateDb support modules store and load transaction data
type StateDB interface {

	// StartCross start the cross transaction
	StartCross(crossID string, content []byte) error

	// FinishCross finish the cross transaction
	FinishCross(crossID string, result []byte, state storetypes.State) error

	// ReadCross load content for the crossID
	ReadCross(crossID string) ([]byte, error)

	// WriteCross write the content for the crossID
	WriteCross(crossID string, content []byte) error

	// WriteChainIDs write the relationship between crossID and chain ids
	WriteChainIDs(crossID string, chainIDs []string) error

	// ReadChainIDs read chain ids for the crossID
	ReadChainIDs(crossID string) ([]string, bool)

	// WriteCrossState write the total state for the crossID
	WriteCrossState(crossID string, state storetypes.State) error

	// ReadCrossState read the total state for the crossID
	ReadCrossState(crossID string) (storetypes.State, []byte, bool)

	// WriteChainCrossState write the state for the crossID and chain
	WriteChainCrossState(crossID, chainID string, state storetypes.State, content []byte) error

	// FinishChainCrossState finish the cross transaction for chain
	FinishChainCrossState(crossID, chainID string, result []byte, state storetypes.State) error

	// ReadChainCrossState read the state for crossID and chain
	ReadChainCrossState(crossID, chainID string) (storetypes.State, []byte, bool)

	// ReadUnfinishedCrossIDs read unfinished crossID array
	ReadUnfinishedCrossIDs() []string

	// DeleteCrossIDFromUnfinished delete crossID from the unfinished crossID array
	DeleteCrossIDFromUnfinished(crossID string) error

	// Close close the state database
	Close()
}
