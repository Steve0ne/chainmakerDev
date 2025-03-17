module chainmaker.org/chainmaker-cross/handler

go 1.15

require (
	chainmaker.org/chainmaker-cross/adapter v0.0.0
	chainmaker.org/chainmaker-cross/conf v0.0.0
	chainmaker.org/chainmaker-cross/event v0.0.0
	chainmaker.org/chainmaker-cross/logger v0.0.0
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker-cross/router v0.0.0
	chainmaker.org/chainmaker-cross/store v0.0.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
)

replace (
	chainmaker.org/chainmaker-cross/adapter => ../adapter
	chainmaker.org/chainmaker-cross/channel => ../channel
	chainmaker.org/chainmaker-cross/conf => ../conf
	chainmaker.org/chainmaker-cross/event => ../event
	chainmaker.org/chainmaker-cross/logger => ../logger
	chainmaker.org/chainmaker-cross/net => ../net
	chainmaker.org/chainmaker-cross/pb/protogo => ../pb/protogo
	chainmaker.org/chainmaker-cross/prover => ../prover
	chainmaker.org/chainmaker-cross/router => ../router
	chainmaker.org/chainmaker-cross/store => ../store
	chainmaker.org/chainmaker-cross/utils => ../utils
)
