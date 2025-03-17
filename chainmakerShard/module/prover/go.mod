module chainmaker.org/chainmaker-cross/prover

go 1.15

require (
	chainmaker.org/chainmaker-cross/conf v0.0.0
	chainmaker.org/chainmaker-cross/event v0.0.0
	chainmaker.org/chainmaker-cross/logger v0.0.0
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker/spv/v2 v2.1.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
)

replace (
	chainmaker.org/chainmaker-cross/conf => ../conf
	chainmaker.org/chainmaker-cross/event => ../event
	chainmaker.org/chainmaker-cross/logger => ../logger
	chainmaker.org/chainmaker-cross/pb/protogo => ../pb/protogo
	chainmaker.org/chainmaker-cross/utils => ../utils
)
