/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package chainmaker

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	"chainmaker.org/chainmaker/pb-go/v2/common"

	"chainmaker.org/chainmaker-cross/conf"

	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/prover"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

const (
	WaitTimeOut             int64 = 10  // 单位：秒
	RetryCount                    = 150 // 150 * 0.2 = 30s
	RetryTimePeriod               = 200 * time.Millisecond
	ProofContractParamKey         = "proofKey"
	ProofContractParamValue       = "txProof"
	ContractResultCode_OK         = 0
)

// ChainMakerAdapter adapter of chainmaker
type ChainMakerAdapter struct {
	chainID       string                   // chainID
	proofContract *conf.ProofContract      // 证据保存的合约信息
	dispatcher    *prover.ProverDispatcher // chainmaker 交易证明的证明模块分发入口
	sdk           sdk.SDKInterface         // chainmaker sdk 实例
	logger        *zap.SugaredLogger       // 日志模块
}

type result struct {
	Code    int32  `json:"code"`
	Message string `json:"message:"`
	Result  []byte `json:"result"`
}

// NewChainMakerAdapter create new instance of chainmaker adapter
func NewChainMakerAdapter(chainID, configPath string, proofContract *conf.ProofContract, logger *zap.SugaredLogger) (*ChainMakerAdapter, error) {
	chainMakerSdk, err := sdk.NewChainClient(
		sdk.WithConfPath(configPath),
	)
	if err != nil {
		return nil, err
	}
	return &ChainMakerAdapter{
		chainID:       chainID,
		proofContract: proofContract,
		dispatcher:    prover.GetProverDispatcher(),
		sdk:           chainMakerSdk,
		logger:        logger,
	}, nil
}

// GetChainID return chain id
func (c *ChainMakerAdapter) GetChainID() string {
	return c.chainID
}

// Prove prove the proof
func (c *ChainMakerAdapter) Prove(txProof *eventproto.Proof) bool {
	if txProof == nil {
		c.logger.Info("the proof is nil, not need to prove, will return true")
		return true
	}
	proveResult, err := c.dispatcher.Prove(txProof)
	if err != nil {
		c.logger.Errorf("prove the proof error, ", err)
		return false
	}
	if !proveResult {
		c.logger.Error("prove the proof failed, will return false")
	}
	return proveResult
}

// SaveProof save the proof and verify in the chain
func (c *ChainMakerAdapter) SaveProof(crossID, proofKey string, txProof *eventproto.Proof, verifyResult bool) (*eventproto.TxResponse, error) {
	c.logger.Infof("start save proof for cross[%s]->chain[%s]", crossID, txProof.GetChainID())
	if c.proofContract == nil {
		return nil, fmt.Errorf("proofContract is not configured for chain[%s]", c.chainID)
	}
	var (
		pr    prover.Prover
		exist bool
	)
	if pr, exist = c.dispatcher.GetProver(txProof.ChainId); !exist {
		return nil, fmt.Errorf("can not find prover for chain[%s]", txProof.ChainId)
	}
	verifiedProof := eventproto.NewVerifiedProof(txProof, verifyResult, fmt.Sprintf("%v", pr.GetType()), "")
	// 允许重新保存
	return c.saveProof(crossID, proofKey, verifiedProof)
}

// Invoke transfer transaction-event which include the check of transaction prove
func (c *ChainMakerAdapter) Invoke(txEvent *eventproto.TransactionEvent) (*eventproto.TxResponse, error) {
	// 直接调用sdk处理该交易
	return c.invoke(txEvent)
}

// QueryByTxKey query transaction response by txkey
func (c *ChainMakerAdapter) QueryByTxKey(txKey string) (*event.CommonTxResponse, error) {
	if len(txKey) == 0 {
		return nil, fmt.Errorf("TxKey is <nil>")
	}
	transactionInfo, err := c.sdk.GetTxByTxId(txKey)
	if err != nil {
		return nil, err
	}
	// 将TransactionInfo转换为TxResponse
	return convertToTxResponse(transactionInfo), nil
}

// QueryTx query transaction and return response
func (c *ChainMakerAdapter) QueryTx(payload []byte) (*event.CommonTxResponse, error) {
	txRequest := &common.TxRequest{}
	if err := proto.Unmarshal(payload, txRequest); err != nil {
		return nil, fmt.Errorf("unmarshal transaction payload failed, %s", err.Error())
	}
	txKey := txRequest.Payload.TxId
	c.logger.Infof("unmarshal find txKey = [%s]", txKey)
	return c.QueryByTxKey(txKey)
}

// saveProof
func (c *ChainMakerAdapter) saveProof(crossID, proofKey string, verifiedProof *eventproto.VerifiedProof) (*eventproto.TxResponse, error) {
	// 表示该交易未上链，可重新上链操作
	jsonText, err := json.Marshal(verifiedProof)
	if err != nil {
		c.logger.Errorf("marshal verified proof error crossID = [%s], ", crossID, err)
		return nil, err
	}
	chainID := verifiedProof.TxProof.GetChainID()
	params := []*common.KeyValuePair{
		{
			ProofContractParamKey,
			[]byte(proofKey),
		},
		{
			ProofContractParamValue,
			jsonText,
		},
	}
	txResponse, err := c.sdk.InvokeContract(c.proofContract.Name, c.proofContract.Method, "", params, WaitTimeOut, false)
	c.logger.Infof("send save proof for cross[%s]->chain[%s]", crossID, chainID)
	if err != nil {
		c.logger.Errorf("cross[%s]->chain[%s]'s tx-request send failed, ", crossID, chainID, err)
		return nil, err
	}
	txID := txResponse.GetTxId()
	txInfo, err := c.loadTransactionInfo(crossID, chainID, txID)
	if err != nil {
		return nil, err
	}
	return c.convertToTxResponse(crossID, txInfo)
}

// invoke transfer transaction event and return response
func (c *ChainMakerAdapter) invoke(txEvent *eventproto.TransactionEvent) (*eventproto.TxResponse, error) {
	payload := txEvent.GetPayload()
	// 将该payload转换成为ChainMaker的请求
	txRequest := &common.TxRequest{}
	if err := proto.Unmarshal(payload, txRequest); err != nil {
		return nil, fmt.Errorf("unmarshal transaction payload failed, %s", err.Error())
	}
	c.logger.Infof("cross[%s]->chain[%s]'s tx-request unmarshalled", txEvent.GetCrossID(), txEvent.GetChainID())
	txInfo, err := c.sendTxRequest(txEvent.GetCrossID(), txEvent.GetChainID(), txRequest)
	if err != nil {
		return nil, err
	}
	return c.convertToTxResponse(txEvent.GetCrossID(), txInfo)
}

// convertToTxResponse
func (c *ChainMakerAdapter) convertToTxResponse(crossID string, txInfo *common.TransactionInfo) (*eventproto.TxResponse, error) {
	contract, err := convertToContract(txInfo)
	txId, chainID := txInfo.Transaction.Payload.TxId, txInfo.Transaction.Payload.ChainId
	if err != nil {
		c.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] convert to contract failed, ", crossID, chainID, txId, err)
		return nil, fmt.Errorf("transaction[%s]'s type is not invoke user transaction, %s", txId, err.Error())
	}
	c.logger.Infof("cross[%s]->chain[%s]'s tx[%s] invoke success", crossID, chainID, txId)
	// 操作成功，则封装TxResponse
	txResp := event.NewTxResponse(chainID, txId, int64(txInfo.GetBlockHeight()), -1, contract, nil)
	return txResp, nil
}

// sendTxRequest inner interface for send tx request
func (c *ChainMakerAdapter) sendTxRequest(crossID, chainID string, txRequest *common.TxRequest) (*common.TransactionInfo, error) {
	txResponse, err := c.sdk.SendTxRequest(txRequest, WaitTimeOut, false)
	if err != nil {
		c.logger.Errorf("cross[%s]->chain[%s]'s tx-request send failed, ", crossID, chainID, err)
		return nil, err
	}
	if txResponse.Code != common.TxStatusCode_SUCCESS {
		return nil, errors.New(txResponse.Message)
	}
	//if txResponse.ContractResult == nil { //todo这里结果为nil是因为withSyncResult为false
	//	c.logger.Errorf("cross[%s]->chain[%s]'s tx-request get empty response, ", crossID, chainID, err)
	//	return nil, err
	//}
	return c.loadTransactionInfo(crossID, chainID, txResponse.GetTxId())
}

func (c ChainMakerAdapter) loadTransactionInfo(crossID, chainID, txId string) (*common.TransactionInfo, error) {
	c.logger.Infof("start get cross[%s]->chain[%s]'s tx-request[%s]'s result", crossID, chainID, txId)
	var (
		txInfo *common.TransactionInfo
		err    error
	)
	err = retry.Retry(func(uint) error {
		c.logger.Infof("cross[%s]->chain[%s]'s tx[%s] get......", crossID, chainID, txId)
		txInfo, err = c.sdk.GetTxByTxId(txId)
		if err != nil {
			return err
		}
		return nil
	},
		strategy.Limit(RetryCount),
		strategy.Wait(RetryTimePeriod), // 指定超时等待
	)
	if err != nil {
		c.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] load failed, ", crossID, chainID, txId, err)
		return nil, fmt.Errorf("cross[%s]->chain[%s]'s tx[%s] load failed, %s", crossID, chainID, txId, err.Error())
	}
	if txInfo == nil || txInfo.Transaction == nil || txInfo.Transaction.Result == nil {
		c.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] load failed, ", crossID, chainID, txId, err)
		return nil, fmt.Errorf("cross[%s]->chain[%s]'s tx[%s] load failed", crossID, chainID, txId)
	}

	if txInfo.Transaction.Result.Code != common.TxStatusCode_SUCCESS {
		c.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] invoke failed, result code = ", crossID, chainID, txId,
			txInfo.Transaction.Result.Code)
		return txInfo, errors.New(txInfo.Transaction.GetResult().ContractResult.Message)
	}
	// 交易失败或合约执行失败
	if txInfo.Transaction.GetResult().ContractResult.Code != ContractResultCode_OK {
		c.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] invoke failed, contract result code = ", crossID, chainID, txId,
			txInfo.Transaction.GetResult().ContractResult.Code)
		return txInfo, errors.New(txInfo.Transaction.GetResult().ContractResult.Message)
	}

	if txInfo.Transaction.GetResult().ContractResult.Result != nil {
		ret := &result{}
		err = json.Unmarshal(txInfo.Transaction.GetResult().ContractResult.Result, ret)
		if err != nil {
			c.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] contract result unmarshal fail: [%v]", crossID, chainID, txId, err)
			return txInfo, err
		}
		if ret.Code != 0 {
			c.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] contract result is fail:%s", crossID, chainID, txId,
				ret.Message)
			return txInfo, errors.New(ret.Message)
		}
	}

	return txInfo, nil
}

// convertToContract convert transactionInfo to ContractInfo
func convertToContract(info *common.TransactionInfo) (*eventproto.ContractInfo, error) {
	txType := info.GetTransaction().GetPayload().GetTxType()
	if txType == common.TxType_INVOKE_CONTRACT {
		//var transactPayload common.QueryPayload
		transactPayload := info.GetTransaction().GetPayload()
		//err := proto.Unmarshal(info.GetTransaction().GetRequestPayload(), &transactPayload)
		//if err != nil {
		//	return nil, err
		//}
		// 不关心version
		contract := event.NewContract(transactPayload.ContractName, "", transactPayload.Method, nil)
		txParams := transactPayload.Parameters
		for _, param := range txParams {
			innerParam := param
			contract.AddParameter(event.NewContractParameter(innerParam.Key, string(innerParam.Value)))
		}
		return contract, nil
	} else {
		return nil, errors.New("transaction is not invoke user contract")
	}
}

// convertToTxResponse
func convertToTxResponse(info *common.TransactionInfo) *event.CommonTxResponse {
	contractInfo, _ := convertToContract(info)
	txResponse := event.NewTxResponse(info.Transaction.Payload.ChainId, info.Transaction.Payload.TxId, int64(info.BlockHeight), -1, contractInfo, nil)
	if info.Transaction.Result.Code == common.TxStatusCode_SUCCESS {
		return event.NewCommonTxResponse(txResponse, event.SuccessResp, "")
	}
	return event.NewCommonTxResponse(txResponse, event.FailureResp, info.Transaction.Result.ContractResult.Message)
}
