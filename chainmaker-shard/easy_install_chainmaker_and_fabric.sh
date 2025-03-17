#!/usr/bin/env bash
#
# Copyright (C) BABEC. All rights reserved.
# Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

VERSION=V1.0.0
DATETIME=$(date "+%Y%m%d%H%M%S")
PLATFORM=$(uname -m)

PROJECT_PATH=$(cd `dirname $0`;pwd)
RELEASE_PATH=${PROJECT_PATH}/release
CONFIG_PATH=${PROJECT_PATH}/config
CM_CONFIG_PATH=${CONFIG_PATH}/chainmaker
FABRIC_CONFIG_PATH=${CONFIG_PATH}/fabric
CM_CHAIN_ID="chain1"
FABRIC_CHAIN_ID="mychannel"
IP="127.0.0.1"
PORT1="12301"
PORT2="12302"
CM_TRANSACTION_CONTRACT="Transaction"
CM_SAVE_PROOF_METHOD="SaveProof"
FABRIC_TRANSACTION_CONTRACT="Transaction"
FABRIC_SAVE_PROOF_METHOD="SaveProof"

CHAIN_ID1="chain1"
CHAIN_TYPE1="chainmaker"
TRANSACTION_CONTRACT1="CROSS_TRANSACTION"
SAVE_PROOF_METHOD1="SaveProof"
SPV_TYPE1="ChainMaker_Light"

CHAIN_ID2="mychannel"
CHAIN_TYPE2="fabric"
TRANSACTION_CONTRACT2="Transaction"
SAVE_PROOF_METHOD2="SaveProof"
SPV_TYPE2="Fabric_SPV"

ADAPTER_CONFIG1="config/chainmaker/chainmaker_sdk.yml"
ADAPTER_CONFIG2="config/fabric/fabric_sdk.yml"

function prepare_env() {
    # 创建文件结构
    echo "prepare directory structure..."
	  mkdir -p ${RELEASE_PATH}/bin
    mkdir -p ${RELEASE_PATH}/lib
	  mkdir -p ${RELEASE_PATH}/config
	  mkdir -p ${RELEASE_PATH}/logs
	  mkdir -p ${RELEASE_PATH}/config/chainmaker
    mkdir -p ${RELEASE_PATH}/config/fabric

	  # 复制脚本
    cp ${PROJECT_PATH}/script/start.sh ${RELEASE_PATH}/bin
    cp ${PROJECT_PATH}/script/shutdown.sh ${RELEASE_PATH}/bin

    # 初始化配置文件
    cp ${CONFIG_PATH}/cross_chain.yml ${RELEASE_PATH}/config
    cp ${CONFIG_PATH}/spv.yml ${RELEASE_PATH}/config/spv.yml
    cp ${CM_CONFIG_PATH}/chainmaker_sdk.yml ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk.yml
    cp ${FABRIC_CONFIG_PATH}/fabric_sdk.yml ${RELEASE_PATH}/config/fabric/fabric_sdk.yml

    # 创建代理节点私钥
    openssl ecparam -out ${RELEASE_PATH}/config/ecprikey.key -name prime256v1 -noout -genkey

    # 将私钥转换身份ID，写入文件

    # 指定 crypto-config 路径
    read -p "请指定 chainmaker crypto-config 文件夹的路径: " CM_CRPTO_CONFIG
    if [ ! -d "$CM_CRPTO_CONFIG" ]; then
      echo "crypto-config given path not exist"
      exit
    fi
    cp -r $CM_CRPTO_CONFIG ${RELEASE_PATH}/config/chainmaker/

    read -p "请指定 fabric organizations 文件夹的路径: " FABRIC_ORGANIZATIONS
    if [ ! -d "$FABRIC_ORGANIZATIONS" ]; then
      echo "organizations given path not exist"
      exit
    fi
    cp -r $FABRIC_ORGANIZATIONS ${RELEASE_PATH}/config/fabric/

    # 修改 chainmaker-sdk.yml 的配置
    xsed "s%{ CHAIN_ID }%$CM_CHAIN_ID%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk.yml
    xsed "s%{ CRYPTO_CONFIG }%$CM_CRPTO_CONFIG%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk.yml
    xsed "s%{ IP }%$IP%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk.yml
    xsed "s%{ PORT1 }%$PORT1%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk.yml
    xsed "s%{ PORT2 }%$PORT2%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk.yml

    # 修改 fabric_sdk.yml 的配置
    xsed "s%{ FABRIC_CHAIN_ID }%$FABRIC_CHAIN_ID%g" ${RELEASE_PATH}/config/fabric/fabric_sdk.yml
    xsed "s%{ CRYPTO_CONFIG }%$FABRIC_ORGANIZATIONS%g" ${RELEASE_PATH}/config/fabric/fabric_sdk.yml


   for i in {1..2}
      do
        v_chain=CHAIN_ID$i
        chain=`eval echo '$'"$v_chain"`
        v_chain_type=CHAIN_TYPE$i
        chain_type=`eval echo '$'"$v_chain_type"`
        v_tx_contract=TRANSACTION_CONTRACT$i
        tx_contract=`eval echo '$'"$v_tx_contract"`
        v_method=SAVE_PROOF_METHOD$i
        method=`eval echo '$'"$v_method"`
        v_spv_type=SPV_TYPE$i
        spv_type=`eval echo '$'"$v_spv_type"`
        v_adapter_config=ADAPTER_CONFIG$i
        adapter_config=`eval echo '$'"$v_adapter_config"`

        xsed "s/{ CHAIN_ID_$i }/$chain/" ${RELEASE_PATH}/config/spv.yml
        xsed "s/{ SPV_TYPE_$i }/$spv_type/" ${RELEASE_PATH}/config/spv.yml
        xsed "s%{ ADAPTER_CONFIG_PATH_$i }%$adapter_config%" ${RELEASE_PATH}/config/spv.yml

        xsed "s/{ CHAIN_TYPE_$i }/$chain_type/" ${RELEASE_PATH}/config/cross_chain.yml
        xsed "s/{ CHAIN_ID_$i }/$chain/" ${RELEASE_PATH}/config/cross_chain.yml
        xsed "s/{ TRANSACTION_CONTRACT_$i }/$tx_contract/" ${RELEASE_PATH}/config/cross_chain.yml
        xsed "s/{ SAVE_PROOF_METHOD_$i }/$method/" ${RELEASE_PATH}/config/cross_chain.yml
        xsed "s%{ ADAPTER_CONFIG_PATH_$i }%$adapter_config%" ${RELEASE_PATH}/config/cross_chain.yml
      done

      read -p "指定tls配置(Y/N) " tls
          case $tls in
          [yY][eE][sS]|[yY])
               # 修改tls配置
                      read -p "请指定 CA 证书路径: " CA_PATH
                      if [ ! -f "$CA_PATH" ]; then
                        echo "CA_PATH not exist"
                        exit
                      fi
                      read -p "请指定用户 TLS 证书路径: " TLS_CRT_PATH
                      if [ ! -f "$TLS_CRT_PATH" ]; then
                        echo "TLS_CRT_PATH not exist"
                        exit
                      fi
                      read -p "请指定用户 TLS 私钥路径: " TLS_KEY_PATH
                      if [ ! -f "$TLS_KEY_PATH" ]; then
                        TLS_KEY_PATH=""
                        echo "TLS_KEY_PATH not exist"
                        exit
                      fi
                      ;;
          esac

          xsed "s%{ TLS_CRT_PATH }%$TLS_CRT_PATH%g" ${RELEASE_PATH}/config/spv.yml
          xsed "s%{ TLS_KEY_PATH }%$TLS_KEY_PATH%g" ${RELEASE_PATH}/config/spv.yml
          xsed "s%{ CA_FILE_PATH }%$CA_PATH%g" ${RELEASE_PATH}/config/spv.yml

          xsed "s%{ TLS_CRT_PATH }%$TLS_CRT_PATH%g" ${RELEASE_PATH}/config/cross_chain.yml
          xsed "s%{ TLS_KEY_PATH }%$TLS_KEY_PATH%g" ${RELEASE_PATH}/config/cross_chain.yml
          xsed "s%{ CA_FILE_PATH }%$CA_PATH%g" ${RELEASE_PATH}/config/cross_chain.yml
          xsed "46,75s/^/#/" ${RELEASE_PATH}/config/cross_chain.yml

    echo "prepare directory structure ok"
}

function build() {
    cd $PROJECT_PATH
    make
    echo "build release finished"
}

function xsed() {
    system=$(uname)

    if [ "${system}" = "Linux" ]; then
        sed -i "$@"
    else
        sed -i '' "$@"
    fi
}

prepare_env
#build
