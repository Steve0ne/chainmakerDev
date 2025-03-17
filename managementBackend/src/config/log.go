/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package config

// log 整体配置
type Log struct {
	System *SystemConf `yaml:"system"`
	Brief  *BriefConf  `yaml:"brief"`
	Event  *EventConf  `yaml:"event"`
}

type SystemConf struct {
	LogLevelDefault string        `yaml:"log_level_default"`
	LogLevels       *LogLevelConf `yaml:"log_levels"`
	FilePath        string        `yaml:"file_path"`
	MaxAge          int           `yaml:"max_age"`
	RotationTime    int           `yaml:"rotation_time"`
	LogInConsole    bool          `yaml:"log_in_console"`
	ShowColor       bool          `yaml:"show_color"`
}

type LogLevelConf struct {
	Core string `yaml:"core"`
	Net  string `yaml:"net"`
}

type BriefConf struct {
	LogLevelDefault string `yaml:"log_level_default"`
	FilePath        string `yaml:"file_path"`
	MaxAge          int    `yaml:"max_age"`
	RotationTime    int    `yaml:"rotation_time"`
	LogInConsole    bool   `yaml:"log_in_console"`
	ShowColor       bool   `yaml:"show_color"`
}

type EventConf struct {
	LogLevelDefault string `yaml:"log_level_default"`
	FilePath        string `yaml:"file_path"`
	MaxAge          int    `yaml:"max_age"`
	RotationTime    int    `yaml:"rotation_time"`
	LogInConsole    bool   `yaml:"log_in_console"`
	ShowColor       bool   `yaml:"show_color"`
}
