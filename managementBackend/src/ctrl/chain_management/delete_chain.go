package chain_management

import (
	"github.com/gin-gonic/gin"
	"management_backend/src/ctrl/common"
	dbchain "management_backend/src/db/chain"
	"management_backend/src/entity"
)

type DeleteChainHandler struct{}

func (deleteChainHandler *DeleteChainHandler) LoginVerify() bool {
	return true
}

func (deleteChainHandler *DeleteChainHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindDeleteChainHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	err := dbchain.DeleteChain(params.ChainId)
	if err != nil {
		log.Error("DeleteChain err : " + err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorDeleteChain)
		return
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
