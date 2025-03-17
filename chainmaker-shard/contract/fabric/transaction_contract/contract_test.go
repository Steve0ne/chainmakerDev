/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"chainmaker.org/chainmaker-cross/utils"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-samples/chaincode/transaction_contract/go/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSmartContract_Execute(t *testing.T) {
	sc, txCtx, fn := setUp(t)
	defer fn()
	// -----------------
	// test normal logic
	// -----------------
	// new cross id
	crossID := utils.GetUUID()
	// execute params
	execute := CallContractParams{
		ContractName: fabcarContract,
		Method:  methodQuery,
		Params: nil,
	}
	executeBz, err := json.Marshal(execute)
	require.NoError(t, err)
	require.NotNil(t, executeBz)
	executeParams := string(executeBz)
	// rollback params
	rollback := CallContractParams{
		ContractName: fabcarContract,
		Method:  methodQuery,
		Params: nil,
	}
	rollbackBz, err := json.Marshal(rollback)
	require.NoError(t, err)
	require.NotNil(t, rollbackBz)
	rollbackParams := string(rollbackBz)
	// do contract method
	respStr, err := sc.Execute(txCtx, crossID, executeParams, rollbackParams)
	require.NoError(t, err)
	require.Equal(t, respStr, string(ExecuteSuccess))

	// -----------------
	// test over range
	// -----------------
	// test empty crossID
	respStr, err = sc.Execute(txCtx, EmptyCrossID, executeParams, rollbackParams)
	require.Equal(t, err, fmt.Errorf("failed to get crossID"))
	require.Equal(t, respStr, "")

	// test dup crossID
	respStr, err = sc.Execute(txCtx, crossID, executeParams, rollbackParams)
	require.Equal(t, err, fmt.Errorf("duplicated crossID: %s", crossID))
	require.Equal(t, respStr, "")

	// test empty params
	respStr, err = sc.Execute(txCtx, utils.GetUUID(), "", rollbackParams)
	require.Equal(t, err, fmt.Errorf("executeParams is nil"))
	require.Equal(t, respStr, "")
	respStr, err = sc.Execute(txCtx, utils.GetUUID(), executeParams, "")
	require.Equal(t, err, fmt.Errorf("rollbackParams is nil"))
	require.Equal(t, respStr, "")

	// test json unmarshall error
	respStr, err = sc.Execute(txCtx, utils.GetUUID(), "test", rollbackParams)
	require.NotNil(t, err)
	require.Equal(t, respStr, "")
	respStr, err = sc.Execute(txCtx, utils.GetUUID(), executeParams, "test")
	require.NotNil(t, err)
	require.Equal(t, respStr, "")

	// test invoke error
	crossID = utils.GetUUID()
	execute = CallContractParams{
		ContractName: fabcarContract,
		Method:  "error",
		Params: nil,
	}
	executeBz, err = json.Marshal(execute)
	require.NoError(t, err)
	require.NotNil(t, executeBz)
	executeParams = string(executeBz)
	respStr, err = sc.Execute(txCtx, crossID, executeParams, rollbackParams)
	require.NotNil(t, err)
	require.Equal(t, respStr, "")
}

func TestSmartContract_Commit(t *testing.T) {
	sc, txCtx, fn := setUp(t)
	defer fn()
	// -----------------
	// test normal logic
	// -----------------
	// init execute
	crossID := initExecute(t, sc, txCtx)
	// do commit
	respStr, err := sc.Commit(txCtx, crossID)
	require.NoError(t, err)
	require.Equal(t, respStr, string(CommitSuccess))

	// -----------------
	// test over range
	// -----------------
	// test commit empty cross id
	respStr, err = sc.Commit(txCtx, EmptyCrossID)
	require.Equal(t, err, fmt.Errorf("failed to get crossID"))
	require.Equal(t, respStr, "")

	// test commit ExecuteFail
	respStr, err = sc.Commit(txCtx, string(ExecuteFail))
	require.Equal(t, err, fmt.Errorf("failed to Commit cross: [%s], %s", ExecuteFail, ExecuteFail))
	require.Equal(t, respStr, "")

	// test commit CommitSuccess
	respStr, err = sc.Commit(txCtx, string(CommitSuccess))
	require.Nil(t, err)
	require.Equal(t, respStr,  string(CommitSuccess))

	// test commit CommitFail
	respStr, err = sc.Commit(txCtx, string(CommitFail))
	require.Nil(t, err)
	require.Equal(t, respStr,  string(CommitSuccess))

	// test commit Rollback
	respStr, err = sc.Commit(txCtx, string(RollbackSuccess))
	require.Equal(t, err, fmt.Errorf("failed to Commit cross: [%s], unexpected pre-state: [%s]", RollbackSuccess, RollbackSuccess))
	require.Equal(t, respStr, "")
	respStr, err = sc.Commit(txCtx, string(RollbackFail))
	require.Equal(t, err, fmt.Errorf("failed to Commit cross: [%s], unexpected pre-state: [%s]", RollbackFail, RollbackFail))
	require.Equal(t, respStr, "")
	respStr, err = sc.Commit(txCtx, string(RollbackIgnore))
	require.Equal(t, err, fmt.Errorf("failed to Commit cross: [%s], unexpected pre-state: [%s]", RollbackIgnore, RollbackIgnore))
	require.Equal(t, respStr, "")

	// test commit unknown state case
	respStr, err = sc.Commit(txCtx, string(StateUnknown))
	require.Equal(t, err, fmt.Errorf("failed to Commit cross: [%s], unknown pre-state: [%s]", StateUnknown, StateUnknown))
	require.Equal(t, respStr, "")
}

func TestSmartContract_Rollback(t *testing.T) {
	sc, txCtx, fn := setUp(t)
	defer fn()
	// -----------------
	// test normal logic
	// -----------------
	// init execute
	crossID := initExecute(t, sc, txCtx)
	// do rollback
	respStr, err := sc.Rollback(txCtx, crossID)
	require.NoError(t, err)
	require.Equal(t, respStr, string(ExecuteSuccess))

	// -----------------
	// test over range
	// -----------------
	// test rollback empty cross id
	respStr, err = sc.Rollback(txCtx, EmptyCrossID)
	require.Equal(t, err, fmt.Errorf("failed to get crossID"))
	require.Equal(t, respStr, "")

	// test rollback unknown state case
	respStr, err = sc.Rollback(txCtx, string(StateUnknown))
	require.NoError(t, err)
	require.Equal(t, respStr, "rollback ignored")
	state, err := sc.ReadState(txCtx, string(StateUnknown))
	require.NoError(t, err)
	require.Equal(t, state, string(RollbackIgnore))

	// test rollback ExecuteSuccess
	respStr, err = sc.Rollback(txCtx, string(ExecuteSuccess))
	require.NoError(t, err)
	require.Equal(t, respStr, string(ExecuteSuccess))

	// test rollback RollbackFail
	respStr, err = sc.Rollback(txCtx, string(RollbackFail))
	require.NoError(t, err)
	require.Equal(t, respStr,  string(ExecuteSuccess))

	// test commit CommitFail
	respStr, err = sc.Commit(txCtx, string(CommitFail))
	require.Nil(t, err)
	require.Equal(t, respStr,  string(CommitSuccess))

	// test commit Rollback
	respStr, err = sc.Commit(txCtx, string(RollbackSuccess))
	require.Equal(t, err, fmt.Errorf("failed to Commit cross: [%s], unexpected pre-state: [%s]", RollbackSuccess, RollbackSuccess))
	require.Equal(t, respStr, "")
	respStr, err = sc.Commit(txCtx, string(RollbackFail))
	require.Equal(t, err, fmt.Errorf("failed to Commit cross: [%s], unexpected pre-state: [%s]", RollbackFail, RollbackFail))
	require.Equal(t, respStr, "")
	respStr, err = sc.Commit(txCtx, string(RollbackIgnore))
	require.Equal(t, err, fmt.Errorf("failed to Commit cross: [%s], unexpected pre-state: [%s]", RollbackIgnore, RollbackIgnore))
	require.Equal(t, respStr, "")


}

func TestSmartContract_ReadState(t *testing.T) {

}

func setUp(t *testing.T) (*SmartContract, contractapi.TransactionContextInterface, func()) {
	dPoSStakeRuntime := NewSmartContract()
	ctrl := gomock.NewController(t)

	txSimContext := mock.NewMockTransactionContextInterface(ctrl)
	shimContext := mock.NewMockChaincodeStubInterface(ctrl)
	cache := mock.NewMemCache()

	txSimContext.EXPECT().GetStub().DoAndReturn(
		func() shim.ChaincodeStubInterface {
			return shimContext
		},
	).AnyTimes()

	shimContext.EXPECT().GetState(gomock.Any()).DoAndReturn(
		func(key string) ([]byte, error) {
			return cache.Get(key), nil
		},
	).AnyTimes()

	shimContext.EXPECT().PutState(gomock.Any(), gomock.Any()).DoAndReturn(
		func(key string, value []byte) error {
			cache.Put(key, value)
			return nil
		},
	).AnyTimes()

	shimContext.EXPECT().GetChannelID().DoAndReturn(
		func() string {
			return channelID
		},
	).AnyTimes()

	shimContext.EXPECT().InvokeChaincode(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(chaincodeName string, args [][]byte, channel string) pb.Response {
			switch chaincodeName {
			case fabcarContract:
				if string(args[0]) == "error" {
					return pb.Response{
						Status: ERROR,
						Message: "",
						Payload: []byte(ExecuteFail),
					}
				}
				return pb.Response{
					Status: SUCCESS200,
					Message: "",
					Payload: []byte(ExecuteSuccess),
				}
			default:
				return pb.Response{
					Status: 10, // 执行错误
					Message: "No Such Contract", // 未定义的合约
				}
			}
		},
	).AnyTimes()

	// init cross state
	var err error
	err = putState(txSimContext, string(StateUnknown), StateUnknown)
	require.NoError(t, err)
	err = putState(txSimContext, string(ExecuteSuccess), ExecuteSuccess)
	require.NoError(t, err)
	err = putState(txSimContext, string(ExecuteFail), ExecuteFail)
	require.NoError(t, err)
	err = putState(txSimContext, string(CommitSuccess), CommitSuccess)
	require.NoError(t, err)
	err = putState(txSimContext, string(CommitFail), CommitFail)
	require.NoError(t, err)
	err = putState(txSimContext, string(RollbackSuccess), RollbackSuccess)
	require.NoError(t, err)
	err = putState(txSimContext, string(RollbackFail), RollbackFail)
	require.NoError(t, err)
	err = putState(txSimContext, string(RollbackIgnore), RollbackIgnore)
	require.NoError(t, err)

	// put RollbackParams
	// rollback params
	rollback := CallContractParams{
		ContractName: fabcarContract,
		Method:  methodQuery,
		Params: nil,
	}
	rollbackBz, err := json.Marshal(rollback)
	require.NoError(t, err)
	require.NotNil(t, rollbackBz)
	err = putRollback(txSimContext, string(ExecuteSuccess), rollbackBz)
	require.NoError(t, err)
	err = putRollback(txSimContext, string(RollbackFail), rollbackBz)
	require.NoError(t, err)

	return dPoSStakeRuntime, txSimContext, ctrl.Finish
}


func initExecute(t *testing.T, sc *SmartContract, txCtx contractapi.TransactionContextInterface) string {
	// new cross id
	crossID := utils.GetUUID()
	// execute params
	execute := CallContractParams{
		ContractName: fabcarContract,
		Method:  methodQuery,
		Params: nil,
	}
	executeBz, err := json.Marshal(execute)
	require.NoError(t, err)
	require.NotNil(t, executeBz)
	executeParams := string(executeBz)
	// rollback params
	rollback := CallContractParams{
		ContractName: fabcarContract,
		Method:  methodQuery,
		Params: nil,
	}
	rollbackBz, err := json.Marshal(rollback)
	require.NoError(t, err)
	require.NotNil(t, rollbackBz)
	rollbackParams := string(rollbackBz)
	// do contract method
	respStr, err := sc.Execute(txCtx, crossID, executeParams, rollbackParams)
	require.NoError(t, err)
	require.Equal(t, respStr, string(ExecuteSuccess))

	return crossID
}