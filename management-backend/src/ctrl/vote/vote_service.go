/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package vote

import (
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
	"management_backend/src/db/chain_participant"
	"management_backend/src/db/connection"

	"management_backend/src/ctrl/common"
	"management_backend/src/ctrl/multi_sign"
	"management_backend/src/db"
	dbcommon "management_backend/src/db/common"
	dbpolicy "management_backend/src/db/policy"
	"management_backend/src/entity"
)

// VoteHandler 投票接口
type VoteHandler struct {
}

func (handler *VoteHandler) LoginVerify() bool {
	return true
}

func (handler *VoteHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindVoteHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}
	var (
		vote   *dbcommon.VoteManagement
		policy *dbcommon.ChainPolicy
		err    error
	)
	vote, err = db.GetVoteManagementById(params.VoteId)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorGetVoteManagement)
		return
	}

	if vote.VoteStatus == 1 {
		common.ConvergeFailureResponse(ctx, common.ErrorAlreadyOnChain)
		return
	}

	// 判断是否满足链策略上配置的权限，如果满足，广播上链
	policy, err = dbpolicy.GetChainPolicy(vote.ChainId, vote.VoteType)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorGetChainPolicy)
		return
	}

	var roleType int
	if policy.RoleType == multi_sign.POLICY_CLIENT {
		roleType = chain_participant.CLIENT
	} else {
		roleType = chain_participant.ADMIN
	}

	_, err = chain_participant.GetUserCertByOrgId(vote.VoteOrgId, roleType)
	if err != nil {
		if roleType == chain_participant.CLIENT {
			common.ConvergeFailureResponse(ctx, common.ErrorGetOrgClientUser)
		} else {
			common.ConvergeFailureResponse(ctx, common.ErrorGetOrgAdminUser)
		}

		return
	}

	if vote.VoteResult == 0 {
		vote.VoteResult = params.VoteResult
		connection.DB.Save(&vote)
	} else {
		common.ConvergeFailureResponse(ctx, common.ErrorAlreadyVoted)
		return
	}

	passOrgs, notPassOrgs, err := db.GetVotedOrgListByMultiId(vote.MultiId)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	passedCnt := len(passOrgs)
	total := passedCnt + len(notPassOrgs)
	needPassedCnt := dbpolicy.GetPassedVoteCnt(total, policy)
	// 禁止调用的操作
	if needPassedCnt == 0 {
		common.ConvergeFailureResponse(ctx, common.ErrorForbiddenPolicy)
		return
	}
	// 不正确的链策略，比如分数没设对
	if needPassedCnt < 0 {
		common.ConvergeFailureResponse(ctx, common.ErrorChainPolicy)
		return
	}
	// !!! 没有去判断Self的权限，因为self很特殊，产品原型已禁掉可以设置为self的选项 !!!

	if passedCnt >= needPassedCnt {
		//	广播上链
		err = db.SetMultiIdVotedCompleted(vote.MultiId)
		if err != nil {
			common.ConvergeHandleFailureResponse(ctx, err)
			return
		}
		err := multi_sign.MultiSignInvoke(vote.Params, vote.VoteType, passOrgs, policy.RoleType, vote.ConfigStatus)
		if err != nil {
			common.ConvergeHandleFailureResponse(ctx, err)
			return
		}

	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

// GetVoteManageHandler 查询投票列表
type GetVoteManageListHandler struct {
}

func (handler *GetVoteManageListHandler) LoginVerify() bool {
	return true
}

func (handler *GetVoteManageListHandler) Handle(user *entity.User, ctx *gin.Context) {
	var (
		voteList   []*dbcommon.VoteManagement
		totalCount int64
		offset     int
		limit      int
	)
	params := BindGetVoteManageListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	offset = params.PageNum * params.PageSize
	limit = params.PageSize
	totalCount, voteList, err := db.GetVoteManagementList(offset, limit, params.ChainId,
		params.OrgId, params.VoteType, params.VoteStatus)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	txInfos := convertToVoteManageViews(voteList)
	common.ConvergeListResponse(ctx, txInfos, totalCount, nil)
}

func convertToVoteManageViews(voteList []*dbcommon.VoteManagement) []interface{} {
	views := arraylist.New()
	for _, vote := range voteList {
		view := NewVoteManagementView(vote)
		views.Add(view)
	}
	return views.Values()
}

// GetVoteDetailHandler 查询投票详情
type GetVoteDetailHandler struct {
}

func (handler *GetVoteDetailHandler) LoginVerify() bool {
	return true
}

func (handler *GetVoteDetailHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetVoteDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	var (
		vote *dbcommon.VoteManagement
		err  error
	)

	if params.VoteId != 0 {
		vote, err = db.GetVoteManagementById(params.VoteId)
		if err != nil {
			common.ConvergeHandleFailureResponse(ctx, err)
			return
		}

	}
	txView := NewVoteManagementView(vote)
	common.ConvergeDataResponse(ctx, txView, nil)
}
