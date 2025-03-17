
package fabric

import (
	"github.com/spf13/viper"
	"strings"
)

// defConfigBackend represents the default config backend
type defConfigBackend struct {
	configViper *viper.Viper
	//opts        options
}

func newBackend() (*defConfigBackend, error) {
	return &defConfigBackend{newViper("FABRIC_SDK")}, nil
}

func newViper(cmdRootPrefix string) *viper.Viper {
	myViper := viper.New()
	myViper.SetEnvPrefix(cmdRootPrefix)
	myViper.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	myViper.SetEnvKeyReplacer(replacer)
	return myViper
}

// Lookup gets the config item value by Key
func (c *defConfigBackend) Lookup(key string) (interface{}, bool) {
	value := c.configViper.Get(key)
	if value == nil {
		return nil, false
	}
	return value, true
}
