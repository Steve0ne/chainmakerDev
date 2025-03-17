#!/usr/bin/env bash

set -e

RELEASE_PATH=$(cd $(dirname "${BASH_SOURCE[0]}") && pwd)

cd $RELEASE_PATH
for file in `ls $RELEASE_PATH`
    do
        if [ -d $file ]; then
            echo "START ==> " $RELEASE_PATH/$file
            cd $file/bin && ./restart.sh && cd - > /dev/null
        fi
    done
