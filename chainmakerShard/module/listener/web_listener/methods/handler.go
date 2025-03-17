/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/

package methods

import (
	"github.com/gin-gonic/gin"
)

// ContextHandler handler of context
type ContextHandler interface {

	// Handle handle the context
	Handle(ctx *gin.Context)
}
