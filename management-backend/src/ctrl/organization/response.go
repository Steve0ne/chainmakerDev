/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package organization

import (
	"management_backend/src/db/relation"
)

type OrgListWithNodeNumView struct {
	Id         int
	OrgName    string
	OrgId      string
	NodeNum    int
	CreateTime int64
}

func NewOrgListWithNodeNumView(org *relation.OrgListWithNodeNum) *OrgListWithNodeNumView {
	return &OrgListWithNodeNumView{
		Id:         org.Id,
		OrgName:    org.OrgName,
		OrgId:      org.OrgId,
		NodeNum:    org.NodeNum,
		CreateTime: org.CreateAt.Unix(),
	}
}
