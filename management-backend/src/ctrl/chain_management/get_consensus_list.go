/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package chain_management

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"chainmaker.org/chainmaker/pb-go/v2/consensus"

	"management_backend/src/ctrl/common"
	"management_backend/src/entity"
)

type GetConsensusListHandler struct{}

func (getConsensusListHandler *GetConsensusListHandler) LoginVerify() bool {
	return true
}

func (getConsensusListHandler *GetConsensusListHandler) Handle(user *entity.User, ctx *gin.Context) {
	consensusListView := NewConsensusListView()
	common.ConvergeListResponse(ctx, consensusListView, int64(len(consensusListView)), nil)
}

var consensusTypeList = arraylist.New("SOLO", "TBFT", "RAFT", "HOTSTUFF")

type ConsensusListView struct {
	ConsensusName string
	ConsensusType int32
}

func NewConsensusListView() []interface{} {
	consensusViews := arraylist.New()
	for consensusName, consensusType := range consensus.ConsensusType_value {
		consensusView := ConsensusListView{
			ConsensusName: consensusName,
			ConsensusType: consensusType,
		}
		if consensusTypeList.Contains(consensusView.ConsensusName) {
			consensusViews.Add(consensusView)
		}
	}
	return consensusViews.Values()
}
