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
CHAIN1="chain1"
CHAIN2="chain2"
IP="127.0.0.1"
PORT1="12301"
PORT2="12302"
CHAIN_TYPE="chainmaker"
SPV_TYPE="ChainMaker_Light"
ADAPTER_CONFIG1="config/chainmaker/chainmaker_sdk1.yml"
ADAPTER_CONFIG2="config/chainmaker/chainmaker_sdk2.yml"

function prepare_env() {
    # 创建文件结构
    echo "prepare directory structure..."
	  mkdir -p ${RELEASE_PATH}/bin
    mkdir -p ${RELEASE_PATH}/lib
	  mkdir -p ${RELEASE_PATH}/config
	  mkdir -p ${RELEASE_PATH}/logs
	  mkdir -p ${RELEASE_PATH}/config/chainmaker

	  # 复制脚本
    cp ${PROJECT_PATH}/script/start.sh ${RELEASE_PATH}/bin
    cp ${PROJECT_PATH}/script/shutdown.sh ${RELEASE_PATH}/bin

    # 初始化配置文件
    cp ${CONFIG_PATH}/cross_chain.yml ${RELEASE_PATH}/config
    cp ${CONFIG_PATH}/spv.yml ${RELEASE_PATH}/config
    cp ${CM_CONFIG_PATH}/chainmaker_sdk.yml ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk1.yml
    cp ${CM_CONFIG_PATH}/chainmaker_sdk.yml ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk2.yml

    # 创建代理节点私钥
    openssl ecparam -out ${RELEASE_PATH}/config/ecprikey.key -name prime256v1 -noout -genkey

    # 将私钥转换身份ID，写入文件

    # 指定 crypto-config 路径
    read -p "请指定 crypto-config 文件夹的路径: " CRPTO_CONFIG
    if [ ! -d "$CRPTO_CONFIG" ]; then
      echo "crypto-config given path not exist"
    fi
    cp -r $CRPTO_CONFIG ${RELEASE_PATH}/config/chainmaker/

    # 修改 chainmaker-sdk.yml 的配置
    xsed "s%{ CHAIN_ID }%$CHAIN1%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk1.yml
    xsed "s%{ CRYPTO_CONFIG }%$CRPTO_CONFIG%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk1.yml
    xsed "s%{ IP }%$IP%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk1.yml
    xsed "s%{ PORT1 }%$PORT1%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk1.yml
    xsed "s%{ PORT2 }%$PORT2%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk1.yml

    xsed "s%{ CHAIN_ID }%$CHAIN2%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk2.yml
    xsed "s%{ CRYPTO_CONFIG }%$CRPTO_CONFIG%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk2.yml
    xsed "s%{ IP }%$IP%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk2.yml
    xsed "s%{ PORT1 }%$PORT1%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk2.yml
    xsed "s%{ PORT2 }%$PORT2%g" ${RELEASE_PATH}/config/chainmaker/chainmaker_sdk2.yml

    for i in {1..2}
    do
      v=CHAIN$i
      chain=`eval echo '$'"$v"`
      v_adapter_config=ADAPTER_CONFIG$i
      adapter_config=`eval echo '$'"$v_adapter_config"`

      xsed "s/{ CHAIN_ID_$i }/$chain/" ${RELEASE_PATH}/config/spv.yml
      xsed "s/{ SPV_TYPE_$i }/$SPV_TYPE/" ${RELEASE_PATH}/config/spv.yml
      xsed "s%{ ADAPTER_CONFIG_PATH_$i }%$adapter_config%" ${RELEASE_PATH}/config/spv.yml

      xsed "s/{ CHAIN_TYPE_$i }/$CHAIN_TYPE/" ${RELEASE_PATH}/config/cross_chain.yml
      xsed "s/{ CHAIN_ID_$i }/$chain/" ${RELEASE_PATH}/config/cross_chain.yml
      xsed "s/{ TRANSACTION_CONTRACT_$i }/CROSS_TRANSACTION/" ${RELEASE_PATH}/config/cross_chain.yml
      xsed "s/{ SAVE_PROOF_METHOD_$i }/SaveProof/" ${RELEASE_PATH}/config/cross_chain.yml
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
    xsed "31,35s/^/#/" ${RELEASE_PATH}/config/spv.yml

    xsed "s%{ TLS_CRT_PATH }%$TLS_CRT_PATH%g" ${RELEASE_PATH}/config/cross_chain.yml
    xsed "s%{ TLS_KEY_PATH }%$TLS_KEY_PATH%g" ${RELEASE_PATH}/config/cross_chain.yml
    xsed "s%{ CA_FILE_PATH }%$CA_PATH%g" ${RELEASE_PATH}/config/cross_chain.yml
    xsed "40,75s/^/#/" ${RELEASE_PATH}/config/cross_chain.yml

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
