/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package conf

import (
	"fmt"
	"strconv"

	"chainmaker.org/chainmaker-cross/logger"
)

const (
	StringToByteIndex = 0
)

// LocalConf Local config struct
type LocalConf struct {
	ListenerConfig *ListenerConfig           `mapstructure:"listener"` // 本地服务配置
	AdapterConfigs AdapterConfigs            `mapstructure:"adapters"` // 转接器配置
	RouterConfigs  []*RouterConfig           `mapstructure:"routers"`  // 路由配置
	ProverConfigs  []*ProverConfig           `mapstructure:"provers"`  // 证明器配置
	StorageConfig  *StorageConfig            `mapstructure:"storage"`  // 存储配置
	LogConfig      []*logger.LogModuleConfig `mapstructure:"log"`      // 日志配置
}

// ListenerConfig Listener config
type ListenerConfig struct {
	WebConfig     *WebConfig     `mapstructure:"web"`     // web服务配置
	ChannelConfig *ChannelConfig `mapstructure:"channel"` // P2p网络配置
	GrpcConfig    *GrpcConfig    `mapstructure:"grpc"`    // grpc服务配置
}

// WebConfig WebListener config
type WebConfig struct {
	Address     string             `mapstructure:"address"`        // web服务监听地址
	Port        int                `mapstructure:"port"`           // web服务监听端口
	OpenTxRoute bool               `mapstructure:"open_tx_router"` // web服务开启事务处理路由
	EnableTLS   bool               `mapstructure:"enable_tls"`     //启用tls
	Security    *TransportSecurity `mapstructure:"security"`       //传输安全配置
}

// ToUrl return url of web config
func (webConfig *WebConfig) ToUrl() string {
	return webConfig.Address + ":" + strconv.Itoa(webConfig.Port)
}

// ChannelConfig ChannelListener config
type ChannelConfig struct {
	Provider      string               `mapstructure:"provider"` // P2p网络类型，如 libp2p，添加 Provider 需要扩展该类型
	LibP2PChannel *LibP2PChannelConfig `mapstructure:"libp2p"`   // libp2p 网络配置
}

// GrpcConfig Grpc config
type GrpcConfig struct {
	Network string `mapstructure:"protocol"`    // 网络类型
	Address string `mapstructure:"listen_port"` // 监听地址
}

// LibP2PChannelConfig LibP2P channel config
type LibP2PChannelConfig struct {
	Address     string `mapstructure:"address"`       // listen address
	PrivKeyFile string `mapstructure:"priv_key_file"` // p2p network peer id derived form private key
	ProtocolID  string `mapstructure:"protocol_id"`   // p2p network protocolID
	Delimit     string `mapstructure:"delimit"`       // p2p network delimit
}

// GetDelimit return delimit of libp2p connection
func (c *LibP2PChannelConfig) GetDelimit() byte {
	bs := []byte(c.Delimit)
	if len(bs) > 1 {
		panic("delimit config more than one rune")
	}
	return bs[StringToByteIndex]
}

// StorageConfig storage config
type StorageConfig struct {
	Provider string         `mapstructure:"provider"` // 存储类型
	LevelDB  *LevelDBConfig `mapstructure:"leveldb"`  // levelBD 配置
}

// LevelDBConfig leveldb config
type LevelDBConfig struct {
	StorePath       string `mapstructure:"store_path"` // 存储路径
	WriteBufferSize int    `mapstructure:"write_buffer_size"`
	BloomFilterBits int    `mapstructure:"bloom_filter_bits"`
}

// RouterConfig the config of router
type RouterConfig struct {
	Provider     string              `mapstructure:"provider"`  // 路由网络类型
	ChainIDs     []string            `mapstructure:"chain_ids"` // 代理节点能直连的链
	LibP2PRouter *LibP2PRouterConfig `mapstructure:"libp2p"`    // libp2p 网络配置
	HttpRouter   *HttpRouterConfig   `mapstructure:"http"`      // http 网络配置
}

// LibP2PRouterConfig the config of libp2p router
type LibP2PRouterConfig struct {
	Address           string `mapstructure:"address"`            // libp2p 网络地址
	ProtocolID        string `mapstructure:"protocol_id"`        // p2p network protocolID
	Delimit           string `mapstructure:"delimit"`            // p2p network delimit
	ReconnectLimit    int    `mapstructure:"reconnect_limit"`    // 连接断开重试次数
	ReconnectInterval int    `mapstructure:"reconnect_interval"` // 连接间隔， 单位毫秒
}

type HttpRouterConfig struct {
	Address         string `mapstructure:"address"` // http 网络地址
	HttpTransport   `mapstructure:",squash"`
	RequestStrategy `mapstructure:",squash"`
}

type HttpTransport struct {
	MaxConnection    int                `mapstructure:"max_connection"`    // 连接池最大连接数
	IdleIConnTimeout int                `mapstructure:"idle_conn_timeout"` // 空闲连接超时时间单位s
	EnableTLS        bool               `mapstructure:"enable_tls"`        //启用tls
	EnableH2         bool               `mapstructure:"enable_h2"`         //是否启用http2
	Security         *TransportSecurity `mapstructure:"security"`          //传输安全配置
}

type RequestStrategy struct {
	RequestTimeout       int `mapstructure:"request_timeout"`        // 请求超时时间
	RequestMaxRetries    int `mapstructure:"request_max_retries"`    //请求失败重试次数
	RequestRetryInterval int `mapstructure:"request_retry_interval"` //请求失败重试间隔时间(ms)
}

type TransportSecurity struct {
	CAFile         string `mapstructure:"ca_file"`   //信任的ca文件，多个文件逗号分割
	EnableCertAuth bool   `mapstructure:"ca_auth"`   //是否开启证书验证
	CertFile       string `mapstructure:"cert_file"` //证书文件
	KeyFile        string `mapstructure:"key_file"`  //私钥文件
}

// GetDelimit return delimit of libp2p router config
func (c *LibP2PRouterConfig) GetDelimit() byte {
	bs := []byte(c.Delimit)
	if len(bs) > 1 {
		panic("delimit config more than one rune")
	}
	return bs[StringToByteIndex]
}

// GetChainIDs return chain-ids of remote cross-chain proxy
func (r *RouterConfig) GetChainIDs() []string {
	return r.ChainIDs
}

// AdapterConfig adapter config
type AdapterConfig struct {
	Provider      string              `mapstructure:"provider"`       // 转接器类型
	ChainID       string              `mapstructure:"chain_id"`       // 转接器连接的链的ID
	ConfigPath    string              `mapstructure:"config_path"`    // 配置路径
	ProofContract *ProofContract      `mapstructure:"proof_contract"` // 证据保存的合约信息
	ExtraConf     map[string][]string `mapstructure:"extra_conf"`     // 各个平行链的个性化配置
}

// ProofContract contract for save proof
type ProofContract struct {
	Name   string `mapstructure:"name"`   // 证据存储的合约名称
	Method string `mapstructure:"method"` // 证据存储的方法名称
}

type AdapterConfigs []*AdapterConfig

func (adapters AdapterConfigs) GetExtraConfigByProvider(provider string) (*AdapterConfig, error) {
	for _, a := range adapters {
		if a.Provider == provider {
			return a, nil
		}
	}
	return nil, fmt.Errorf("cant find config by given provider: %s", provider)
}

func (adapters AdapterConfigs) GetExtraConfigByKey(provider string, key string) ([]string, error) {
	for _, a := range adapters {
		if a.Provider == provider {
			if value, ok := a.ExtraConf[key]; !ok {
				return nil, fmt.Errorf("cant find config by given provider: %s and key: %s", provider, key)
			} else {
				return value, nil
			}
		}
	}
	return nil, fmt.Errorf("cant find config by given provider: %s and key: %s", provider, key)
}

// ProverConfig prover config
type ProverConfig struct {
	Provider   string   `mapstructure:"provider"`    // 证明器类型，例如 spv，trust等
	ChainIDs   []string `mapstructure:"chain_ids"`   // 证明器连接的链的ID
	ConfigPath string   `mapstructure:"config_path"` // 证明器配置路径
}

// GetChainIDs return chain-ids of local provers
func (p *ProverConfig) GetChainIDs() []string {
	return p.ChainIDs
}
