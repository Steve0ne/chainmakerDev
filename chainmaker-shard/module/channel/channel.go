/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"chainmaker.org/chainmaker-cross/event"
)

type TransmissionChanType int32

const (
	InnerTransmissionChan TransmissionChanType = iota
	NetTransmissionChan
)

// TransmissionChannel Transmission channel
type TransmissionChannel interface {

	// GetChanType return channel types defined above
	GetChanType() TransmissionChanType

	// Deliver deliver message
	Deliver(eve *event.TransactionEventContext) error
}
