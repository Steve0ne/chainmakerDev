module chainmaker.org/chainmaker-cross/adapter

go 1.15

require (
	chainmaker.org/chainmaker-cross/conf v0.0.0
	chainmaker.org/chainmaker-cross/event v0.0.0
	chainmaker.org/chainmaker-cross/logger v0.0.0
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker-cross/prover v0.0.0
	chainmaker.org/chainmaker-cross/utils v0.0.0
	chainmaker.org/chainmaker/pb-go/v2 v2.0.0
	chainmaker.org/chainmaker/sdk-go/v2 v2.0.0
	github.com/Rican7/retry v0.1.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
)

replace (
	chainmaker.org/chainmaker-cross/conf => ../conf
	chainmaker.org/chainmaker-cross/event => ../event
	chainmaker.org/chainmaker-cross/logger => ../logger
	chainmaker.org/chainmaker-cross/pb/protogo => ../pb/protogo
	chainmaker.org/chainmaker-cross/pb/protogo/event => ../pb/protogo/event
	chainmaker.org/chainmaker-cross/prover => ../prover
	chainmaker.org/chainmaker-cross/utils => ../utils
)
