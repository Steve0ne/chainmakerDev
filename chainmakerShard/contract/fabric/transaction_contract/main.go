/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

func NewSmartContract() *SmartContract {
	return &SmartContract{
		contractapi.Contract{},
	}
}

func (s *SmartContract) Execute(ctx contractapi.TransactionContextInterface, crossID, executeParams, rollbackParams string) (string, error) {
	var err error
	// check and parse params
	if crossID == EmptyCrossID {
		// will end contract calling
		return "", fmt.Errorf("failed to get crossID")
	}
	if isCrossIDExist(ctx, crossID) {
		return "", fmt.Errorf("duplicated crossID: " + crossID)
	}
	// check execute params
	if executeParams == EmptyParam {
		return "", fmt.Errorf("executeParams is nil")
	}
	// unmarshall execute params
	var execute CallContractParams
	err = json.Unmarshal([]byte(executeParams), &execute)
	if err != nil {
		return "", err
	}
	// check rollback params
	if rollbackParams == EmptyParam {
		return "", fmt.Errorf("rollbackParams is nil")
	}
	// unmarshall rollback params
	var rollback CallContractParams
	err = json.Unmarshal([]byte(rollbackParams), &rollback)
	if err != nil {
		return "", err
	}

	// put data
	err = putExecute(ctx, crossID, []byte(executeParams))
	if err != nil {
		return "", err
	}
	err = putRollback(ctx, crossID, []byte(rollbackParams))
	if err != nil {
		return "", err
	}
	// call execute method
	resp := ctx.GetStub().InvokeChaincode(execute.ContractName, ToArgs(execute.Method, execute.Params), ctx.GetStub().GetChannelID())
	if resp.Status != SUCCESS200 {
		// 返回失败结果
		err = putState(ctx, crossID, ExecuteFail)
		if err != nil {
			return "", err
		}
		res := &Response{
		 	Code: int(ERROR),
			Result: resp.Message,
		}
		bz, err := json.Marshal(res)
		if err != nil {
			return "", err
		}
		return string(bz), nil
	} else {
		// 返回正确结果
		err = putState(ctx, crossID, ExecuteSuccess)
		if err != nil {
			return "", err
		}
		res := &Response{
			Code: int(SUCCESS),
			Result: string(resp.Payload),
		}
		bz, err := json.Marshal(res)
		if err != nil {
			return "", err
		}
		return string(bz), nil
	}
}

func (s *SmartContract) Commit(ctx contractapi.TransactionContextInterface, crossID string) (string, error) {
	// check crossID
	if crossID == EmptyCrossID {
		// will end contract calling
		return "", fmt.Errorf("failed to get crossID")
	} else {
		// check state
		crossState, err := getState(ctx, crossID)
		if err != nil {
			return "", err
		}
		if crossState == ExecuteSuccess {
			// change state
			err = putState(ctx, crossID, CommitSuccess)
			if err != nil {
				return "", err
			}
			res := &Response{
				Code: int(SUCCESS),
				Result: string(CommitSuccess),
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return "", err
			}
			return string(bz), nil
		} else if crossState == ExecuteFail {
			return "", fmt.Errorf("failed to Commit cross: [%s], ExecuteFail", crossID)
		} else if crossState == CommitSuccess {
			return string(CommitSuccess), nil
		} else if crossState == CommitFail {
			// change state
			err = putState(ctx, crossID, CommitSuccess)
			if err != nil {
				return "", err
			}
			res := &Response{
				Code: int(SUCCESS),
				Result: string(CommitSuccess),
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return "", err
			}
			return string(bz), nil
		} else if crossState == RollbackSuccess || crossState == RollbackFail || crossState == RollbackIgnore {
			return "", fmt.Errorf("failed to Commit cross: [%s], unexpected pre-state: [%s]", crossID, crossState)
		}
		return "", fmt.Errorf("failed to Commit cross: [%s], unknown pre-state: [%s]", crossID, crossState)
	}
}

func (s *SmartContract) Rollback(ctx contractapi.TransactionContextInterface, crossID string) (string, error) {
	// check crossID
	if crossID == EmptyCrossID {
		// 返回结果
		return "", fmt.Errorf("failed to get crossID")
	} else {
		// check state
		crossState, err := getState(ctx, crossID)
		if err != nil {
			return "", err
		}
		// check rollback state
		if crossState == StateUnknown {
			// 返回结果
			err = putState(ctx, crossID, RollbackIgnore)
			if err != nil {
				return "", err
			}
			res := &Response{
				Code: int(SUCCESS),
				Result: string(RollbackIgnore),
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return "", err
			}
			return string(bz), nil
		}
		// check pre-state
		if crossState == ExecuteSuccess || crossState == RollbackFail {
			if cp, err := getRollback(ctx, crossID); err != nil {
				return "", fmt.Errorf("failed to get Rollback Data, crossID: " + crossID)
			} else {
				resp := ctx.GetStub().InvokeChaincode(cp.ContractName, ToArgs(cp.Method, cp.Params), ctx.GetStub().GetChannelID())
				if resp.Status != SUCCESS200 {
					err = putState(ctx, crossID, RollbackFail)
					if err != nil {
						return "", err
					}
					// 返回失败结果
					res := &Response{
						Code: int(ERROR),
						Result: resp.Message,
					}
					bz, err := json.Marshal(res)
					if err != nil {
						return "", err
					}
					return string(bz), nil
				} else {
					// 返回结果
					err = putState(ctx, crossID, RollbackSuccess)
					if err != nil {
						return "", err
					}
					res := &Response{
						Code: int(SUCCESS),
						Result: resp.Message,
					}
					bz, err := json.Marshal(res)
					if err != nil {
						return "", err
					}
					return string(bz), nil
				}
			}
		} else if crossState == ExecuteFail || crossState == RollbackIgnore {
			// 返回结果
			err = putState(ctx, crossID, RollbackIgnore)
			if err != nil {
				return "", err
			}
			res := &Response{
				Code: int(SUCCESS),
				Result: string(RollbackIgnore),
			}
			bz, err := json.Marshal(res)
			if err != nil {
				return "", err
			}
			return string(bz), nil
		} else if crossState == CommitSuccess || crossState == CommitFail {
			// 返回结果
			return string(RollbackIgnore), nil
		} else if crossState == RollbackSuccess {
			return string(RollbackSuccess), nil
		}
		return "", fmt.Errorf(string("failed to Rollback, unexpected state: " + crossState))
	}
}

func (s *SmartContract) ReadState(ctx contractapi.TransactionContextInterface, crossID string) (string, error) {
	// check crossID
	if crossID == "" {
		// 返回结果
		return "", fmt.Errorf("failed to get crossID")
	} else {
		state, err := getState(ctx, crossID)
		if err != nil {
			return "", err
		}
		res := &Response{
			Code: int(SUCCESS),
			Result: string(state),
		}
		bz, err := json.Marshal(res)
		if err != nil {
			return "", err
		}
		return string(bz), nil
	}
}

func (s *SmartContract) SaveProof(ctx contractapi.TransactionContextInterface, crossID, proofKey, txProof string) (string, error) {
	// 检测是否已经存储proof 是则返回存储的proof， 否则存储
	ret, err := getProof(ctx, crossID + "/" + proofKey)
	if err != nil {
		return "", err
	}
	if len(ret) > 0 {
		// Proof 已存在，返回历史数据
		return ret, nil
	}
	// 写入Proof
	err = putProof(ctx, crossID + "/" + proofKey, txProof)
	if err != nil {
		return "", err
	}
	// 返回状态
	res := &Response{
		Code: int(SUCCESS),
		Result: "ProofPutSuccess",
	}
	bz, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(bz), nil
}

func (s *SmartContract) ReadProof(ctx contractapi.TransactionContextInterface, crossID, proofKey string) (string, error) {
	// 检测是否已经存储proof 是则返回存储的proof， 否则存储
	ret, err := getProof(ctx, crossID + "/" + proofKey)
	if err != nil {
		return "", err
	}
	// 返回历史数据
	res := &Response{
		Code: int(SUCCESS),
		Result: ret,
	}
	bz, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	return string(bz), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}