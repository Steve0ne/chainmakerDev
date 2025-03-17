/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInnerChannel(t *testing.T) {
	// test get channel
	ic := GetInnerChannel()
	require.NotNil(t, ic)

	// test get channel type
	cType := ic.GetChanType()
	require.Equal(t, cType, InnerTransmissionChan)

	// test deliver
	err := ic.Deliver(nil)
	require.NoError(t, err)

	// test get chan
	c := ic.GetChan()
	require.NotNil(t, c)
}
