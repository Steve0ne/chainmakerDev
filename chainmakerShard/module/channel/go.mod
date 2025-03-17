module chainmaker.org/chainmaker-cross/channel

go 1.15

require (
	chainmaker.org/chainmaker-cross/conf v0.0.0
	chainmaker.org/chainmaker-cross/event v0.0.0
	chainmaker.org/chainmaker-cross/logger v0.0.0
	chainmaker.org/chainmaker-cross/net v0.0.0
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker-cross/utils v0.0.0
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.16.0
)

replace (
	chainmaker.org/chainmaker-cross/conf => ../conf
	chainmaker.org/chainmaker-cross/event => ../event
	chainmaker.org/chainmaker-cross/logger => ../logger
	chainmaker.org/chainmaker-cross/net => ../net
	chainmaker.org/chainmaker-cross/pb/protogo => ../pb/protogo
	chainmaker.org/chainmaker-cross/utils => ../utils
)
