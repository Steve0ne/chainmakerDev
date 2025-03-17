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
SDK_PATH=${PROJECT_PATH}/tools/sdk
SDK_TEMPLATE_PATH=${SDK_PATH}/config/template

CM_CHAIN1="chain1"
CM_ORG_ID1="wx-org1.chainmaker.org"
CM_CHAIN2="chain2"
CM_ORG_ID2="wx-org1.chainmaker.org"
FABRIC_CHANNEL="mychannel"
FABRIC_ORG_ID="Org1MSP"
#FABRIC_ORG_USER="User1"
#ORG_PEER1="peer0.org1.example.com"
#ORG_PEER2="peer0.org2.example.com"


function prepare_env() {
    # 初始化配置文件
    cp ${SDK_TEMPLATE_PATH}/cross_chain_sdk.yml ${RELEASE_PATH}/config

    # 指定 crypto-config 路径
    read -p "请指定跨链模式: 1) chainmaker - chainmaker, 2) chainmaker - fabric " CROSS_TYPE
    if [ "$CROSS_TYPE" = "1" ]
    then
      echo "easy create [ chainmaker - chainmaker ] cross configs by default"
      cross_chainmaker_and_chainmaker

    elif [ "$CROSS_TYPE" = "2" ]
    then
      echo "easy create [ chainmaker - fabric ] cross configs by default"
      cross_chainmaker_and_fabric
    else
      echo
    fi

    echo "prepare sdk config ok"
}

function cross_chainmaker_and_chainmaker() {
    # 修改 chain_id 和 org_id
    xsed "2s%{ CM_CHAIN_ID }%$CM_CHAIN1%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "4s%{ CM_ORG_ID }%$CM_ORG_ID1%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "9s%{ CM_CHAIN_ID }%$CM_CHAIN2%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "11s%{ CM_ORG_ID }%$CM_ORG_ID2%" ${RELEASE_PATH}/config/cross_chain_sdk.yml

    # 输入 证书 和 私钥 路径
    read -p "请输入第一条链的用户证书路径: " SIGN_CRT_PATH1
    read -p "请输入第一条链的用户私钥路径: " SIGN_KEY_PATH1
    # 检查路径文件
    if [ ! -f $SIGN_CRT_PATH1 ]; then
      echo "file: $SIGN_CRT_PATH1 not exist"
      exit
    fi

    if [ ! -f $SIGN_KEY_PATH1 ]; then
      echo "file: $SIGN_KEY_PATH1 not exist"
      exit
    fi
    # 替换路径
    xsed "6s%{ SIGN_CRT_PATH }%$SIGN_CRT_PATH1%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "5s%{ SIGN_KEY_PATH }%$SIGN_KEY_PATH1%" ${RELEASE_PATH}/config/cross_chain_sdk.yml

    # 输入 证书 和 私钥 路径
    read -p "请输入第二条链的用户证书路径: " SIGN_CRT_PATH2
    read -p "请输入第二条链的用户私钥路径: " SIGN_KEY_PATH2
    # 检查路径文件
    if [ ! -f $SIGN_CRT_PATH2 ]; then
      echo "file: $SIGN_CRT_PATH2 not exist"
      exit
    fi

    if [ ! -f $SIGN_KEY_PATH2 ]; then
      echo "file: $SIGN_KEY_PATH2 not exist"
      exit
    fi
    # 替换路径
    xsed "13s%{ SIGN_CRT_PATH }%$SIGN_CRT_PATH2%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "12s%{ SIGN_KEY_PATH }%$SIGN_KEY_PATH2%" ${RELEASE_PATH}/config/cross_chain_sdk.yml

    # 注释非引入参数
    xsed "16,31s/^/#/" ${RELEASE_PATH}/config/cross_chain_sdk.yml
}

function cross_chainmaker_and_fabric() {
    # 修改 chain_id 和 org_id
    xsed "2s%{ CM_CHAIN_ID }%$CM_CHAIN1%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "4s%{ CM_ORG_ID }%$CM_ORG_ID1%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "16s%{ FABRIC_CHAIN_ID }%$FABRIC_CHANNEL%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "18s%{ FABRIC_ORG_ID }%$FABRIC_ORG_ID%" ${RELEASE_PATH}/config/cross_chain_sdk.yml

    # 输入 tls 证书路径
    read -p "请输入跨链SDK的 TLS-CA 证书路径: " CA_PATH
    if [ ! -f $CA_PATH ]; then
      echo "file: $CA_PATH not exist"
      exit
    fi
    # 替换证书路径
    xsed "s%{ CA_PATH }%$CA_PATH%" ${RELEASE_PATH}/config/cross_chain_sdk.yml

    # 输入 tls 证书路径
    read -p "请输入跨链SDK的 TLS 证书路径: " CRT_PATH
    read -p "请输入跨链SDK的 TLS 私钥路径: " KEY_PATH
    if [ ! -f $CRT_PATH ]; then
      echo "file: $CRT_PATH not exist"
      exit
    fi
    if [ ! -f $KEY_PATH ]; then
      echo "file: $KEY_PATH not exist"
      exit
    fi

    # 替换 证书 和 私钥 路径
    xsed "s%{ CRT_PATH }%$CRT_PATH%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "s%{ KEY_PATH }%$KEY_PATH%" ${RELEASE_PATH}/config/cross_chain_sdk.yml

    # 输入 证书 和 私钥 路径
    read -p "请输入长安链的用户证书路径: " CM_SIGN_CRT_PATH
    read -p "请输入长安链的用户私钥路径: " CM_SIGN_KEY_PATH
    # 检查路径文件
    if [ ! -f $SIGN_CRT_PATH ]; then
      echo "file: $SIGN_CRT_PATH not exist"
      exit
    fi
    if [ ! -f $SIGN_KEY_PATH ]; then
      echo "file: $SIGN_KEY_PATH not exist"
      exit
    fi
    # 替换路径
    xsed "s%{ CM_SIGN_CRT_PATH }%$CM_SIGN_CRT_PATH%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "s%{ CM_SIGN_KEY_PATH }%$CM_SIGN_KEY_PATH%" ${RELEASE_PATH}/config/cross_chain_sdk.yml

    # 输入 证书 和 私钥 路径
    read -p "请输入Fabric链的用户证书路径: " FABRIC_SIGN_CRT_PATH
    read -p "请输入Fabric链的用户私钥路径: " FABRIC_SIGN_KEY_PATH
    # 检查路径文件
    if [ ! -f $FABRIC_SIGN_CRT_PATH ]; then
      echo "file: $FABRIC_SIGN_CRT_PATH not exist"
      exit
    fi

    if [ ! -f $FABRIC_SIGN_KEY_PATH ]; then
      echo "file: $FABRIC_SIGN_KEY_PATH not exist"
      exit
    fi
    # 替换路径
    xsed "s%{ FABRIC_SIGN_KEY_PATH }%$FABRIC_SIGN_KEY_PATH%" ${RELEASE_PATH}/config/cross_chain_sdk.yml
    xsed "s%{ FABRIC_SIGN_CRT_PATH }%$FABRIC_SIGN_CRT_PATH%" ${RELEASE_PATH}/config/cross_chain_sdk.yml

    # 注释非引入参数
    xsed "9,14s/^/#/" ${RELEASE_PATH}/config/cross_chain_sdk.yml
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
