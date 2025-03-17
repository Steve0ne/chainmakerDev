/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package connection

import (
	"crypto/md5"
	"fmt"
	"time"

	"management_backend/src/config"
	"management_backend/src/db/common"
	"management_backend/src/utils"

	"github.com/jinzhu/gorm"

	// 必须要添加，解决找不到mysql驱动问题
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	ADMIN      = "admin"
	USER_START = 0
	SaltLength = 32
)

const (
	START    = 0
	NO_START = 1
	NO_WORK  = 2
)

var (
	DB *gorm.DB
)

// InitDbConn init database connection
func InitDbConn(dbConfig *config.DBConf) {
	var err error
	DB, err = gorm.Open(config.MySql, dbConfig.ToUrl())
	if err != nil {
		panic(err)
	} else {
		DB.DB().SetMaxIdleConns(config.DbMaxIdleConns)
		DB.DB().SetMaxOpenConns(config.DbMaxOpenConns)
		DB.DB().SetConnMaxLifetime(time.Minute)
		DB.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false)
		DB.SingularTable(true)
		DB.LogMode(true)
	}
	InitDB(DB)
}

func InitDB(engine *gorm.DB) {
	salt := utils.RandomString(SaltLength)
	data := []byte(config.GlobalConfig.WebConf.Password)
	passwdHexString := ToPasswdHexString(salt, fmt.Sprintf("%x", md5.Sum(data)))
	dbUser := &common.User{
		UserName: ADMIN,
		Name:     ADMIN,
		Salt:     salt,
		Passwd:   passwdHexString,
		Status:   USER_START,
	}

	err := engine.AutoMigrate(
		new(common.Block),
		new(common.Transaction),
		new(common.Cert),
		new(common.ChainOrg),
		new(common.Contract),
		new(common.Node),
		new(common.ChainOrgNode),
		new(common.Chain),
		new(common.Upload),
		new(common.ChainPolicy),
		new(common.User),
		new(common.ChainPolicyOrg),
		new(common.InvokeRecords),
		new(common.Org),
		new(common.OrgNode),
		new(common.VoteManagement),
		new(common.ChainConfig),
		new(common.ChainErrorLog),
	).Create(dbUser).Error
	if err != nil {
		fmt.Println(err.Error())
	}
	chains, err := getStartChainList()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, chain := range chains {
		columns := make(map[string]interface{})
		columns["status"] = NO_WORK
		if err := DB.Debug().Model(chain).Where("chain_id = ?", chain.ChainId).
			UpdateColumns(columns).Error; err != nil {
			fmt.Println(err.Error())
		}
	}
}

func getStartChainList() ([]*common.Chain, error) {
	var chains []*common.Chain
	if err := DB.Where("Status = ?", START).Find(&chains).Error; err != nil {
		return nil, err
	}
	return chains, nil
}

func ToPasswdHexString(salt string, password string) string {
	passwd := salt + "-" + password
	return utils.Sha256HexString([]byte(passwd))
}
