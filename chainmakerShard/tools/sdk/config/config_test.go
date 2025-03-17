/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitLocalConfigByFilepath(t *testing.T) {
	ymlPath := "./template/cross_chain_sdk.yml"
	configM, err := InitConfigByFilepath(ymlPath)
	require.NoError(t, err)
	t.Log(configM)
}
