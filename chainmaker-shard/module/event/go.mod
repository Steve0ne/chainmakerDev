module chainmaker.org/chainmaker-cross/event

go 1.15

require (
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker-cross/utils v0.0.0
	github.com/json-iterator/go v1.1.10
	github.com/stretchr/testify v1.4.0
	go.uber.org/zap v1.16.0
)

replace (
	chainmaker.org/chainmaker-cross/pb/protogo => ../pb/protogo
	chainmaker.org/chainmaker-cross/utils => ../utils
)
