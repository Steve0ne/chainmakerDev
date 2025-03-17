/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package ca

import (
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db"
	"management_backend/src/entity"
	"management_backend/src/utils"
)

const (
	FileName               = "File"
	DownloadKeyArrayLength = 3
	DownloadIdIdx          = 0
	DownloadUserIdIdx      = 1
	DownloadHashIdx        = 2
	SaveKeySeparation      = "."
)

type UploadHandler struct{}

func (uploadHandler *UploadHandler) LoginVerify() bool {
	return true
}

func (uploadHandler *UploadHandler) Handle(user *entity.User, ctx *gin.Context) {
	file, err := ctx.FormFile(FileName)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	fileName := file.Filename
	innerFile, err := file.Open()
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	bytes, err := ioutil.ReadAll(innerFile)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	fileKey, err := saveBytesToDb(user, fileName, bytes)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	common.ConvergeDataResponse(ctx, common.NewUploadResponse(fileKey), nil)
}

func saveBytesToDb(user *entity.User, fileName string, bytes []byte) (string, error) {
	hash, err := utils.Sha256(bytes)
	if err != nil {
		return "", err
	}
	hashHex := hex.EncodeToString(hash)
	upload := db.NewUpload(user.Id, hashHex, fileName, bytes)
	id, err := db.CreateUpload(upload)
	if err != nil {
		return "", err
	}
	return ToUploadKey(id, user.Id, hashHex), nil
}

func ToUploadKey(id, userId int64, hash string) string {
	return strconv.FormatInt(id, 10) + SaveKeySeparation + strconv.FormatInt(userId, 10) + SaveKeySeparation + hash
}

func ResolveUploadKey(key string) (id, userId int64, hash string, err error) {
	splitText := strings.Split(key, SaveKeySeparation)
	if len(splitText) != DownloadKeyArrayLength {
		err = errors.New("param num is wrong")
		return
	}
	id, err = strconv.ParseInt(splitText[DownloadIdIdx], 10, 64)
	if err != nil {
		err = errors.New("parse id wrong")
		return
	}
	userId, err = strconv.ParseInt(splitText[DownloadUserIdIdx], 10, 64)
	if err != nil {
		err = errors.New("parse userId wrong")
		return
	}
	hash = splitText[DownloadHashIdx]
	if len(hash) == 0 {
		err = errors.New("parse hash wrong")
		return
	}
	return
}
