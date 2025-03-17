#### 跨链SDK使用说明

跨链SDK是提供给业务方使用的SDK，通过该SDK，业务方用户可以构建跨链请求，发送至对应的跨链网关（主跨链网关）。

前置需求：

> 部署跨链代理，能访问代理节点的web服务
>
> 配置平行链的证书
>
> 在平行链上分别部署事物合约，业务合约

上述需求请参考跨链代理部署文档

以下是跨链SDK的使用方式：

> 调用 `SDK` 源码
```go
//用配置文件创建一个跨链SDK实例
crossSDK, err := NewCrossSDK(WithConfigFile("./config/template/cross_chain_sdk.yml"))
require.NoError(t, err)
//创建构造交易的上下文
tx1Ctx := NewCrossTxBuildCtx(
//链名字
"chain1",
//交易在跨链交易中的索引值
0,
//事务合约Excute时，要执行的业务合约信息
builder.NewContract("balance_003", "Plus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
//事务合约Rollback时，要执行的业务合约信息
builder.NewContract("balance_003", "Minus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
)
tx2Ctx := NewCrossTxBuildCtx(
//链名字
"chain2",
//交易在跨链交易中的索引值
1,
//事务合约Excute时，要执行的业务合约信息
builder.NewContract("balance_003", "GetProof", builder.NewParamsWithMap(map[string]string{"number": "1"})),
//事务合约Rollback时，要执行的业务合约信息
builder.NewContract("balance_003", "Minus", builder.NewParamsWithMap(map[string]string{"number": "1"})),
//如果业务合约中使用proof，此可选项会将proofkey写入合约参数中，chainmaker中存放在key为proofKey的参数中, fabric写在参数列表的第一个参数
builder.WithUseProofKey(true),
)
//生成跨链事件
crossEvent, err := crossSDK.GenCrossEvent(tx1Ctx, tx2Ctx)
require.NoError(t, err)

//发送跨链事件，参数syncResult代表是否同步等待跨链结果
//SendCrossEvent(event *CrossEventContext, url string, syncResult bool, opts ...EventSendOption)
ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
res, err := crossSDK.SendCrossEvent(crossEvent, "https://localhost:8080", true, WithContextOpt(ctx), WithSyncStrategyOpt(SyncStrategy{
0, 2 * time.Second, 1 * time.Second,
}))
require.NoError(t, err)
//获取crossID对应的跨链结果
res, err = crossSDK.QueryCrossResult(crossEvent.GetCrossID(), "https://localhost:8080")
require.NoError(t, err)
```

> 使用命令行工具

```shell script
cd cmd/cli
go build -o cross-chain-sdk-cli

## Deliver a CrossEvent
cross-chain-sdk-cli deliver
-c
/PathToYourProject/chainmaker-cross-chain/tools/sdk/config/template/cross_chain_sdk.yml
-u
http://localhost:8080
--params
/PathToYourProject/chainmaker-cross-chain/tools/sdk/config/template/cross_chain_params.yml

## Query Cross Result
cross-chain-sdk-cli show
-u
http://localhost:8080
--crossID
"XXXXXXX"
```