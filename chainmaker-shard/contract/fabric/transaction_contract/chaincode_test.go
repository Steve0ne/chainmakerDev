package main

const (
	configFile 		= "./fabric_sdk.yml"
	channelID      	= "mychannel"
	orgName        	= "Org1"
	orgAdmin       	= "Admin"
	orgUser1       	= "User1"
	ordererOrgName 	= "OrdererOrg"
	peer0org1       = "peer0.org1.example.com"
	peer0org2       = "peer0.org2.example.com"

	crossID 		= "89d7efbd1c892156bdaa057af6b9ae23d938ef3ba64264d5c177c09ac86b4488"
	fabcarContract	= "fabcar"
	methodQuery  	= "QueryAllCars"
	methodQueryWithArg 	= "QueryCar"
	transactionContract = "transaction_19"
	methodExecute 	= "Execute"
	methodCommit  	= "Commit"
	methodRollback	= "Rollback"
)

//func TestExecuteOverContract(t *testing.T) {
//	sdk, err := fabsdk.New(config.FromFile(configFile))
//	require.NoError(t, err)
//	require.NotNil(t, sdk)
//
//	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1))
//	cc, err := channel.New(ccp)
//	require.NoError(t, err)
//	require.NotNil(t, cc)
//
//	executeParams := CallContractParams{
//		ContractName: fabcarContract,
//		Method: methodQuery,
//		Params: nil,
//	}
//	executePayload, err := json.Marshal(executeParams)
//	require.NoError(t, err)
//
//	rollbackParams := CallContractParams{
//		ContractName: fabcarContract,
//		Method: methodQuery,
//		Params: nil,
//	}
//	rollbackPayload, err := json.Marshal(rollbackParams)
//	require.NoError(t, err)
//
//	// new channel request for query
//	req := channel.Request{
//		ChaincodeID: transactionContract,
//		Fcn:         methodExecute,
//		Args:        packArgs(crossID, string(executePayload), string(rollbackPayload)),
//	}
//
//	// send request and handle response
//	reqPeers := channel.WithTargetEndpoints(peer0org1, peer0org2)
//	resp, err := cc.Execute(req, reqPeers)
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(string(resp.Payload))
//}
//
//func TestExecuteOverContractWithArgs(t *testing.T) {
//	sdk, err := fabsdk.New(config.FromFile(configFile))
//	require.NoError(t, err)
//	require.NotNil(t, sdk)
//
//	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser(orgUser1))
//	cc, err := channel.New(ccp)
//	require.NoError(t, err)
//	require.NotNil(t, cc)
//
//	executeParams := CallContractParams{
//		ContractName: fabcarContract,
//		Method: methodQueryWithArg,
//		Params: []string{"CAR0"},
//	}
//	executePayload, err := json.Marshal(executeParams)
//	require.NoError(t, err)
//
//	rollbackParams := CallContractParams{
//		ContractName: fabcarContract,
//		Method: methodQueryWithArg,
//		Params: []string{"CAR0"},
//	}
//	rollbackPayload, err := json.Marshal(rollbackParams)
//	require.NoError(t, err)
//
//	// new channel request for query
//	req := channel.Request{
//		ChaincodeID: transactionContract,
//		Fcn:         methodExecute,
//		Args:        packArgs(crossID, string(executePayload), string(rollbackPayload)),
//	}
//
//	// send request and handle response
//	reqPeers := channel.WithTargetEndpoints(peer0org1, peer0org2)
//	resp, err := cc.Execute(req, reqPeers)
//	require.NoError(t, err)
//	require.NotNil(t, resp)
//
//	t.Log(string(resp.Payload))
//}
//
//func TestToArgs(t *testing.T) {
//	executeParams := CallContractParams{
//		ContractName: fabcarContract,
//		Method: methodQuery,
//		Params: nil,
//	}
//	executePayload, err := json.Marshal(executeParams)
//	require.NoError(t, err)
//	require.NotNil(t, executePayload)
//
//	rollbackParams := CallContractParams{
//		ContractName: fabcarContract,
//		Method: methodQuery,
//		Params: nil,
//	}
//	rollbackPayload, err := json.Marshal(rollbackParams)
//	require.NoError(t, err)
//	require.NotNil(t, rollbackPayload)
//
//	var execute, rollback CallContractParams
//	err = json.Unmarshal(executePayload, &execute)
//	require.NoError(t, err)
//	err = json.Unmarshal(rollbackPayload, &rollback)
//	require.NoError(t, err)
//
//	executeBytes := ToArgs(execute.Method, execute.Params)
//	rollbackBytes := ToArgs(rollback.Method, rollback.Params)
//
//	require.NotNil(t, executeBytes)
//	require.NotNil(t, rollbackBytes)
//}
//
//func packArgs(s ...string) [][]byte {
//	res := make([][]byte, 0, 0)
//	for _, str := range s {
//		res = append(res, []byte(str))
//	}
//	return res
//}

