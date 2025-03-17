/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package logger

import (
	"io"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LOG_LEVEL the level of log
type LOG_LEVEL int

const (
	LEVEL_DEBUG LOG_LEVEL = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
)

// 日志级别，配置文件定义的常量
const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	WARN  = "WARN"
	ERROR = "ERROR"
)

// GetLogLevel return LOG_LEVEL by string
func GetLogLevel(lvl string) LOG_LEVEL {
	switch lvl {
	case ERROR:
		return LEVEL_ERROR
	case WARN:
		return LEVEL_WARN
	case INFO:
		return LEVEL_INFO
	case DEBUG:
		return LEVEL_DEBUG
	}

	return LEVEL_DEBUG
}

// 日志切割默认配置
const (
	DEFAULT_MAX_AGE       = 365 // 日志最长保存时间，单位：天
	DEFAULT_ROTATION_TIME = 6   // 日志滚动间隔，单位：小时
)

//type color int
//
//const (
//	ColorBlack color = iota + 30
//	ColorRed
//	ColorGreen
//	ColorYellow
//	ColorBlue
//	ColorMagenta
//	ColorCyan
//	ColorWhite
//)
//
//var colorList = [...]color{ColorRed, ColorGreen, ColorYellow, ColorBlue, ColorMagenta}

var hookMap = make(map[string]struct{})

// Config is config of logger print
type Config struct {
	Module       string    // module: module name
	LogPath      string    // logPath: log file save path
	LogLevel     LOG_LEVEL // logLevel: log level
	MaxAge       int       // maxAge: the maximum number of days to retain old log files
	RotationTime int       // RotationTime: rotation time
	JsonFormat   bool      // jsonFormat: log file use json format
	ShowLine     bool      // showLine: show filename and line number
	LogInConsole bool      // logInConsole: show logs in console at the same time
	ShowColor    bool      // if true, show color log
}

// InitSugarLogger init and create SugaredLogger by config
func InitSugarLogger(loggerConfig *Config) (*zap.SugaredLogger, zap.AtomicLevel) {
	var level zapcore.Level
	switch loggerConfig.LogLevel {
	case LEVEL_DEBUG:
		level = zap.DebugLevel
	case LEVEL_INFO:
		level = zap.InfoLevel
	case LEVEL_WARN:
		level = zap.WarnLevel
	case LEVEL_ERROR:
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	aLevel := zap.NewAtomicLevel()
	aLevel.SetLevel(level)

	sugaredLogger := newLogger(loggerConfig, aLevel).Sugar()

	return sugaredLogger, aLevel
}

func newLogger(loggerConfig *Config, level zap.AtomicLevel) *zap.Logger {
	var (
		hook io.Writer
		ok   bool
		err  error
	)

	_, ok = hookMap[loggerConfig.LogPath]
	if !ok {
		hook, err = getHook(loggerConfig.LogPath, loggerConfig.MaxAge, loggerConfig.RotationTime)
		if err != nil {
			log.Fatalf("new logger get hook failed, %s", err)
		}
		hookMap[loggerConfig.LogPath] = struct{}{}
	} else {
		hook, err = getHook(loggerConfig.LogPath, loggerConfig.MaxAge, 0)
		if err != nil {
			log.Fatalf("new logger get hook failed, %s", err)
		}
	}

	var syncer zapcore.WriteSyncer

	if loggerConfig.LogInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook))
	} else {
		syncer = zapcore.AddSync(hook)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    CustomLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	var encoder zapcore.Encoder
	if loggerConfig.JsonFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		syncer,
		level,
	)

	logger := zap.New(core).Named(loggerConfig.Module)
	defer logger.Sync()

	if loggerConfig.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return logger
}

// CustomLevelEncoder
func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// CustomTimeEncoder
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
