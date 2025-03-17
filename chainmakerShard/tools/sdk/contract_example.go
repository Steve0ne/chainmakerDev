package sdk

import (
	"chainmaker.org/chainmaker-cross/sdk/builder"
	"github.com/gogo/protobuf/proto"

	//conf "chainmaker.org/chainmaker-cross/sdk/config"
	//"chainmaker.org/chainmaker-cross/utils"
	"fmt"
	"io/ioutil"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	cmSDK "chainmaker.org/chainmaker/sdk-go/v2"
	"github.com/stretchr/testify/require"
)

const (
	createContractTimeout = 5
)

const (
	certPathPrefix = "./testdata"
	tlsHostName    = "chainmaker.org"
	certPathFormat = "/crypto-config/%s/ca"
	userKeyPath    = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.key"
	userCrtPath    = certPathPrefix + "/crypto-config/%s/user/client1/client1.tls.crt"

	userSignKeyPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.sign.key"
	userSignCrtPath = certPathPrefix + "/crypto-config/%s/user/client1/client1.sign.crt"

	adminKeyPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.key"
	adminCrtPath = certPathPrefix + "/crypto-config/%s/user/admin1/admin1.tls.crt"

	orgId = "wx-org1.chainmaker.org"
)

var caPaths = []string{
	certPathPrefix + fmt.Sprintf(certPathFormat, orgId),
}

//func TestCreateTransactionContract(t *testing.T) {
//	orgId 		:= "wx-org1.chainmaker.org"
//	chainId 	:= "chain2"
//	nodeAddress := "192.168.1.33:22301"
//	// create chain-maker sdk client
//	client, err := createChainMakerClient("../testdata/sdk_config.yml", orgId, userCrtPath, userKeyPath, userSignCrtPath, userSignKeyPath)
//	require.NoError(t, err)
//	// create node
//	node := createNode(nodeAddress, 5)
//	// create admin client
//	adminClient, err := cmSDK.NewChainClient(
//		cmSDK.WithChainClientOrgId(orgId),
//		cmSDK.WithChainClientChainId(chainId),
//		cmSDK.WithChainClientLogger(getDefaultLogger()),
//		cmSDK.WithUserKeyFilePath(fmt.Sprintf(adminKeyPath, orgId)),
//		cmSDK.WithUserCrtFilePath(fmt.Sprintf(adminCrtPath, orgId)),
//		cmSDK.AddChainClientNodeConfig(node),
//	)
//	require.NoError(t, err)
//
//	fmt.Println("====================== 创建合约 ======================")
//	contractName 	:= "TransactionStable"
//	version 		:= "1.0.0"
//	byteCodePath 	:= "/Users/shawnshen/workspace/ChainMaker/chainmaker-cross-chain/contract/transaction_contract/transaction.wasm"
//	payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, common.RuntimeType_GASM, []*common.KeyValuePair{})
//	require.NoError(t, err)
//
//	// 各组织Admin权限用户签名
//	signedPayloadBytes1, err := adminClient.SignContractManagePayload(payloadBytes)
//	require.NoError(t, err)
//
//	// 收集并合并签名
//	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1})
//	require.NoError(t, err)
//
//	// 发送创建合约请求
//	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, true)
//	require.NoError(t, err)
//
//	// 检查返回数据
//	err = checkProposalRequestResp(resp)
//	require.NoError(t, err)
//}
//
//func TestCreateBalanceContract(t *testing.T) {
//	orgId 		:= "wx-org1.chainmaker.org"
//	chainId 	:= "chain2"
//	nodeAddress := "192.168.1.33:22301"
//	// create chain-maker sdk client
//	client, err := createChainMakerClient("../testdata/sdk_config.yml", orgId, userCrtPath, userKeyPath, userSignCrtPath, userSignKeyPath)
//	require.NoError(t, err)
//	// create node
//	node := createNode(nodeAddress, 5)
//	// create admin client
//	adminClient, err := cmSDK.NewChainClient(
//		cmSDK.WithChainClientOrgId(orgId),
//		cmSDK.WithChainClientChainId(chainId),
//		cmSDK.WithChainClientLogger(getDefaultLogger()),
//		cmSDK.WithUserKeyFilePath(fmt.Sprintf(adminKeyPath, orgId)),
//		cmSDK.WithUserCrtFilePath(fmt.Sprintf(adminCrtPath, orgId)),
//		cmSDK.AddChainClientNodeConfig(node),
//	)
//	require.NoError(t, err)
//
//	fmt.Println("====================== 创建合约 ======================")
//	contractName 	:= "BalanceStable"
//	version 		:= "1.0.0"
//	byteCodePath 	:= "/Users/shawnshen/workspace/ChainMaker/chainmaker-cross-chain/contract/balance/balance.wasm"
//	payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, common.RuntimeType_GASM, []*common.KeyValuePair{})
//	require.NoError(t, err)
//
//	// 各组织Admin权限用户签名
//	signedPayloadBytes1, err := adminClient.SignContractManagePayload(payloadBytes)
//	require.NoError(t, err)
//
//	// 收集并合并签名
//	mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1})
//	require.NoError(t, err)
//
//	// 发送创建合约请求
//	resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, true)
//	require.NoError(t, err)
//
//	// 检查返回数据
//	err = checkProposalRequestResp(resp)
//	require.NoError(t, err)
//}

func TestExecuteContract(t *testing.T) {
	//生成CrossSDK实例
	crossSDK, err := NewCrossSDK(WithConfigFile("/Users/shawnshen/workspace/ChainMaker/cross-chain/release/config/cross_chain_sdk.yml"))
	require.NoError(t, err)
	require.NotNil(t, crossSDK)
	//构造跨链事件
	crossEvent, err := crossSDK.GenCrossEvent(
		NewCrossTxBuildCtx(
			"chain1", 0,
			builder.NewContract("BalanceStable2", "Plus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
			builder.NewContract("BalanceStable2", "Reset", nil),
		),
		NewCrossTxBuildCtx(
			"mychannel", 1,
			builder.NewContract("fabcar", "ChangeCarOwner", builder.NewParamsNoKeys("CAR0", "OWNER")),
			builder.NewContract("fabcar", "QueryAllCars", nil),
		),
	)

	// 创建客户端
	client, err := createChainMakerClient("./testdata/sdk_config.yml", orgId, userCrtPath, userKeyPath, userSignCrtPath, userSignKeyPath)
	require.NoError(t, err)

	////获取平行链的 MultiTXBuilder2
	//configM, err := conf.InitConfigByFilepath("./config/template/cross_chain_sdk.yml")
	//require.NoError(t, err)
	//
	//mtxBuilder1, err := sdk.MultiTxFactory(parallels.ChainChainMaker, "chain1", configM)
	//require.NoError(t, err)
	//require.NotNil(t, mtxBuilder1)
	//// 可以通过获取 client 实例，调用非 SDKInterface 提供的方法，提供平行链非通用接口的入口
	//crossTx1, err := mtxBuilder1.GetClient().(*chainmaker.ChainClient).BuildTxRequestData(
	//	utils.NewUUID(),
	//	"BalanceStable",
	//	"Plus", map[string]string{"number": "1"},
	//	"Reset", map[string]string{},
	//	0,
	//)
	//require.NoError(t, err)

	crossTx1 := crossEvent.event.TxEvents.Events[0]

	tmp := common.TxRequest{}
	err = proto.Unmarshal(crossTx1.GetExecutePayload(), &tmp)
	if err != nil {
		fmt.Print(err)
	}
	resp, err := client.SendTxRequest(&tmp, 10000, true)
	require.NoError(t, err)
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}

//func TestRollbackContract(t *testing.T) {
//	client, err := createChainMakerClient("../testdata/sdk_config.yml", orgId, userCrtPath, userKeyPath, userSignCrtPath, userSignKeyPath)
//	require.NoError(t, err)
//
//	contractName := "TransactionV2"
//	method := "Rollback"
//
//	params := map[string]string{
//		"crossID": uuid.New().String(),
//	}
//
//	resp, err := client.InvokeContract(contractName, method, "", params, 10000, true)
//	require.NoError(t, err)
//	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
//}

func TestUseBalance(t *testing.T) {
	client, err := createChainMakerClient("./testdata/sdk_config.yml", orgId, userCrtPath, userKeyPath, userSignCrtPath, userSignKeyPath)
	require.Nil(t, err)

	contractName := "BalanceStable2"
	method := "Reset"

	//params := map[string]string{"number": "100"}
	params := []*common.KeyValuePair{
		{
			Key:   "number",
			Value: []byte("100"),
		},
	}

	resp, err := client.InvokeContract(contractName, method, "", params, 10000, true)
	require.Nil(t, err)
	fmt.Printf("QUERY counter-go contract resp: %+v\n", resp)
}

func createChainMakerClient(configPath, orgId, userCrtPath, userKeyPath, userSignCrtPath, userSignKeyPath string) (*cmSDK.ChainClient, error) {
	userCrtBytes, err := ioutil.ReadFile(fmt.Sprintf(userCrtPath, orgId))
	if err != nil {
		return nil, err
	}
	userKeyBytes, err := ioutil.ReadFile(fmt.Sprintf(userKeyPath, orgId))
	if err != nil {
		return nil, err
	}
	userSignCrtBytes, err := ioutil.ReadFile(fmt.Sprintf(userSignCrtPath, orgId))
	if err != nil {
		return nil, err
	}
	userSignKeyBytes, err := ioutil.ReadFile(fmt.Sprintf(userSignKeyPath, orgId))
	if err != nil {
		return nil, err
	}
	chainClient, err := cmSDK.NewChainClient(
		cmSDK.WithConfPath(configPath),
		cmSDK.WithUserCrtBytes(userCrtBytes),
		cmSDK.WithUserKeyBytes(userKeyBytes),
		cmSDK.WithUserSignKeyBytes(userSignKeyBytes),
		cmSDK.WithUserSignCrtBytes(userSignCrtBytes),
	)

	if err != nil {
		return nil, err
	}

	return chainClient, nil
}

//func createNode(nodeAddr string, connCnt int) *cmSDK.NodeConfig {
//	node := cmSDK.NewNodeConfig(
//		// 节点地址，格式：127.0.0.1:12301
//		cmSDK.WithNodeAddr(nodeAddr),
//		// 节点连接数
//		cmSDK.WithNodeConnCnt(connCnt),
//		// 节点是否启用TLS认证
//		cmSDK.WithNodeUseTLS(true),
//		// 根证书路径，支持多个
//		cmSDK.WithNodeCAPaths(caPaths),
//		// TLS Hostname
//		cmSDK.WithNodeTLSHostName(tlsHostName),
//	)
//
//	return node
//}
//
//func getDefaultLogger() *zap.SugaredLogger {
//	config := log.LogConfig{
//		Module:       "[SDK]",
//		LogPath:      "./sdk.log",
//		LogLevel:     log.LEVEL_DEBUG,
//		MaxAge:       30,
//		JsonFormat:   false,
//		ShowLine:     true,
//		LogInConsole: true,
//	}
//	logger, _ := log.InitSugarLogger(&config)
//	return logger
//}
//
//func checkProposalRequestResp(resp *common.TxResponse) error {
//	if resp.Code != common.TxStatusCode_SUCCESS {
//		return errors.New(resp.Message)
//	}
//
//	if resp.ContractResult == nil {
//		return fmt.Errorf("contract result is nil")
//	}
//
//	if resp.ContractResult.Code != common.ContractResultCode_OK {
//		return errors.New(resp.ContractResult.Message)
//	}
//
//	return nil
//}
//
//func TestChannel(t *testing.T) {
//	//ctx := context.WithCancel(context.Background())
//	x := make(chan int, 1)
//	go func() {
//		select {
//		case y := <- x:
//			t.Log("close channel", y)
//		}
//	}()
//	time.Sleep(time.Second * 2)
//	close(x)
//}
//
