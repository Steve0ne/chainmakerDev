/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"fmt"
	"sort"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/stretchr/testify/require"
)

const (
	configFile     = "./fabric_sdk.yml"
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	orgUser1       = "User1"
	ordererOrgName = "OrdererOrg"
	peer0org1      = "peer0.org1.example.com"
	peer0org2      = "peer0.org2.example.com"
	fabcarContract = "fabcar"
	qsccContract   = "qscc"
)

func TestConfig(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	cfg, err := sdk.Config()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	i, ok := cfg.Lookup("channels.mychannel")
	require.Equal(t, ok, true)
	require.NotNil(t, i)

}

func TestConfigGetKeyByIndex(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	cfg, err := sdk.Config()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	m, ok := cfg.Lookup("channels.mychannel")
	require.Equal(t, ok, true)
	require.NotNil(t, m)

	key := getKeyByIndex(m.(map[string]interface{}), 1)
	require.Equal(t, key, "peers")
}

func getKeyByIndex(m map[string]interface{}, index int) string {
	var j = 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	sort.Strings(keys)
	return keys[index]
}

func TestLookUpRecourseMap(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	cfg, err := sdk.Config()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	i, ok := cfg.Lookup("channels")
	require.Equal(t, ok, true)
	require.NotNil(t, i)
}

func TestCreateChainCode(t *testing.T) {
	// New fab SDK implement
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))

	// New returns a resource management client instance.
	resMgmtClient, err := resmgmt.New(clientContext)
	require.NoError(t, err)
	require.NotNil(t, resMgmtClient)

	//resMgmtClient.InstallCC()
	pkg, err := gopackager.NewCCPackage("fabcar", "/Users/shawnshen/go/")
	require.NoError(t, err)

	req := resmgmt.InstallCCRequest{
		Name:    "fabcar_1",
		Path:    "fabcar",
		Version: "1.0",
		Package: pkg,
	}

	reqPeers := resmgmt.WithTargetEndpoints(peer0org1, peer0org2)
	resp, err := resMgmtClient.InstallCC(req, reqPeers)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestCreateChannelClient(t *testing.T) {
	// New fab SDK implement
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1))
	cc, err := channel.New(ccp)
	require.NoError(t, err)
	require.NotNil(t, cc)
}

func TestQuery(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1))
	cc, err := channel.New(ccp)
	require.NoError(t, err)
	require.NotNil(t, cc)

	// new channel request for query
	req := channel.Request{
		ChaincodeID: fabcarContract,
		Fcn:         "QueryAllCars",
		Args:        nil,
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer0org1)
	resp, err := cc.Query(req, reqPeers)
	require.NoError(t, err)
	require.NotNil(t, resp)

	//fmt.Printf("resp: %v", resp)
	fmt.Printf("Querey chaincode response:\n TxId: %s\n Balance: %s\n",
		resp.TransactionID,
		resp.Payload,
	)
}

func TestQueryWithArg(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1), fabsdk.WithUser(orgUser1))
	cc, err := channel.New(ccp)
	require.NoError(t, err)
	require.NotNil(t, cc)

	// new channel request for query
	req := channel.Request{
		ChaincodeID: fabcarContract,
		Fcn:         "QueryCar",
		Args:        packArgs("CAR0"),
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer0org1)
	resp, err := cc.Query(req, reqPeers)
	require.NoError(t, err)
	require.NotNil(t, resp)

	//fmt.Printf("resp: %v", resp)
	fmt.Printf("Querey chaincode response:\n TxId: %s\n Balance: %s\n",
		resp.TransactionID,
		resp.Payload,
	)
}

func TestQueryByTxKey(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1))
	cc, err := ledger.New(ccp)
	require.NoError(t, err)
	require.NotNil(t, cc)

	// new ledger request for query
	txID := fab.TransactionID("6365315862fa66f38af56726b26b11a58e035d487b63d26b5cbb1e7e3780d7d9")

	// send request and handle response
	reqPeers := ledger.WithTargetEndpoints(peer0org1, peer0org2)

	// send request and handle response
	resp, err := cc.QueryTransaction(txID, reqPeers)
	require.NoError(t, err)
	require.NotNil(t, resp)

	//fmt.Printf("resp: %v", resp)
	fmt.Printf("Querey chaincode response:\n TxId: %s\n Balance: %s\n",
		resp,
		resp,
	)
}

func TestQueryByTxID(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1))
	cc, err := channel.New(ccp)
	require.NoError(t, err)
	require.NotNil(t, cc)

	// new channel request for query
	req := channel.Request{
		ChaincodeID: qsccContract,
		Fcn:         "GetTransactionByID",
		Args:        packArgs("mychannel", "89d7efbd1c892156bdaa057af6b9ae23d938ef3ba64264d5c177c09ac86b4484"),
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer0org1, peer0org2)
	resp, err := cc.Query(req, reqPeers)
	require.NoError(t, err)
	require.NotNil(t, resp)

	//fmt.Printf("resp: %v", resp)
	fmt.Printf("Querey chaincode response:\n TxId: %s\n Balance: %s\n",
		resp.TransactionID,
		resp.Payload,
	)
}

func TestInvoke(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1))
	cc, err := channel.New(ccp)
	require.NoError(t, err)
	require.NotNil(t, cc)

	// new channel request for query
	req := channel.Request{
		ChaincodeID: "transaction_tmp1",
		Fcn:         "ReadState",
		Args:        packArgs("5bb8d7c5096341df98e75c09731bb91a"),
	}
	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer0org1, peer0org2)
	resp, err := cc.Execute(req, reqPeers)
	require.NoError(t, err)
	require.NotNil(t, resp)

	//fmt.Printf("resp: %v", resp)
	fmt.Printf("Querey chaincode response:\n TxId: %s\n",
		resp.TransactionID,
	)
}

func TestQueryTransactionContract(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	require.NoError(t, err)
	require.NotNil(t, sdk)

	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1))
	cc, err := channel.New(ccp)
	require.NoError(t, err)
	require.NotNil(t, cc)

	// new channel request for query
	req := channel.Request{
		ChaincodeID: "fabcar",
		Fcn:         "QueryCar",
		Args:        packArgs("CAR0"),
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer0org1, peer0org2)
	resp, err := cc.Execute(req, reqPeers)
	require.NoError(t, err)
	require.NotNil(t, resp)

	ppPayload := &peer.ProposalResponsePayload{}
	err = proto.Unmarshal(resp.Responses[0].Payload, ppPayload)
	require.NoError(t, err)

	action := &peer.ChaincodeAction{}
	err = proto.Unmarshal(ppPayload.Extension, action)
	require.NoError(t, err)

	//fmt.Printf("resp: %v", resp)
	fmt.Printf("Querey chaincode response:\n TxId: %s\n",
		resp.TransactionID)

	//fabAdapter, err := NewFabricAdapter(chainID, configFile, logger.GetLogger(logger.ModuleAdapter))
	//require.Nil(t, err)
	//require.NotNil(t, fabAdapter)
	//
	//// prepare payload
	//req := channel.Request{
	//	ChaincodeID: "transaction_19",
	//	Fcn:         "ReadState",
	//	Args:        packArgs("2f095ce3c2854bd289e8438e6ec1ddcc"),
	//}
	//bz, err := json.Marshal(req)
	//require.Nil(t, err)
	//require.NotNil(t, bz)
	//eve := event.NewExecuteTransactionEvent("crossID", chainID, bz, nil)
	//
	//initAdapterConf()
	//resp, err := fabAdapter.Invoke(eve)
	//require.Nil(t, err)
	//require.NotNil(t, resp)
}

func packArgs(s ...string) [][]byte {
	res := make([][]byte, 0, 0)
	for _, str := range s {
		res = append(res, []byte(str))
	}
	return res
}
