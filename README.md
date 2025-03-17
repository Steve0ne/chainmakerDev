## 项目结构

| 项目                   | 版本    | 功能概括                                                                                                 |
| :------------------- | :---- | :--------------------------------------------------------------------------------------------------- |
| chainmakerGo        | 2.1.0 | 1. 可以通过内置的脚本与配置文件来生成长安链节点的执行文件。<br>2. 可以通过内置的**cmc命令行工具**在终端输入命令的方式对长安链的节点进行管理（功能丰富，但参数复杂）           |
| chainmakerCryptogen | 2.1.0 | 生成各类证书（公钥）及签名，从而辅助chainmaker-go完成节点执行文件的生成。                                                          |
| managementBackend   | 2.1.0 | 1. 可以通过图形化的前端页面来生成长安链节点的执行文件。<br>2. 可以通过图形化的前端页面对长安链的节点进行管理（功能较少，但操作简单）<br>3. 可以较为方便的查看链上的交易信息（主要用户） |
| chainmakerShard     | 2.1.0 | 完成分片事务                                                                                               |

## 身份权限管理模型

- 在长安链2.X版本中, 一共实现了三种身份权限管理模型：

1. **PermissionedWithCert**：基于数字证书的用户标识体系、基于角色的权限控制体系;
    
2. PermissionedWithKey：基于公钥的用户标识体系、基于角色的权限控制体系。
    
3. Public：基于公钥的用户标识体系、基于角色的权限控制体系。

|         |         PermissionWithCert         |
| :-----: | :--------------------------------: |
|  模式名称   |                证书模式                |
|  模式简称   |               cert模式               |
|  账户类型   | 节点账户(共识节点、同步节点、轻节点),用户账户(管理员、普通用户) |
|  账户标识.  |                数字证书                |
| 是否需要准入  |             是，证书需要CA签发             |
| 账户与组织关系 |              账户属于某个组织              |
|  共识算法   |         TBFT、RAFT、 MaxBFT          |
|  适用场景   |                联盟链                 |

### 节点类型

共识节点 consensus：有权**参与区块共识流程**的链上节点。

同步节点common：无权参与区块共识流程，但可在**链上同步数据的节点**。


### 用户类型

管理员 admin ：**可代表组织进行链上治理的用户**。
    
普通用户 client ：无权进行链上治理，但**可发送和查询交易的用户**。
    
轻节点用户 light ：无权进行链上治理，无权发送交易，**只可查询、订阅自己组织的区块、交易数据**。


## 前置软件

| 软件名称           | 版本       |
| -------------- | -------- |
| git            | -        |
| gcc            | >7.3     |
| golang         | 1.16     |
| docker         | >20.10.7 |
| docker-compose | >1.29.2  |
| make           | -        |
| tree           | -        |
| net-tools      | -        |


### 命令

``` shell
sudo apt install git
```

``` shell
sudo apt install gcc
```

``` shell
sudo snap install go --channel=1.16/stable --classic
```

``` shell
sudo apt install make
```

``` shell
sudo apt install tree
```

``` shell
apt install net-tools
```

**Docker安装请参考tools目录下的进阶使用指南.pdf**

## 简单部署-chainmakerGo

- 编译证书生成工具：

``` shell
cd chainmakerCryptogen
```

``` shell
make
```

- 将编译好的chainmakerCryptogen，软连接到chainmakerGo/tools目录：

``` shell
cd ../chainmakerGo/tools
ln -s ../../chainmakerCryptogen/ .
```

- 生成节点可执行文件:

``` shell
# 前往scripts目录
cd ../scripts
./prepare.sh 4 1
```

``` shell
# 前往scripts目录
cd ../scripts
./build_release.sh
```

- 运行节点可执行文件:

``` shell
# 前往scripts目录
cd ../scripts
./cluster_quick_start.sh normal
```

- 查看节点是否启用：

``` shell
ps -ef|grep chainmaker | grep -v grep
```

``` shell
netstat -lptn | grep 123
```

## 简单部署-management-backend

- 修改docker-compose.yml文件的版本参数：

``` shell
cd managementBackend
sudo nano docker-compose.yml
```

将**version: "3.9"** 更改为 **version: "3.3"**


- 安装镜像：

``` shell
docker-compose up
```

- 检查是否运行：

``` shell
docker ps
```

- 访问前端界面：

```
浏览器输入：localhost
账号：admin
密码：a123456
```

## cmc的编译与使用

- 编译：

``` shell
# 切换至cmc工作目录
cd ./tools/cmc
go build
```

- 证书拷贝：

``` shell
# 将./prepare.sh生成的证书拷贝至testdata子目录下
cp -rf ../../build/crypto-config ./testdata/
```

- 修改sdk配置：

``` shell
sudo nano ./testdata/sdk_config.yml
```

1. `chain_id`：链id，与bc.yml保持一致即可。
2. `org_id`: 组织id，表明该用户的组织id。
3. `user_key_file_path`、`user_key_file_path`、`user_key_file_path`、`user_key_file_path`：指向**证书配置**中的用户证书路径即可。
4. `trust_root_paths`: 信任根地址（组织证书地址），与bc.yml保持一致即可。

- cmc的功能总览：

|    功能名称    |                   功能描述                    |
| :--------: | :---------------------------------------: |
|    私钥管理    |                  私钥生成功能                   |
|    证书管理    | 包括生成ca证书、生成crl列表、生成csr、颁发证书、根据证书获取节点Id等功能 |
|  **交易功能**  |   **主要包括链管理、用户合约发布、升级、吊销、冻结、调用、查询等功能**    |
| **查询链上数据** |         **查询链上block和transaction**         |
|  **链配置**   |               **查询及更新链配置**                |
|  归档&恢复功能   |    将链上数据转移到独立存储上，归档后的数据具备可查询、可恢复到链上的特性    |
|    线上多签    |        通过系统合约实现线上多签的请求发起、投票和查询等功能         |
|  系统合约开放管理  |          管理系统合约的开放权限、查询弃用系统合约列表           |
|    证书别名    |                  证书别名管理                   |

受限于篇幅，本指南在cmc使用部分仅提及重要的功能，部分功能若未讲解，可以自行查阅：[1. 命令行工具 — chainmaker-docs v2.1.0 documentation](https://docs.chainmaker.org.cn/v2.1.0/html/dev/%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%B7%A5%E5%85%B7.html)

## 后记

关于本项目更详细地使用请参考tools目录下的**进阶使用指南.pdf**
