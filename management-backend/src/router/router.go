/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package router

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"management_backend/src/config"
	"management_backend/src/ctrl"
	user_service "management_backend/src/ctrl/user"
	"management_backend/src/entity"
	loggers "management_backend/src/logger"
	"management_backend/src/session"
)

var (
	log = loggers.GetLogger(loggers.ModuleWeb)
)

func HttpServe(webConf *config.WebConf) {
	// 启动Web服务(默认Debug级别)
	gin.SetMode(gin.ReleaseMode)
	// 生成route
	ginRouter := gin.Default()
	// 初始化路由配置
	InitRouter(ginRouter, webConf)
	// 初始化Captcha
	user_service.InitCaptcha(webConf.CaptchaConf)
	// 启动Http服务
	err := ginRouter.Run(webConf.ToUrl())
	if err != nil {
		panic(err)
	}
}

func InitRouter(router *gin.Engine, webConf *config.WebConf) {
	// 处理跨域请求，安装nginx的情况下理论上不需要跨域
	if webConf.CrossDomain {
		log.Info("start cross domain processing...")
		router.Use(Cors())
	}
	store := session.NewSessionStore(webConf.SessionAge)
	router.Use(sessions.Sessions(session.SessionID, store))
	group := router.Group("/")
	initControllers(group) // 定义接口
}

// initControllers 初始化Controller配置
func initControllers(routeGroup *gin.RouterGroup) {
	routeGroup.POST(entity.Project, ctrl.Dispatcher)
}

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers",
				"Authorization, Content-Length, X-CSRF-Token, Token,session, Content-Type")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers",
				"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				panic(err)
			}
		}()
		c.Next()
	}
}
