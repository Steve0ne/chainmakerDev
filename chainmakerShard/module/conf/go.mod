module chainmaker.org/chainmaker-cross/conf

go 1.15

require (
	chainmaker.org/chainmaker-cross/logger v0.0.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.4.0
)

replace chainmaker.org/chainmaker-cross/logger => ../logger
