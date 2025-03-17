module github.com/hyperledger/fabric-samples/chaincode/transaction_contract/go

go 1.16

require (
	chainmaker.org/chainmaker-cross/utils v0.0.0-00010101000000-000000000000
	github.com/golang/mock v1.4.3
	github.com/golang/protobuf v1.3.3
	github.com/google/go-cmp v0.3.0 // indirect
	github.com/hyperledger/fabric-chaincode-go v0.0.0-20200424173110-d7076418f212
	github.com/hyperledger/fabric-contract-api-go v1.1.1
	github.com/hyperledger/fabric-protos-go v0.0.0-20200707132912-fee30f3ccd23
	github.com/stretchr/testify v1.5.1
	golang.org/x/sys v0.0.0-20190801041406-cbf593c0f2f3 // indirect
	google.golang.org/grpc v1.29.1 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)

replace chainmaker.org/chainmaker-cross/utils => ../../../module/utils
