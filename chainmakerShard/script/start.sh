#!/usr/bin/env bash
#
# Copyright (C) BABEC. All rights reserved.
# Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

set -e

CURRENT_PATH=$(cd `dirname $0`;pwd)
RELEASE_PATH=$(dirname "${CURRENT_PATH}")
CONFIG_PATH=${RELEASE_PATH}/config
BINARY_PATH=${RELEASE_PATH}/bin
LIB_PATH=${RELEASE_PATH}/lib
BINARY_FILE=${LIB_PATH}/cross-chain

function check_binary() {
    if  [ ! -d $CONFIG_PATH ] ;then
        echo $CONFIG_PATH" is missing"
        exit 1
    fi

    if  [ ! -d $BINARY_PATH ] ;then
        echo $BINARY_PATH" is missing"
        exit 1
    fi

    if  [ ! -e $LIB_PATH ] ;then
        echo $LIB_PATH" is missing"
        exit 1
    fi

    if  [ ! -e $BINARY_FILE ] ;then
        echo $BINARY_FILE" is missing"
        exit 1
    fi

}

function start() {
    pid=`ps -ef | grep cross_chain | grep -v grep | awk '{print $2}'`
    if [ -z ${pid} ];then
        nohup ${BINARY_FILE} start -c ${CONFIG_PATH}/cross_chain.yml -d ${RELEASE_PATH} > ${RELEASE_PATH}/logs/console.log &
        echo "cross-chain start, pls check log..."
    else
        echo "cross-chain is already started"
    fi
}

check_binary
start

