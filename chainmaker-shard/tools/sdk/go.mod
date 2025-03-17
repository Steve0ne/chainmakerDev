module chainmaker.org/chainmaker-cross/sdk

go 1.16

require (
	chainmaker.org/chainmaker-cross/adapter v0.0.0
	chainmaker.org/chainmaker-cross/conf v0.0.0
	chainmaker.org/chainmaker-cross/event v0.0.0
	chainmaker.org/chainmaker-cross/mock v0.0.0-00010101000000-000000000000
	chainmaker.org/chainmaker-cross/net v0.0.0
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker/common v0.0.0-20210722032200-380ced605d25
	chainmaker.org/chainmaker/pb-go v0.0.0-20210719032153-653bd8436ef6
	chainmaker.org/chainmaker/pb-go/v2 v2.0.0
	chainmaker.org/chainmaker/sdk-go/v2 v2.0.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
)

replace (
	chainmaker.org/chainmaker-cross/adapter => ../../module/adapter
	chainmaker.org/chainmaker-cross/conf => ../../module/conf
	chainmaker.org/chainmaker-cross/event => ../../module/event
	chainmaker.org/chainmaker-cross/logger => ../../module/logger
	chainmaker.org/chainmaker-cross/mock => ../../mock
	chainmaker.org/chainmaker-cross/net => ../../module/net
	chainmaker.org/chainmaker-cross/pb/protogo => ../../module/pb/protogo
	chainmaker.org/chainmaker-cross/prover => ../../module/prover
	chainmaker.org/chainmaker-cross/sdk => ./
	chainmaker.org/chainmaker-cross/utils => ../../module/utils
)
