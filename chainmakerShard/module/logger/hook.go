/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package logger

import (
	"io"

	"gopkg.in/natefinch/lumberjack.v2"
)

func getHook(filename string, maxAge, rotationTime int) (io.Writer, error) {
	hook := &lumberjack.Logger{
		Filename:   filename, // 日志文件名
		MaxSize:    100,      // megabytes
		MaxBackups: 1,
		MaxAge:     maxAge, // days
		Compress:   false,  // disabled by default
	}

	return hook, nil
}
