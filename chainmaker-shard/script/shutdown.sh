#!/usr/bin/env bash
#
# Copyright (C) BABEC. All rights reserved.
# Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

function shutdown() {
    pid=`ps -ef | grep cross_chain | grep -v grep | awk '{print $2}'`
    if [ ! -z ${pid} ];then
        kill -9 $pid
    fi
    echo "cross_chain is stopped"
}

shutdown