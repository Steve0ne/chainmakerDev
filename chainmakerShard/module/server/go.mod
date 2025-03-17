module chainmaker.org/chainmaker-cross/server

go 1.16

require (
	chainmaker.org/chainmaker-cross/adapter v0.0.0
	chainmaker.org/chainmaker-cross/event v0.0.0
	chainmaker.org/chainmaker-cross/handler v0.0.0
	chainmaker.org/chainmaker-cross/listener v0.0.0
	chainmaker.org/chainmaker-cross/logger v0.0.0
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker-cross/prover v0.0.0
	chainmaker.org/chainmaker-cross/router v0.0.0
	chainmaker.org/chainmaker-cross/store v0.0.0
	chainmaker.org/chainmaker-cross/transaction v0.0.0
	chainmaker.org/chainmaker-cross/utils v0.0.0
	go.uber.org/zap v1.16.0
)

replace (
	chainmaker.org/chainmaker-cross/adapter => ../adapter
	chainmaker.org/chainmaker-cross/channel => ../channel
	chainmaker.org/chainmaker-cross/conf => ../conf
	chainmaker.org/chainmaker-cross/event => ../event
	chainmaker.org/chainmaker-cross/handler => ../handler
	chainmaker.org/chainmaker-cross/listener => ../listener
	chainmaker.org/chainmaker-cross/logger => ../logger
	chainmaker.org/chainmaker-cross/net => ./../net
	chainmaker.org/chainmaker-cross/pb/protogo => ../pb/protogo
	chainmaker.org/chainmaker-cross/prover => ../prover
	chainmaker.org/chainmaker-cross/router => ../router
	chainmaker.org/chainmaker-cross/store => ../store
	chainmaker.org/chainmaker-cross/transaction => ../transaction
	chainmaker.org/chainmaker-cross/utils => ../utils
)
