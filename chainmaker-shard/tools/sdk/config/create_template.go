package conf

import (
	"io/ioutil"
	"os"
)

var CM_SDK_TEMPLATE = "chain_client:\n  # 链ID\n  chain_id: CHAIN_ID\n  # 组织ID\n  org_id: ORG_ID\n  # 客户端用户私钥路径，使用绝对路径\n  user_key_file_path: USER_KEY_PATH\n  # 客户端用户证书路径，使用绝对路径\n  user_crt_file_path: USER_KEY_PATH\n  # 客户端用户交易签名私钥路径(若未设置，将使用user_key_file_path)，使用绝对路径\n  user_sign_key_file_path: SIGN_KEY_PATH\n  # 客户端用户交易签名证书路径(若未设置，将使用user_crt_file_path)，使用绝对路径\n  user_sign_crt_file_path: SIGN_CRT_PATH\n\n  nodes:\n    - # 节点地址，格式为：IP:端口:连接数\n      node_addr: IP:PORT\n      # 节点连接数\n      conn_cnt: 5\n      # RPC连接是否启用双向TLS认证\n      enable_tls: false\n      # 信任证书池路径，使用绝对路径\n      trust_root_paths:\n      # TLS hostname\n      tls_host_name: TLS_HOST_NAME\n"

var FABRIC_SDK_TEMPLATE = "version: 1.0.0\n\norganizations:\n  Org:\n    mspid: OrgMSP\n    users:\n      User:\n        key:\n          path: PRIVATE_KEY_PATH\n        cert:\n          path: CERT_PATH\n"

type delFunc func()

func createChainMakerSdkTmpFile() (delFunc, string, error) {
	var err error
	// 创建临时文件目录
	p, err := ioutil.TempDir(os.TempDir(), "cm_sdk_tpl")
	if err != nil {
		return nil, "", err
	}
	path := p + ".yml"
	// 写入模版
	err = ioutil.WriteFile(path, []byte(CM_SDK_TEMPLATE), 0600)
	if err != nil {
		return nil, "", err
	}
	// 返回回调函数
	var delF delFunc = func() {
		os.RemoveAll(path)
	}

	return delF, path, nil
}

func createFabricSdkTmpFile() (delFunc, string, error) {
	var err error
	// 创建临时文件目录
	p, err := ioutil.TempDir(os.TempDir(), "fabric_sdk_tpl")
	if err != nil {
		return nil, "", err
	}
	path := p + ".yml"
	// 写入模版
	err = ioutil.WriteFile(path, []byte(FABRIC_SDK_TEMPLATE), 0600)
	if err != nil {
		return nil, "", err
	}
	// 返回回调函数
	var delF delFunc = func() {
		os.RemoveAll(path)
	}

	return delF, path, nil
}
