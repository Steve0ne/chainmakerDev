/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
*/

package fabric

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"time"

	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"
	"chainmaker.org/chainmaker-cross/prover"
	"github.com/golang/protobuf/proto"
	fabcommon "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"go.uber.org/zap"
)

const (
	FabricUser     = "org_user"
	FabricPeer     = "org_peer"
	FabricProvider = "fabric"
)

const (
	WaitTimeOut     int64 = 15 // 单位：秒
	RetryCount            = 10
	RetryTimePeriod       = 2 * time.Second
	//ProofContractParamKey         = "proofKey"
	//ProofContractParamValue       = "proofValue"
	//ContractResultCode_OK         = 0
)

// FabricAdapter adapter of fabric
type FabricAdapter struct {
	chainID    		string                   	// chainID
	proofContract 	*conf.ProofContract   		// 证据保存的合约信息
	dispatcher 		*prover.ProverDispatcher 	// fabric 交易证明的证明模块分发入口
	sdk        		*fabsdk.FabricSDK        	// fabric sdk 实例
	logger     		*zap.SugaredLogger       	// 日志模块
}

// NewFabricAdapter create new instance of fabric adapter
func NewFabricAdapter(chainID, configPath string, proofContract *conf.ProofContract, logger *zap.SugaredLogger) (*FabricAdapter, error) {
	fabricSDK, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		return nil, err
	}
	return &FabricAdapter{
		chainID:    chainID,
		proofContract: proofContract,
		dispatcher: prover.GetProverDispatcher(),
		sdk:        fabricSDK,
		logger:     logger,
	}, nil
}

// SaveProof save the proof and verify in the chain
func (f *FabricAdapter) SaveProof(crossID, proofKey string, txProof *eventproto.Proof, verifyResult bool) (*eventproto.TxResponse, error) {
	f.logger.Infof("start save proof for cross[%s]->chain[%s]", crossID, txProof.GetChainID())
	var (
		pr    prover.Prover
		exist bool
	)
	if pr, exist = f.dispatcher.GetProver(f.chainID); !exist {
		return nil, fmt.Errorf("can not find prover for chain[%s]", f.chainID)
	}
	verifiedProof := eventproto.NewVerifiedProof(txProof, verifyResult, fmt.Sprintf("%v", pr.GetType()), "")
	// 允许重新保存
	return f.saveProof(crossID, proofKey, verifiedProof)
}

// saveProof
func (f *FabricAdapter) saveProof(crossID, proofKey string, verifiedProof *eventproto.VerifiedProof) (*eventproto.TxResponse, error) {
	// 表示该交易未上链，可重新上链操作
	jsonText, err := json.Marshal(verifiedProof)
	if err != nil {
		f.logger.Errorf("marshal verified proof error crossID = [%s], ", crossID, err)
		return nil, err
	}
	chainID := verifiedProof.TxProof.GetChainID()

	// get user and org peer
	user, err := conf.Config.AdapterConfigs.GetExtraConfigByKey(FabricProvider, FabricUser)
	if err != nil {
		return nil, fmt.Errorf("get adapter fabric user failed, ChainID: %s, UserKey: %s, %s", chainID, FabricUser, err)
	}
	peers, err := conf.Config.AdapterConfigs.GetExtraConfigByKey(FabricProvider, FabricPeer)
	if err != nil {
		return nil, fmt.Errorf("get adapter fabric peers failed, ChainID: %s, PeerKey: %s, %s", chainID, FabricPeer, err)
	}

	// new channel client
	ccp := f.sdk.ChannelContext(channelID, fabsdk.WithUser(user[0]))
	cc, err := channel.New(ccp)
	if err != nil {
		f.logger.Errorf("create new fabric channel error crossID = [%s], ", crossID, err)
		return nil, err
	}
	// build request
	req := channel.Request{
		ChaincodeID: f.proofContract.Name,
		Fcn:         f.proofContract.Method,
		Args:        packArgs(crossID, proofKey, string(jsonText)),
	}
	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := cc.Execute(req, reqPeers)
	if err != nil {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx-request send failed, ", crossID, chainID, err)
		return nil, err
	}
	// parse response
	if resp.Responses == nil || resp.TransactionID == "" || resp.TxValidationCode != peer.TxValidationCode_VALID {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] load failed, ", crossID, chainID, resp.TransactionID, err)
		return nil, fmt.Errorf("cross[%s]->chain[%s]'s tx[%s] load failed", crossID, chainID, resp.TransactionID)
	}
	// 交易失败或合约执行失败
	if resp.ChaincodeStatus != 200 {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] invoke failed, contract result code = ", crossID, chainID, resp.TransactionID,
			resp.ChaincodeStatus)
		return nil, errors.New(string(resp.Payload))
	}
	contract, err := executePayloadToContract(resp.Proposal.Payload)
	if err != nil {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] convert to contract failed, ", crossID, chainID, resp.TransactionID, err)
		return nil, fmt.Errorf("transaction[%s]'s type is not invoke user transaction, %s", resp.TransactionID, err.Error())
	}
	f.logger.Infof("cross[%s]->chain[%s]'s tx[%s] invoke success", crossID, chainID, resp.TransactionID)
	// 操作成功，则封装TxResponse
	txResp := event.NewTxResponse(chainID, string(resp.TransactionID), 0, -1, contract, nil)
	return txResp, nil
}

// Prove prove the proof
func (f *FabricAdapter) Prove(txProof *eventproto.Proof) bool {
	if txProof == nil {
		f.logger.Info("the proof is nil, not need to prove, will return true")
		return true
	}
	proveResult, err := f.dispatcher.Prove(txProof)
	if err != nil {
		f.logger.Errorf("prove the proof error, ", err)
		return false
	}
	if !proveResult {
		f.logger.Error("prove the proof failed, will return false")
	}
	return proveResult
}

// GetChainID return chain id, in fabric equals channel name
func (f *FabricAdapter) GetChainID() string {
	return f.chainID
}

// Invoke transfer transaction-event which include the check of transaction prove
func (f *FabricAdapter) Invoke(txEvent *eventproto.TransactionEvent) (*eventproto.TxResponse, error) {
	if txEvent.TxProof != nil {
		f.logger.Infof("cross[%s]->chain[%s]'s tx-proof need be prove", txEvent.CrossId, txEvent.ChainId)
		// 表示需要进行证明
		proveResult, err := f.dispatcher.Prove(txEvent.TxProof)
		if err != nil {
			f.logger.Errorf("cross[%s]->chain[%s]'s tx-proof prove failed, ", txEvent.CrossId, txEvent.ChainId, err)
			return nil, err
		}
		if !proveResult {
			f.logger.Errorf("cross[%s]->chain[%s]'s tx-proof prove failed, ", txEvent.CrossId, txEvent.ChainId, err)
			return nil, fmt.Errorf("cross[%s]->chain[%s]'s tx-proof prove failed, err: %s", txEvent.CrossId, txEvent.ChainId, err)
		}
		f.logger.Infof("cross[%s]->chain[%s]'s tx-proof prove success", txEvent.CrossId, txEvent.ChainId)
	} else {
		f.logger.Infof("cross[%s]->chain[%s]'s tx-proof not need be prove", txEvent.CrossId, txEvent.ChainId)
	}
	// 直接调用sdk处理该交易
	return f.invoke(txEvent)
}

// QueryByTxKey query transaction response by tx key
func (f *FabricAdapter) QueryByTxKey(txKey string) (*event.CommonTxResponse, error) {
	if len(txKey) == 0 {
		return nil, fmt.Errorf("adapter query by tx key err: length of TxKey is 0")
	}
	txID := fab.TransactionID(txKey)

	// get adapter config
	provider, err := conf.Config.AdapterConfigs.GetExtraConfigByProvider(FabricProvider)
	if err != nil {
		return nil, fmt.Errorf("get adapter provider config failed, Provider: %s, %s", FabricProvider, err.Error())
	}
	user, err := conf.Config.AdapterConfigs.GetExtraConfigByKey(FabricProvider, FabricUser)
	if err != nil {
		return nil, fmt.Errorf("get adapter user config failed, ChainID: %s, UserKey: %s, %s", provider.ChainID, FabricUser, err.Error())
	}
	peers, err := conf.Config.AdapterConfigs.GetExtraConfigByKey(FabricProvider, FabricPeer)
	if err != nil {
		return nil, fmt.Errorf("get adapter fabric peers failed, ChainID: %s, PeerKey: %s, %s", provider.ChainID, FabricPeer, err.Error())
	}

	// create ledger client
	if len(user) < 1 {
		return nil, fmt.Errorf("adapter fabric wrong config, ChainID: %s, PeerKey: %s, %s", provider.ChainID, FabricUser, fmt.Errorf("no org user set up"))
	}
	channelCtx := f.sdk.ChannelContext(provider.ChainID, fabsdk.WithUser(user[0]))
	ledgerClient, err := ledger.New(channelCtx)
	if err != nil {
		return nil, fmt.Errorf("create fabric ledger client failed, ChainID: %s, UserKey: %s, %s", provider.ChainID, FabricUser, err.Error())
	}

	// define target endpoint options
	reqPeers := ledger.WithTargetEndpoints(peers...)

	// query
	resp, err := ledgerClient.QueryTransaction(txID, reqPeers)
	if err != nil {
		return nil, err
	}

	// 将 response 转换为 contract
	contractInfo, err := payloadToContract(resp.TransactionEnvelope.Payload)
	if err != nil {
		return nil, fmt.Errorf("convert payload to contract data failed, ChainID: %s, UserKey: %s, %s", provider.ChainID, FabricUser, err.Error())
	}
	// create tx response event
	txResponse := event.NewTxResponse(provider.ChainID, string(txID), 0, -1, contractInfo, nil)
	if resp.ValidationCode == 0 {
		return event.NewCommonTxResponse(txResponse, event.SuccessResp, ""), nil
	}
	return event.NewCommonTxResponse(txResponse, event.FailureResp, ""), nil
}

// QueryTx query transaction and return response
func (f *FabricAdapter) QueryTx(payload []byte) (*event.CommonTxResponse, error) {
	txRequest := &TxRequest{}
	if err := json.Unmarshal(payload, txRequest); err != nil {
		return nil, fmt.Errorf("unmarshal txRequest failed, err: %s", err.Error())
	}
	txKey := txRequest.Header.TxId
	f.logger.Infof("unmarshal find txKey = [%s]", txKey)
	return f.QueryByTxKey(txKey)
}

// invoke transfer transaction event and return response
func (f *FabricAdapter) invoke(txEvent *eventproto.TransactionEvent) (*eventproto.TxResponse, error) {
	// get adapter config
	provider, err := conf.Config.AdapterConfigs.GetExtraConfigByProvider(FabricProvider)
	if err != nil {
		return nil, fmt.Errorf("get adapter config failed, Provider: %s, %s", FabricProvider, err.Error())
	}
	user, err := conf.Config.AdapterConfigs.GetExtraConfigByKey(FabricProvider, FabricUser)
	if err != nil {
		return nil, fmt.Errorf("get adapter fabric user failed, ChainID: %s, UserKey: %s, %s", provider.ChainID, FabricUser, err.Error())
	}

	// create channel context
	if len(user) < 1 {
		return nil, fmt.Errorf("adapter fabric wrong config, ChainID: %s, PeerKey: %s, %s", provider.ChainID, FabricUser, fmt.Errorf("no org user set up"))
	}
	channelCtx := f.sdk.ChannelContext(provider.ChainID, fabsdk.WithUser(user[0]))
	// create channel context
	fabClient, err := channel.New(channelCtx)
	if err != nil {
		return nil, fmt.Errorf("create fabric channel client failed, ChainID: %s, User: %s, %s", provider.ChainID, user, err.Error())
	}
	// create ledger client
	ledgerClient, err := ledger.New(channelCtx)
	if err != nil {
		return nil, fmt.Errorf("create fabric ledger client failed, ChainID: %s, User: %s, %s", provider.ChainID, user, err.Error())
	}

	f.logger.Infof("cross[%s]->chain[%s]'s tx-request unmarshalled", txEvent.CrossId, txEvent.ChainId)
	// 将该 payload 转换成为 TxRequest 后转发请求
	payload := txEvent.GetPayload()
	txRequest := &TxRequest{}
	if err := json.Unmarshal(payload, txRequest); err != nil {
		return nil, fmt.Errorf("unmarshal TxRequest payload failed, %s", err.Error())
	}
	f.logger.Infof("cross[%s]->chain[%s]'s tx-request unmarshalled", txEvent.CrossId, txEvent.ChainId)

	// pre execute proposal
	// 不同于 SaveProof 接口可以直接调用 SDK 的 Execute 方法，
	// 此处进行交易转发，只能使用预先构建好的 proposal 内容做预执行和确认步骤
	resp, err := SendProposal(fabClient, txRequest)
	if err != nil {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx-request send failed, ", txEvent.CrossId, txEvent.ChainId, err)
		return nil, err
	}

	// unmarshall request to channel request for commit
	req := &channel.Request{}
	err = json.Unmarshal(txRequest.Request, req)
	if err != nil {
		f.logger.Errorf("cross[%s]->chain[%s]'s parse request failed, ", txEvent.CrossId, txEvent.ChainId, err)
		return nil, err
	}

	// commit tx
	err = CommitResp(fabClient, req, resp)
	if err != nil {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] commit failed, contract result code = ", txEvent.CrossId, txEvent.ChainId, resp.TransactionID,
			resp.ChaincodeStatus)
		return nil, errors.New(string(resp.Payload))
	}

	// query tx
	txId := resp.TransactionID
	peersURLs, err := conf.Config.AdapterConfigs.GetExtraConfigByKey(FabricProvider, FabricPeer)
	if err != nil {
		return nil, fmt.Errorf("get adapter fabric user failed, ChainID: %s, UserKey: %s, %s", provider.ChainID, FabricPeer, err.Error())
	}
	blockHigh, err := RetryQueryTx(ledgerClient, txId, peersURLs)
	if err != nil {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] query failed, contract result code = ", txEvent.CrossId, txEvent.ChainId, resp.TransactionID,
			resp.ChaincodeStatus)
		return nil, fmt.Errorf("txId: [%s] query error", txId)
	}

	// parse response
	if resp.Responses == nil || resp.TransactionID == "" || resp.TxValidationCode != peer.TxValidationCode_VALID {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] load failed, ", txEvent.CrossId, txEvent.ChainId, txId, err)
		return nil, fmt.Errorf("cross[%s]->chain[%s]'s tx[%s] load failed", txEvent.CrossId, txEvent.ChainId, txId)
	}
	// 交易失败或合约执行失败
	if resp.ChaincodeStatus != 200 {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] invoke failed, contract result code = ", txEvent.CrossId, txEvent.ChainId, txId,
			resp.ChaincodeStatus)
		return nil, errors.New(string(resp.Payload))
	}
	// 检查执行结果的跨合约调用结果，解析预定义的执行结果结构
	var state = SUCCESS
	var message string
	var response Response
	for _, res := range resp.Responses {
		err := json.Unmarshal(res.ProposalResponse.Response.Payload, &response)
		if err != nil {
			f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] invoke failed, unmarshall response code error", txEvent.CrossId, txEvent.ChainId, txId,
				resp.ChaincodeStatus)
			return nil, errors.New(string(res.ProposalResponse.Response.Payload))
		}
		if response.Code == ERROR {
			state = ERROR
			message = response.Result
			break
		}
	}
	if state != SUCCESS {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] invoke failed, cross contract invoke error: [%s]", txEvent.CrossId, txEvent.ChainId, txId,
			resp.ChaincodeStatus, message)
		return nil, errors.New("cross contract invoke error" + message)
	}
	contract, err := executePayloadToContract(resp.Proposal.Payload)
	if err != nil {
		f.logger.Errorf("cross[%s]->chain[%s]'s tx[%s] convert to contract failed, ", txEvent.CrossId, txEvent.ChainId, txId, err)
		return nil, fmt.Errorf("transaction[%s]'s type is not invoke user transaction, %s", txId, err.Error())
	}
	f.logger.Infof("cross[%s]->chain[%s]'s tx[%s] invoke success", txEvent.CrossId, txEvent.ChainId, txId)
	// 操作成功，则封装TxResponse
	txResp := event.NewTxResponse(txEvent.ChainId, string(txId), int64(blockHigh), -1, contract, nil)
	return txResp, nil
}

// payloadToContract convert ProcessedResponse to ContractInfo
func payloadToContract(blob []byte) (*eventproto.ContractInfo, error) {
	// heavy unmarshal from envelope payload
	var payload fabcommon.Payload
	err := proto.Unmarshal(blob, &payload)
	if err != nil {
		return nil, err
	}

	var tx peer.Transaction
	err = proto.Unmarshal(payload.Data, &tx)
	if err != nil {
		return nil, err
	}

	if len(tx.Actions) < 1 {
		return nil, fmt.Errorf("tx has no action")
	}

	var action peer.ChaincodeActionPayload
	err = proto.Unmarshal(tx.Actions[0].Payload, &action)
	if err != nil {
		return nil, err
	}

	var proposal peer.ChaincodeProposalPayload
	err = proto.Unmarshal(action.ChaincodeProposalPayload, &proposal)
	if err != nil {
		return nil, err
	}

	var spec peer.ChaincodeInvocationSpec
	err = proto.Unmarshal(proposal.Input, &spec)
	if err != nil {
		return nil, err
	}

	// 组装 construct
	contract := event.NewContract(spec.ChaincodeSpec.ChaincodeId.GetName(), spec.ChaincodeSpec.ChaincodeId.GetVersion(), string(spec.ChaincodeSpec.Input.GetArgs()[0]), nil)
	// Input.GetArgs()[1:] 去掉第一个参数，第一个参数是合约调用方法
	txParams := spec.ChaincodeSpec.Input.GetArgs()[1:]
	for _, param := range txParams {
		contract.AddParameter(event.NewContractParameter(string(param), string(param)))
	}
	return contract, nil
}

// executePayloadToContract convert channel channel response to ContractInfo
func executePayloadToContract(blob []byte) (*eventproto.ContractInfo, error) {
	var payload peer.ChaincodeProposalPayload
	err := proto.Unmarshal(blob, &payload)
	if err != nil {
		return nil, err
	}

	var spec peer.ChaincodeInvocationSpec
	err = proto.Unmarshal(payload.Input, &spec)
	if err != nil {
		return nil, err
	}

	// convert to contract
	contract := event.NewContract(spec.ChaincodeSpec.ChaincodeId.GetName(), spec.ChaincodeSpec.ChaincodeId.GetVersion(), string(spec.ChaincodeSpec.Input.GetArgs()[0]), nil)
	txParams := spec.ChaincodeSpec.Input.GetArgs()[1:]
	for _, param := range txParams {
		p := param
		contract.AddParameter(event.NewContractParameter(string(p), string(p)))
	}
	return contract, nil
}
