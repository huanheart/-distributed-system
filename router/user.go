package router

import (
	"MyChat/controller/music"
	"MyChat/controller/user"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.RouterGroup) {
	{
		r.POST("/register", user.Register)
		r.POST("/login", user.Login)
		r.POST("/captcha", user.HandleCaptcha)
	}
}

func AfterUserRouter(r *gin.RouterGroup) {
	{
		r.POST("/like", user.Like)
		r.GET("/rankings", music.Rankings)
	}
}
