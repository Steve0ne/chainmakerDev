#### 跨链SDK使用说明

跨链SDK是提供给业务方使用的SDK，通过该SDK，业务方用户可以构建跨链请求，发送至对应的跨链网关（主跨链网关）。

前置需求：

> 部署跨链代理，能访问代理节点的web服务
>
> 配置平行链的证书
>
> 在平行链上分别部署事物合约，业务合约
>
> 编译 sdk 的 cli 可执行文件

上述需求请参考跨链代理部署文档

以下是跨链SDK的 CLI 使用方式：

```shell script
## Deliver a CrossEvent
cross-chain-sdk-cli deliver
-c
/PathToYourProject/chainmaker-cross-chain/tools/sdk/config/template/cross_chain_sdk.yml
-u
http://localhost:8080
--params
/PathToYourProject/chainmaker-cross-chain/tools/sdk/config/template/cross_chain_params.yml

# Return
... #  Logs
crossID:  0c1a3b099fd54162b187b9386499b9b3   # 本次跨链业务的ID


## Query Cross Result
cross-chain-sdk-cli show
-u
http://localhost:8080
--crossID
"XXXXXXX"

# Return
CrossID: 0c1a3b099fd54162b187b9386499b9b3, Code: SuccessResp, Msg: cross chain success

其中:
CrossID 为查询的跨链ID
Code 为跨链状态码, 包括: SuccessResp 表示成功, FailureResp 表示失败, ErrorResp 表示存在异常, UnknownResp 表示异常退出
Msg 为跨链事件附带的消息, 如跨链失败或异常的具体信息
```
