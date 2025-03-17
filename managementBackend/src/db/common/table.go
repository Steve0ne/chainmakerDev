/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package common

const (
	TableBlock          = "cmb_block"
	TableChain          = "cmb_chain"
	TableTransaction    = "cmb_transaction"
	TableCert           = "cmb_cert"
	TableNode           = "cmb_node"
	TableOrg            = "cmb_org"
	TableContract       = "cmb_contract"
	TableInvokeRecords  = "cmb_invoke_records"
	TableChainPolicy    = "cmb_chain_policy"
	TableChainPolicyOrg = "cmb_chain_policy_org"
	TableChainOrg       = "cmb_chain_org"
	TableChainOrgNode   = "cmb_chain_org_node"
	TableOrgNode        = "cmb_org_node"
	TableUploadContent  = "cmb_upload_content"
	TableUser           = "cmb_user"
	TableVoteManagement = "cmb_vote_management"
	TableChainConfig    = "cmb_chain_config"
	TableChainErrorLog  = "cmb_chain_error_log"
)

// package chain
type Block struct {
	CommonIntField
	PreBlockHash  string `json:"preBlockHash"`
	ChainId       string `gorm:"unique_index:chain_id_block_height_index;index:chain_id_block_hash_index" json:"chainId"`
	BlockHash     string `gorm:"index:chain_id_block_hash_index" json:"blockHash"`            //区块哈希
	BlockHeight   uint64 `gorm:"unique_index:chain_id_block_height_index" json:"blockHeight"` //区块高度
	OrgId         string // 组织ID
	Timestamp     int64  `json:"timestamp"`
	DagHash       string `json:"dagHash"`                       //DAG哈希
	TxCount       int    `json:"txCount"`                       //交易数量
	RwSetHash     string `json:"rwSetHash"`                     //读写集生成merkle树根哈希
	TxRootHash    string `json:"txRootHash"`                    //交易merkle树根哈希
	Proposer      string `gorm:"type:longtext" json:"proposer"` //提案节点标识
	ProposerType  string `json:"proposer_type"`
	ConsensusArgs string `gorm:"type:longtext" json:"consensusArgs"` //共识参数
	ProposerId    string `json:"proposer_id"`                        // 打包节点
}

func (*Block) TableName() string {
	return TableBlock
}

type Chain struct {
	CommonIntField
	ChainId           string `gorm:"unique_index:chain_id_index"` //链ID，链唯一标识
	ChainName         string //链名称
	Consensus         string // 当前链使用的共识算法
	TxTimeout         uint32 // 交易时间戳的过期时间(秒)
	BlockTxCapacity   uint32 // 区块中最大交易数量
	BlockInterval     uint32 // 区块最大间隔时间，单位：ms
	Policy            string // 多签策略
	Status            int    // 链状态，1：未启动，0：启动 2: 失效
	Version           string // 链版本
	Sequence          string //链配置版本
	Monitor           int    //监控开关 1：开启 0：不开启
	ChainmakerImprove int    //是否参与长安链改进计划 1：开启，0：不开启
	Address           string //管理平台地址
	AutoReport        int    // 是否自动上报错误日志 1: 开启， 0：不开启
	TLS               int    // 是否开启TLS 0: 开启， 1：不开启 （默认开启）
	DockerVm          int    // 是否开启DockerVm 0: 不开启， 1：开启 （默认不开启）
}

func (*Chain) TableName() string {
	return TableChain
}

type Transaction struct {
	CommonIntField
	ChainId string `gorm:"unique_index:chain_id_tx_id_index;index:chain_id_org_id_block_height_block_hash_timestamp_index" json:"chainId"`
	TxId    string `gorm:"unique_index:chain_id_tx_id_index" json:"txId"` //交易id
	// OrgId 发起交易的组织
	OrgId  string `gorm:"index:chain_id_org_id_block_height_block_hash_timestamp_index" json:"orgId"`
	Sender string `json:"sender"`
	// BlockHeight 区块高度
	BlockHeight uint64 `gorm:"index:chain_id_org_id_block_height_block_hash_timestamp_index" json:"blockHeight"`
	BlockHash   string
	TxType      string `json:"txType"` //交易类型
	// Timestamp 交易时间戳
	Timestamp           int64  `gorm:"index:chain_id_org_id_block_height_block_hash_timestamp_index" json:"timestamp"`
	TxStatusCode        string `json:"txStatusCode"`                            //交易状态码
	ContractName        string `json:"contractName"`                            //合约名称
	ContractMethod      string `json:"contractMethod"`                          //合约方法
	ContractParameters  string `gorm:"type:longtext" json:"contractParameters"` //查询参数
	ContractVersion     string `json:"contractVersion"`                         //合约版本
	ContractRuntimeType string `json:"contractRuntimeType"`                     //合约运行时版本
	Endorsers           string `gorm:"type:longtext" json:"endorsers"`          //签名者签名集合
	Sequence            uint64 `json:"sequence"`                                //配置序列
	TXResult
}

type TXResult struct {
	ContractResult []byte `gorm:"type:mediumblob" json:"contractResult"` //合约结果
	ResultCode     string `json:"result_code"`                           //合约结果码
	ResultMessage  string `json:"result_message"`                        //合约结果信息
	RwSetHash      string `json:"rwSetHash"`                             //读写集哈希

	// contract result 的信息解析
	ContractResultCode uint32 `json:"contractResultCode"`
	//ContractResultResult []byte `gorm:"type:mediumblob" json:"contractResult"`   //合约结果
	ContractResultMessage string `json:"contractResultMessage"` //合约结果信息
	Gas                   uint64 `json:"gas"`
}

func (*Transaction) TableName() string {
	return TableTransaction
}

// package chain_participant
type Cert struct {
	CommonIntField
	CertType     int    `gorm:"index:cert_type_index"` // 证书类型 5：普通节点证书，4：共识节点证书，3：用户client证书，2：用户admin证书, 1：ca证书， 0：根证书
	CertUse      int    // 证书用途 1：tls，0：签名
	Cert         string `gorm:"type:text"` //证书值
	PrivateKey   string `gorm:"type:text"` //私钥值
	OrgId        string // 组织id
	OrgName      string // 组织名称
	CertUserName string // 证书用户名
	NodeName     string // 节点名
}

func (*Cert) TableName() string {
	return TableCert
}

type Node struct {
	CommonIntField
	NodeName string `gorm:"index:node_name_node_id_index"` // 节点名称
	NodeId   string `gorm:"index:node_name_node_id_index"` // 节点id
	NodeIp   string // 节点ip
	NodePort string // 节点端口
	Type     int
}

func (*Node) TableName() string {
	return TableNode
}

type Org struct {
	CommonIntField
	OrgId   string `gorm:"unique_index:org_id_index"` // 组织id
	OrgName string // 组织名称
}

func (*Org) TableName() string {
	return TableOrg
}

// package contract
type Contract struct {
	CommonIntField
	ChainId          string `gorm:"unique_index:chain_id_name_version_index" json:"chainId"` //子链标识
	Name             string `gorm:"unique_index:chain_id_name_version_index" json:"name"`    //合约名称
	Version          string `gorm:"unique_index:chain_id_name_version_index" json:"version"` //合约版本
	RuntimeType      int    `json:"runtimeType"`                                             //运行时版本
	SourceSaveKey    string // 合约源码存储的key
	EvmAbiSaveKey    string // evm合约abi存储的key
	EvmFunctionType  int    // 0：正常方法 1：构造函数
	EvmAddress       string // evm链上合约名
	ContractOperator string // 合约发布的操作员
	MgmtParams       string `gorm:"type:mediumtext"` // 合约操作的参数
	Methods          string `gorm:"type:mediumtext"` // 合约方法
	ContractStatus   int    // 合约状态，-1：未知（可能在链上，未在管理平台）；0：已存储；1：发布成功；2：发布失败
	BlockHeight      uint64 // 当前合约操作所在区块高度
	TxId             string // 创建合约的交易id
	OrgId            string `json:"org_id"` //合约的发起组织
	MultiSignStatus  int
	Timestamp        int64 `gorm:"column:timestamp"` // 创建时间
	TxNum            int64
}

func (*Contract) TableName() string {
	return TableContract
}

type InvokeRecords struct {
	CommonIntField
	ChainId      string // 链id
	OrgId        string // 组织id
	OrgName      string // 组织名称
	ContractName string // 合约名
	TxId         string `gorm:"index:tx_id_index"` // 证书类型 5：普通节点证书，4：共识节点证书，3：用户client证书，2：用户admin证书, 1：ca证书， 0：根证书
	TxStatus     int    // 交易状态 ，0：成功；其余：失败
	Status       int    // 上链状态 ，0：上链中；1：已上链，2：上链失败
	Method       string
	UserName     string
}

func (*InvokeRecords) TableName() string {
	return TableInvokeRecords
}

// package policy
type ChainPolicy struct {
	CommonIntField
	ChainId string `gorm:"index:chain_id_index"`
	// 权限类型
	// 0:NODE_ADDR_UPDATE; 1:TRUST_ROOT_UPDATE 2:CONSENSUS_EXT_DELETE; 3:BLOCK_UPDATE; 4: INIT_CONTRACT;
	// 5:UPGRADE_CONTRACT; 6：FREEZE_CONTRACT; 7: UNFREEZE_CONTRACT; 8:REVOKE_CONTRACT
	Type       int
	PolicyType int //策略类型 0:Majority; 1:Any; 2:Self; 3:All 4:Forbidden; 5:All
	RoleType   int //角色类型 0:Admin; 1:Client; 2:All
	PercentNum string
}

func (*ChainPolicy) TableName() string {
	return TableChainPolicy
}

type ChainPolicyOrg struct {
	CommonIntField
	ChainPolicyId int64  `gorm:"unique_index:chain_policy_id_org_id_index"`
	OrgId         string `gorm:"unique_index:chain_policy_id_org_id_index"`
	OrgName       string
	Status        int
}

func (*ChainPolicyOrg) TableName() string {
	return TableChainPolicyOrg
}

// package relation
type ChainOrg struct {
	CommonIntField
	ChainId string `gorm:"index:org_id_chain_id_index"`
	OrgId   string `gorm:"index:org_id_chain_id_index"`
	OrgName string
}

func (*ChainOrg) TableName() string {
	return TableChainOrg
}

type ChainOrgNode struct {
	CommonIntField
	ChainId  string `gorm:"index:chain_id_index"`
	OrgId    string
	OrgName  string
	NodeId   string `gorm:"index:node_id_index"`
	NodeName string
}

func (*ChainOrgNode) TableName() string {
	return TableChainOrgNode
}

type OrgNode struct {
	CommonIntField
	OrgId    string `gorm:"index:org_id_index"`
	OrgName  string
	NodeId   string
	NodeName string
}

func (*OrgNode) TableName() string {
	return TableOrgNode
}

// package db
type Upload struct {
	CommonIntField
	UserId    int64 `gorm:"index:user_id_index"`
	Hash      string
	FileName  string
	Extension string
	Content   []byte `gorm:"type:mediumblob"`
}

func (*Upload) TableName() string {
	return TableUploadContent
}

type User struct {
	CommonIntField
	UserName string `gorm:"unique_index:user_name_index"` // 用户名，唯一
	Name     string // 用户姓名
	Salt     string // 盐
	Passwd   string // 密码
	ParentId int64  `gorm:"index:parent_id_index"` // 父用户
	Status   int    // 用户状态 0:启用；1:禁用
}

func (*User) TableName() string {
	return TableUser
}

type VoteManagement struct {
	CommonIntField
	MultiId      string `json:"multiId"`      //多签事件唯一标识
	ChainId      string `json:"chainId"`      //子链标识
	StartOrgId   string `json:"blockHash"`    //发起组织
	StartOrgName string `json:"startOrgName"` //发起组织名称
	VoteOrgId    string `json:"voteOrgId"`    //投票组织ID
	VoteOrgName  string `json:"voteOrgName"`  // 投票组织名称
	// 投票时间类型
	// 0:NODE_ADDR_UPDATE; 1:TRUST_ROOT_UPDATE 2:CONSENSUS_EXT_DELETE; 3:BLOCK_UPDATE; 4: INIT_CONTRACT;
	// 5:UPGRADE_CONTRACT; 6：FREEZE_CONTRACT; 7: UNFREEZE_CONTRACT; 8:REVOKE_CONTRACT
	VoteType     int    `gorm:"index:vote_type_index" json:"voteType"`
	PolicyType   int    `json:"policyType"`                                //策略类型 0:Majority; 1:Any; 2:Self; 3:All 4:Forbidden; 5:All
	PassPercent  string `json:"passPercent"`                               //通过率
	VoteResult   int    `json:"voteResult"`                                //投票结果 0:未投票; 1:同意; 2:反对
	VoteStatus   int    `gorm:"index:vote_status_index" json:"voteStatus"` //投票状态 0:投票中; 1:投票完成
	Params       string `gorm:"type:mediumtext"`
	Reason       string `gorm:"type:mediumtext" json:"reason"` //发起投票原因
	VoteDetail   string `gorm:"type:longtext"`
	ConfigStatus int    //更新配置块 0：修改链权限； 1：修改链配置；2：其他
}

func (*VoteManagement) TableName() string {
	return TableVoteManagement
}

type ChainConfig struct {
	CommonIntField
	ChainId     string `gorm:"unique_index:chain_id_block_unique_index" json:"chainId"`     //链id
	BlockHeight uint64 `gorm:"unique_index:chain_id_block_unique_index" json:"blockHeight"` // 配置块的高度
	BlockTime   int64  `json:"blockTime"`                                                   // 配置块的时间
	Config      string `gorm:"type:longtext" json:"config"`                                 //链配置
}

func (*ChainConfig) TableName() string {
	return TableChainConfig
}

type ChainErrorLog struct {
	CommonIntField
	ChainId string `gorm:"index:chain_id_index" json:"chain_id"` //链id
	NodeId  string `json:"node_id"`                              // 节点id
	Type    string `json:"type"`                                 // 错误类型
	LogId   string `json:"log_id"`                               // 错误日志id
	LogTime int64  `json:"log_time"`                             // 日志时间
	Log     string `gorm:"type:longtext" json:"config"`          // 日志详细内容
}

func (*ChainErrorLog) TableName() string {
	return TableChainErrorLog
}
