/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package config

// chainmaker 整体配置
type Chainmaker struct {
	ChainLogConf   *ChainLogConf     `yaml:"log"`
	BlockchainConf []*BlockchainConf `yaml:"blockchain"`
	NodeConf       *NodeConf         `yaml:"node"`
	NetConf        *NetConf          `yaml:"net"`
	TxpoolConf     *TxpoolConf       `yaml:"txpool"`
	RpcConf        *RpcConf          `yaml:"rpc"`
	MonitorConf    *MonitorConf      `yaml:"monitor"`
	PprofConf      *PprofConf        `yaml:"pprof"`
	StorageConf    *StorageConf      `yaml:"storage"`
	CoreConf       *CoreConf         `yaml:"core"`
	SchedulerConf  *SchedulerConf    `yaml:"scheduler"`
	VmConf         *VmConf           `yaml:"vm"`
}

type ChainLogConf struct {
	ConfigFile string `yaml:"config_file"`
}

type BlockchainConf struct {
	ChainId string `yaml:"chainId"`
	Genesis string `yaml:"genesis"`
}

type NodeConf struct {
	NodeType        string      `yaml:"type"`
	OrgId           string      `yaml:"org_id"`
	PrivKeyFile     string      `yaml:"priv_key_file"`
	CertFile        string      `yaml:"cert_file"`
	SignerCacheSize int         `yaml:"signer_cache_size"`
	CertCacheSize   int         `yaml:"cert_cache_size"`
	Pkcs11          *Pkcs11Conf `yaml:"pkcs11"`
}

type NetConf struct {
	Provider   string   `yaml:"provider"`
	ListenAddr string   `yaml:"listen_addr"`
	Seeds      []string `yaml:"seeds"`
	Tls        *TlsConf `yaml:"tls"`
}

type TlsConf struct {
	Enabled     bool   `yaml:"enabled"`
	PrivKeyFile string `yaml:"priv_key_file"`
	CertFile    string `yaml:"cert_file"`
}

type TxpoolConf struct {
	MaxTxpoolSize       int `yaml:"max_txpool_size"`
	MaxConfigTxpoolSize int `yaml:"max_config_txpool_size"`
	FullNotifyAgainTime int `yaml:"full_notify_again_time"`
	BatchMaxSize        int `yaml:"batch_max_size"`
	BatchCreateTimeout  int `yaml:"batch_create_timeout"`
}

type RpcConf struct {
	Provider                               string          `yaml:"provider"`
	Port                                   int             `yaml:"port"`
	CheckChainConfTrustRootsChangeInterval int             `yaml:"check_chain_conf_trust_roots_change_interval"`
	Ratelimit                              *RatelimitConf  `yaml:"ratelimit"`
	Subscriber                             *SubscriberConf `yaml:"subscriber"`
	Tls                                    *RpcTlsConf     `yaml:"tls"`
}

type RpcTlsConf struct {
	Mode        string `yaml:"mode"`
	PrivKeyFile string `yaml:"priv_key_file"`
	CertFile    string `yaml:"cert_file"`
}

type SubscriberConf struct {
	Ratelimit *RatelimitConf `yaml:"ratelimit"`
}

type RatelimitConf struct {
	TokenPerSecond  int `yaml:"token_per_second"`
	TokenBucketSize int `yaml:"token_bucket_size"`
}

type MonitorConf struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type PprofConf struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port"`
}

type StorageConf struct {
	StorePath              string              `yaml:"store_path"`
	UnarchiveBlockHeight   int64               `yaml:"unarchive_block_height"`
	BlockdbConfig          *StorageDbConf      `yaml:"blockdb_config"`
	StatedbConfig          *StorageDbConf      `yaml:"statedb_config"`
	HistorydbConfig        *StorageDbConf      `yaml:"historydb_config"`
	ResultdbConfig         *StorageDbConf      `yaml:"resultdb_config"`
	DisableContractEventdb bool                `yaml:"disable_contract_eventdb"`
	ContractEventdbConfig  *StorageEventDbConf `yaml:"contract_eventdb_config"`
}

type StorageEventDbConf struct {
	Provider    string            `yaml:"provider"`
	SqldbConfig *StorageSqlDbConf `yaml:"sqldb_config"`
}

type StorageSqlDbConf struct {
	SqldbType string `yaml:"sqldb_type"`
	Dsn       string `yaml:"dsn"`
}

type StorageDbConf struct {
	Provider      string         `yaml:"provider"`
	LeveldbConfig *LevelDbDbConf `yaml:"leveldb_config"`
}

type LevelDbDbConf struct {
	StorePath string `yaml:"store_path"`
}

type Pkcs11Conf struct {
	Enabled          bool   `yaml:"enabled"`
	Library          string `yaml:"library"`
	Label            string `yaml:"label"`
	Password         string `yaml:"password"`
	SessionCacheSize int    `yaml:"session_cache_size"`
	Hash             string `yaml:"hash"`
}

type CoreConf struct {
	Evidence bool `yaml:"evidence"`
}

type SchedulerConf struct {
	RwsetLog bool `yaml:"rwset_log"`
}

type VmConf struct {
	EnableDockervm        bool   `yaml:"enable_dockervm"`
	DockervmContainerName string `yaml:"dockervm_container_name"`
	DockervmMountPath     string `yaml:"dockervm_mount_path"`
	DockervmLogPath       string `yaml:"dockervm_log_path"`
	LogInConsole          bool   `yaml:"log_in_console"`
	LogLevel              string `yaml:"log_level"`
	UdsOpen               bool   `yaml:"uds_open"`
	TxSize                int    `yaml:"tx_size"`
	UserNum               int    `yaml:"user_num"`
	TimeLimit             int    `yaml:"time_limit"`
}
