/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package main

// 安装合约时会执行此方法，必须
//export init_contract
func initContract() {}

// 升级合约时会执行此方法，必须
//export upgrade
func upgrade() {}

//export Execute
func Execute() {
	// check and parse params
	crossID, executeParams, rollbackParams := UnpackUploadParams(Args())
	if crossID == EmptyCrossID {
		// will end contract calling
		ErrorResult("failed to get crossID")
		return
	} else {
		if isCrossIDExist(crossID) {
			ErrorResult("duplicated crossID: " + crossID)
			return
		}
		// check execute params
		if executeParams == nil {
			ErrorResult("executeParams is nil")
			return
		}
		eMap := callParamsToMap(executeParams)
		if eMap == nil {
			ErrorResult("failed to parse execute params")
			return
		}
		// check execute params
		if rollbackParams == nil {
			ErrorResult("rollbackParams is nil")
			return
		}
		rMap := callParamsToMap(rollbackParams)
		if rMap == nil {
			ErrorResult("failed to parse rollback params")
			return
		}
		// put data
		putExecute(crossID, eMap)
		putRollback(crossID, rMap)
	}

	// call execute method
	if bz, resultCode := CallContract(executeParams.ContractName, executeParams.Method, executeParams.Params); resultCode != SUCCESS {
		// 返回失败结果
		putState(crossID, ExecuteFail)
		resp := &Response{
			Code: int(ERROR),
			Result: string(bz),
		}
		respStr := ResponseToJsonString(resp)
		SuccessResult(respStr)
		return
	} else {
		// 返回正确结果
		putState(crossID, ExecuteSuccess)
		resp := &Response{
			Code: int(SUCCESS),
			Result: string(bz),
		}
		respStr := ResponseToJsonString(resp)
		SuccessResult(respStr)
		return
	}
}

//export Commit
func Commit() {
	// get crossID
	if crossID, resultCode := Arg("crossID"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get crossID")
		return
	} else {
		// check state
		crossState := getState(string(crossID))
		if crossState == ExecuteSuccess {
			// change state
			putState(string(crossID), CommitSuccess)
			resp := &Response{
				Code: int(SUCCESS),
				Result: string(CommitSuccess),
			}
			respStr := ResponseToJsonString(resp)
			SuccessResult(respStr)
			return
		} else if crossState == ExecuteFail {
			ErrorResult(string("failed to Commit, unexpected pre-state: " + crossState))
			return
		} else if crossState == CommitSuccess {
			resp := &Response{
				Code: int(SUCCESS),
				Result: string(CommitSuccess),
			}
			respStr := ResponseToJsonString(resp)
			SuccessResult(respStr)
			return
		} else if crossState == CommitFail {
			// change state
			putState(string(crossID), CommitSuccess)
			resp := &Response{
				Code: int(SUCCESS),
				Result: string(CommitSuccess),
			}
			respStr := ResponseToJsonString(resp)
			SuccessResult(respStr)
			return
		} else if crossState == RollbackSuccess || crossState == RollbackFail || crossState == RollbackIgnore {
			ErrorResult(string("failed to Commit, unexpected pre-state: " + crossState))
			return
		}
		ErrorResult(string("failed to Commit, unexpected pre-state: " + crossState))
		return
	}
}

//export Rollback
func Rollback() {
	// get crossID
	if crossID, resultCode := Arg("crossID"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get crossID")
		return
	} else {
		// check state
		crossState := getState(string(crossID))
		// check rollback state
		if crossState == StateUnknown {
			// 返回结果
			putState(string(crossID), RollbackIgnore)
			resp := &Response{
				Code: int(SUCCESS),
				Result: string(RollbackIgnore),
			}
			respStr := ResponseToJsonString(resp)
			SuccessResult(respStr)
			return
		}
		// check pre-state
		if crossState == ExecuteSuccess || crossState == RollbackFail {
			if m, resultCode := getRollback(string(crossID)); resultCode != SUCCESS {
				ErrorResult("failed to get Rollback Data, crossID: " + string(crossID))
				return
			} else {
				cp := callParamsFromMap(m)
				if bz, resultCode := CallContract(cp.ContractName, cp.Method, cp.Params); resultCode != SUCCESS {
					// 返回失败结果
					putState(string(crossID), RollbackFail)
					resp := &Response{
						Code: int(ERROR),
						Result: string(bz),
					}
					respStr := ResponseToJsonString(resp)
					SuccessResult(respStr)
					return
				} else {
					// 返回结果
					putState(string(crossID), RollbackSuccess)
					resp := &Response{
						Code: int(SUCCESS),
						Result: string(bz),
					}
					respStr := ResponseToJsonString(resp)
					SuccessResult(respStr)
					return
				}
			}
		} else if crossState == ExecuteFail || crossState == RollbackIgnore {
			// 返回结果
			putState(string(crossID), RollbackIgnore)
			resp := &Response{
				Code: int(SUCCESS),
				Result: string(RollbackIgnore),
			}
			respStr := ResponseToJsonString(resp)
			SuccessResult(respStr)
			return
		} else if crossState == CommitSuccess || crossState == CommitFail {
			// 返回结果
			resp := &Response{
				Code: int(SUCCESS),
				Result: string(RollbackIgnore),
			}
			respStr := ResponseToJsonString(resp)
			SuccessResult(respStr)
			return
		} else if crossState == RollbackSuccess {
			resp := &Response{
				Code: int(SUCCESS),
				Result: string(RollbackSuccess),
			}
			respStr := ResponseToJsonString(resp)
			SuccessResult(respStr)
			return
		}
		ErrorResult(string("failed to Rollback, unexpected state: " + crossState))
		return
	}
}

//export ReadState
func ReadState() {
	// get crossID
	if crossID, resultCode := Arg("crossID"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get crossID")
		return
	} else {
		state := getState(string(crossID))
		resp := &Response{
			Code: int(SUCCESS),
			Result: string(state),
		}
		respStr := ResponseToJsonString(resp)
		SuccessResult(respStr)
		return
	}
}

//export SaveProof
func SaveProof() {
	var crossID, proofKey, txProof []byte
	var resultCode ResultCode
	// get crossID
	if crossID, resultCode = Arg("crossID"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get crossID")
		return
	}
	// get proofKey
	if proofKey, resultCode = Arg("proofKey"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get proofKey")
		return
	}
	// get txProof
	if txProof, resultCode = Arg("txProof"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get txProof")
		return
	}
	// 检测是否已经存储proof 是则返回存储的proof， 否则存储
	ret, resultCode := getProof(string(crossID) + "." + string(proofKey))
	if resultCode != SUCCESS {
		// 写入Proof
		resultCode = putProof(string(crossID) + "." + string(proofKey), string(txProof))
		if resultCode != SUCCESS {
			ErrorResult("failed to putProof, crossID: " + string(crossID) + "proofKey: " + string(proofKey))
			return
		}
		// 返回状态
		resp := &Response{
			Code: int(SUCCESS),
			Result: "ProofPutSuccess",
		}
		respStr := ResponseToJsonString(resp)
		SuccessResult(respStr)
		return
	} else {
		// Proof 已存在，返回历史数据
		resp := &Response{
			Code: int(SUCCESS),
			Result: ret,
		}
		respStr := ResponseToJsonString(resp)
		SuccessResult(respStr)
		return
	}
}

//export ReadProof
func ReadProof() {
	var crossID, proofKey []byte
	var resultCode ResultCode
	// get crossID
	if crossID, resultCode = Arg("crossID"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get crossID")
		return
	}
	// get proofKey
	if proofKey, resultCode = Arg("proofKey"); resultCode != SUCCESS {
		// 返回结果
		ErrorResult("failed to get proofKey")
		return
	}
	// 读取 Proof
	ret, resultCode := getProof(string(crossID) + "." + string(proofKey))
	if resultCode != SUCCESS {
		ErrorResult("failed to call ReadProof, crossID: " + string(crossID) + "proofKey: " + string(proofKey))
		return
	}
	// 返回数据
	resp := &Response{
		Code: int(SUCCESS),
		Result: ret,
	}
	respStr := ResponseToJsonString(resp)
	SuccessResult(respStr)
	return
}

func main() {}
