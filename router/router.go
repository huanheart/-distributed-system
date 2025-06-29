package router

import (
	"MyChat/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {

	r := gin.Default()
	enterRouter := r.Group("/api/v1")
	{
		RegisterUserRouter(enterRouter.Group("/user"))
	}
	//由于音乐功能是登录之后的功能，固然需要对jwt进行一个校验,注册一个中间件机制
	// 注册 /music 路由组，并添加 JWT 鉴权中间件
	{
		musicGroup := enterRouter.Group("/music")
		musicGroup.Use(jwt.Auth())
		MusicRouter(musicGroup)
	}

	{
		loginGroup := enterRouter.Group("/login")
		loginGroup.Use(jwt.Auth())
		AfterUserRouter(loginGroup)
	}
	return r
}
