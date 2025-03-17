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
    cp ${CM_CONFIG_PATH}/chainmaker_sdk.yml ${RELEASE_PATH}/config/chainmaker
    cp ${FABRIC_CONFIG_PATH}/fabric_sdk.yml ${RELEASE_PATH}/config/fabric
    cp ${CONFIG_PATH}/spv.yml ${RELEASE_PATH}/config

    # 创建代理节点私钥
    openssl ecparam -out ${RELEASE_PATH}/config/ecprikey.key -name prime256v1 -noout -genkey
    echo "directory structure ok"
}

function build() {
    cd $PROJECT_PATH
    make
    echo "next step, you need to config chain spv and sdk configurations."
    echo "GOTO ${RELEASE_PATH}/config/cross_chain.yml"
    echo "CHANGE config_path of prover AND config_path of adapters."
    echo "REFER details to ${CONFIG_PATH}/chainmaker configuration."
}

prepare_env
build
