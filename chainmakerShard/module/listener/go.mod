module chainmaker.org/chainmaker-cross/listener

go 1.15

require (
	chainmaker.org/chainmaker-cross/channel v0.0.0
	chainmaker.org/chainmaker-cross/conf v0.0.0
	chainmaker.org/chainmaker-cross/event v0.0.0
	chainmaker.org/chainmaker-cross/handler v0.0.0
	chainmaker.org/chainmaker-cross/logger v0.0.0
	chainmaker.org/chainmaker-cross/net v0.0.0
	chainmaker.org/chainmaker-cross/pb/protogo v0.0.0
	chainmaker.org/chainmaker-cross/utils v0.0.0
	github.com/gin-gonic/gin v1.7.2
	github.com/libp2p/go-libp2p v0.13.0
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.36.0
)

replace (
	chainmaker.org/chainmaker-cross/adapter => ../adapter
	chainmaker.org/chainmaker-cross/channel => ../channel
	chainmaker.org/chainmaker-cross/conf => ../conf
	chainmaker.org/chainmaker-cross/event => ../event
	chainmaker.org/chainmaker-cross/handler => ../handler
	chainmaker.org/chainmaker-cross/logger => ../logger
	chainmaker.org/chainmaker-cross/net => ./../net
	chainmaker.org/chainmaker-cross/pb/protogo => ../pb/protogo
	chainmaker.org/chainmaker-cross/prover => ../prover
	chainmaker.org/chainmaker-cross/router => ../router
	chainmaker.org/chainmaker-cross/store => ../store
	chainmaker.org/chainmaker-cross/utils => ../utils
)
