package conf

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCrossChainConf_GetExtraParamsByKey(t *testing.T) {
	conf := CrossChainConf{
		ExtraParams: map[string]interface{}{
			"string": "string",
			"slice": []string{"1", "2"},
		},
	}
	str, err := conf.GetExtraParamsByKey("string")
	require.NoError(t, err)
	require.NotNil(t, str)

	sli, err := conf.GetExtraParamsByKey("slice")
	require.NoError(t, err)
	require.NotNil(t, sli)
}