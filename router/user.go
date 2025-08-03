package router

import (
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
		r.POST("/query_like_infos", user.QueryLikeInfos)
	}
}
