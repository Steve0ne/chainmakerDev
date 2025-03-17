module chainmaker.org/chainmaker-cross/net

go 1.15

require (
	chainmaker.org/chainmaker-cross/conf v0.0.0
	chainmaker.org/chainmaker-cross/event v0.0.0
	chainmaker.org/chainmaker-cross/logger v0.0.0
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker-cross/utils v0.0.0
	github.com/Rican7/retry v0.1.0
	github.com/btcsuite/btcd v0.21.0-beta
	github.com/gin-gonic/gin v1.7.2
	github.com/google/uuid v1.2.0
	github.com/libp2p/go-libp2p v0.13.0
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
)

replace (
	chainmaker.org/chainmaker-cross/conf => ../conf
	chainmaker.org/chainmaker-cross/event => ../event
	chainmaker.org/chainmaker-cross/logger => ../logger
	chainmaker.org/chainmaker-cross/pb/protogo => ../pb/protogo
	chainmaker.org/chainmaker-cross/utils => ../utils
)
