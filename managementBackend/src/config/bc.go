/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package config

// bc 整体配置
type Bc struct {
	ChainId          string                `yaml:"chain_id"`
	Version          string                `yaml:"version"`
	Sequence         int                   `yaml:"sequence"`
	AuthType         string                `yaml:"auth_type"`
	Crypto           *CryptoConf           `yaml:"crypto"`
	Contract         *ContractConf         `yaml:"contract"`
	Block            *BlockConf            `yaml:"block"`
	Core             *BcCoreConf           `yaml:"core"`
	Snapshot         *EvidenceConf         `yaml:"snapshot"`
	Scheduler        *EvidenceConf         `yaml:"scheduler"`
	Consensus        *ConsensusConf        `yaml:"consensus"`
	TrustRoots       []*TrustRootsConf     `yaml:"trust_roots"`
	ResourcePolicies []*ResourcePolicyConf `yaml:"resource_policies"`
}

type ResourcePolicyConf struct {
	ResourceName string      `yaml:"resource_name"`
	Policy       *PolicyConf `yaml:"policy"`
}

type PolicyConf struct {
	Rule     string   `yaml:"rule"`
	OrgList  []string `yaml:"org_list"`
	RoleList []string `yaml:"role_list"`
}

type ConsensusConf struct {
	Type       int32        `yaml:"type"`
	Nodes      []*NodesConf `yaml:"nodes"`
	ExtConfig  []*KvConf    `yaml:"ext_config"`
	DposConfig []*KvConf    `yaml:"dpos_config"`
}

type NodesConf struct {
	OrgId  string   `yaml:"org_id"`
	NodeId []string `yaml:"node_id"`
}

type KvConf struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type TrustRootsConf struct {
	OrgId string   `yaml:"org_id"`
	Root  []string `yaml:"root"`
}

type CryptoConf struct {
	Hash string `yaml:"hash"`
}

type ContractConf struct {
	EnableSqlSupport bool `yaml:"enable_sql_support"`
}

type BlockConf struct {
	TxTimestampVerify bool   `yaml:"tx_timestamp_verify"`
	TxTimeout         uint32 `yaml:"tx_timeout"`
	BlockTxCapacity   uint32 `yaml:"block_tx_capacity"`
	BlockSize         int    `yaml:"block_size"`
	BlockInterval     int    `yaml:"block_interval"`
}

type BcCoreConf struct {
	TxSchedulerTimeout         int `yaml:"tx_scheduler_timeout"`
	TxSchedulerValidateTimeout int `yaml:"tx_scheduler_validate_timeout"`
}

type EvidenceConf struct {
	EnableEvidence bool `yaml:"enable_evidence"`
}
