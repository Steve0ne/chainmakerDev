/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package transaction

import (
	"fmt"
	"testing"
	"time"

	eventproto "chainmaker.org/chainmaker-cross/pb/protogo/event"

	"chainmaker.org/chainmaker-cross/adapter"
	"chainmaker.org/chainmaker-cross/conf"
	"chainmaker.org/chainmaker-cross/event"
	"chainmaker.org/chainmaker-cross/handler"
	"chainmaker.org/chainmaker-cross/listener/inner_listener"
	"chainmaker.org/chainmaker-cross/logger"
	"chainmaker.org/chainmaker-cross/prover"
	"chainmaker.org/chainmaker-cross/prover/impl"
	"chainmaker.org/chainmaker-cross/router"
	"chainmaker.org/chainmaker-cross/store"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	chain1 = "chain1"
	chain2 = "chain2"
)

var proverChange = false

func TestInitManager(t *testing.T) {
	event.InitLog(getLogger())
	conf.Config.StorageConfig = &conf.StorageConfig{
		Provider: "memory",
	}
	stateDB := store.InitStateDB()
	defer stateDB.Close()
	require.NotNil(t, stateDB)
	manager := InitManager(stateDB)
	require.NotNil(t, manager)
}

func TestTransactionManager(t *testing.T) {
	event.InitLog(getLogger())
	conf.Config.StorageConfig = &conf.StorageConfig{
		Provider: "memory",
	}
	stateDB := store.InitStateDB()
	defer stateDB.Close()
	require.NotNil(t, stateDB)
	manager := InitManager(stateDB)
	require.NotNil(t, manager)
	handler.InitEventHandlers(stateDB, manager.GetEventChan())
	innerListener := inner_listener.NewInnerListener()
	err := innerListener.ListenStart()
	require.Nil(t, err)
	err = routerDispatcher()
	require.Nil(t, err)
	pd := prover.GetProverDispatcher()
	pd.Register(NewProverMock([]string{chain1, chain2}))
	adapterDispatcher := adapter.GetChainAdapterDispatcher()
	adapterDispatcher.SetLog(getLogger())
	adapterDispatcher.Register(NewChainAdapterMock(chain1, chainInvoke, chain1QueryByTxKey, nil))
	adapterDispatcher.Register(NewChainAdapterMock(chain2, chainInvoke, chain2QueryByTxKey, nil))
	eventChan := manager.GetEventChan()
	require.NotNil(t, eventChan)
	err = manager.Start()
	defer manager.Stop()
	require.Nil(t, err)
	// 创建一个跨链交易
	crossTxs := make([]*eventproto.CrossTx, 0)
	crossTxs = append(crossTxs, initCrossTxs(chain1, 0))
	crossTxs = append(crossTxs, initCrossTxs(chain2, 1))
	crossEvent := event.NewCrossEvent(crossTxs)
	result, err := handleCrossEvent(manager, crossEvent)
	require.Nil(t, err)
	fmt.Println(result)
	fmt.Println("--- success cross event handle over ---")
	fmt.Println("--- start handle unsupport chain ---")
	crossTxs = make([]*eventproto.CrossTx, 0)
	crossTxs = append(crossTxs, initCrossTxs(chain1, 0))
	crossTxs = append(crossTxs, initCrossTxs("chain3", 1))
	crossEvent = event.NewCrossEvent(crossTxs)
	result, err = handleCrossEvent(manager, crossEvent)
	require.Nil(t, err)
	fmt.Println(result)
	crossTxs = make([]*eventproto.CrossTx, 0)
	crossTxs = append(crossTxs, initCrossTxs("chain4", 0))
	crossTxs = append(crossTxs, initCrossTxs("chain3", 1))
	crossEvent = event.NewCrossEvent(crossTxs)
	result, err = handleCrossEvent(manager, crossEvent)
	require.Nil(t, err)
	fmt.Println(result)
	fmt.Println("--- unsupport chain cross event handle over ---")
	fmt.Println("--- start handle prover error ---")
	crossTxs = make([]*eventproto.CrossTx, 0)
	crossTxs = append(crossTxs, initCrossTxs(chain1, 0))
	crossTxs = append(crossTxs, initCrossTxs(chain2, 1))
	crossEvent = event.NewCrossEvent(crossTxs)
	proverChange = true
	result, err = handleCrossEvent(manager, crossEvent)
	require.Nil(t, err)
	fmt.Println(result)
}

func handleCrossEvent(manager *Manager, crossEvent *eventproto.CrossEvent) (interface{}, error) {
	manager.handle(crossEvent, true)
	time.Sleep(time.Second * 3) // 确保处理完成
	// 查询对应结果
	eveHandler, _ := handler.GetEventHandlerTools().GetHandler(handler.CrossSearch)
	return eveHandler.Handle(event.NewCrossSearchEvent(crossEvent.GetCrossID()), true)
}

func initCrossTxs(chainID string, index int32) *eventproto.CrossTx {
	return event.NewCrossTx(chainID, index, nil, nil, nil)
}

func chainInvoke(eve *eventproto.TransactionEvent) (*eventproto.TxResponse, error) {
	response := event.NewTxResponse(eve.GetChainID(), "", 1024, 0, nil, nil)
	return response, nil
}

func chain1QueryByTxKey(txKey string) (*event.CommonTxResponse, error) {
	response := event.NewTxResponse(chain1, txKey, 1024, 0, nil, nil)
	return event.NewCommonTxResponse(response, event.SuccessResp, ""), nil
}

func chain2QueryByTxKey(txKey string) (*event.CommonTxResponse, error) {
	response := event.NewTxResponse(chain2, txKey, 1024, 0, nil, nil)
	return event.NewCommonTxResponse(response, event.SuccessResp, ""), nil
}

func routerDispatcher() error {
	routerDispatcher := router.GetDispatcher()
	testLogger := getLogger()
	event.InitLog(testLogger)
	routerDispatcher.SetLogger(testLogger)
	// 注册innerrouter
	var chainIDs = []string{"chain1", "chain2"}
	innerRouter := router.GetInnerRouter()
	innerRouter.Init(chainIDs)
	return routerDispatcher.Register(innerRouter)
}

func getLogger() *zap.SugaredLogger {
	config := []*logger.LogModuleConfig{
		{
			ModuleName:   "default",
			LogLevel:     logger.INFO,
			FilePath:     "logs/default.log",
			MaxAge:       365,
			RotationTime: 1,
			LogInConsole: false,
			ShowColor:    true,
		},
	}
	logger.InitLogConfig(config)
	return logger.GetLogger("default")
}

type ChainAdapterMock struct {
	chainID      string
	invoke       func(*eventproto.TransactionEvent) (*eventproto.TxResponse, error)
	queryByTxKey func(string) (*event.CommonTxResponse, error)
	queryTx      func([]byte) (*event.CommonTxResponse, error)
}

func NewChainAdapterMock(chainID string, invokeFunc func(*eventproto.TransactionEvent) (*eventproto.TxResponse, error),
	queryByTxKey func(string) (*event.CommonTxResponse, error),
	queryTx func([]byte) (*event.CommonTxResponse, error)) *ChainAdapterMock {
	return &ChainAdapterMock{
		chainID:      chainID,
		invoke:       invokeFunc,
		queryByTxKey: queryByTxKey,
		queryTx:      queryTx,
	}
}

func (c *ChainAdapterMock) GetChainID() string {
	return c.chainID
}

func (c *ChainAdapterMock) Invoke(txEvent *eventproto.TransactionEvent) (*eventproto.TxResponse, error) {
	if c.invoke != nil {
		return c.invoke(txEvent)
	}
	return nil, nil
}

func (c *ChainAdapterMock) QueryByTxKey(txKey string) (*event.CommonTxResponse, error) {
	if c.queryByTxKey != nil {
		return c.queryByTxKey(txKey)
	}
	return nil, nil
}

func (c *ChainAdapterMock) QueryTx(payload []byte) (*event.CommonTxResponse, error) {
	if c.queryTx != nil {
		return c.queryTx(payload)
	}
	return nil, nil
}

type ProverMock struct {
	chainIDs []string
}

func NewProverMock(chainIDs []string) *ProverMock {
	return &ProverMock{
		chainIDs: chainIDs,
	}
}

func (p *ProverMock) GetType() impl.ProverType {
	return impl.TrustProverType
}

func (p *ProverMock) GetChainIDs() []string {
	return p.chainIDs
}

func (p *ProverMock) ToProof(chainID, txKey string, blockHeight int64, index int32, contract *eventproto.ContractInfo, extra []byte) (*eventproto.Proof, error) {
	return event.NewProof(chainID, txKey, blockHeight, index, contract, extra), nil
}

func (p *ProverMock) Prove(proof *eventproto.Proof) (bool, error) {
	if proverChange {
		return false, nil
	}
	return true, nil
}
