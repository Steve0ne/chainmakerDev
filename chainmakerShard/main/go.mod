module chainmaker.org/chainmaker-cross/main

go 1.15

require (
	chainmaker.org/chainmaker-cross/conf v0.0.0
	chainmaker.org/chainmaker-cross/logger v0.0.0
	chainmaker.org/chainmaker-cross/server v0.0.0
	github.com/google/martian v2.1.0+incompatible
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
)

replace (
	chainmaker.org/chainmaker-cross/adapter => ../module/adapter
	chainmaker.org/chainmaker-cross/adapter/fabric => ../module/adapter/fabric
	chainmaker.org/chainmaker-cross/channel => ../module/channel
	chainmaker.org/chainmaker-cross/conf => ../module/conf
	chainmaker.org/chainmaker-cross/event => ../module/event
	chainmaker.org/chainmaker-cross/handler => ../module/handler
	chainmaker.org/chainmaker-cross/listener => ../module/listener
	chainmaker.org/chainmaker-cross/logger => ../module/logger
	chainmaker.org/chainmaker-cross/net => ./../module/net
	chainmaker.org/chainmaker-cross/pb/protogo => ../module/pb/protogo
	chainmaker.org/chainmaker-cross/prover => ../module/prover
	chainmaker.org/chainmaker-cross/router => ../module/router
	chainmaker.org/chainmaker-cross/server => ../module/server
	chainmaker.org/chainmaker-cross/store => ../module/store
	chainmaker.org/chainmaker-cross/transaction => ../module/transaction
	chainmaker.org/chainmaker-cross/utils => ../module/utils
)
