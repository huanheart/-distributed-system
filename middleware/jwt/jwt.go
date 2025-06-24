package jwt

import (
	"MyChat/common/code"
	"MyChat/controller"
	"MyChat/utils/myjwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 从头中读取jwt
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		res := new(controller.Response)

		// 从 Header 中读取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		// 截取 token 字符串（去掉 Bearer 前缀）
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析 token
		claimsId, ok := myjwt.ParseToken(token)
		if !ok {
			c.JSON(http.StatusOK, res.CodeOf(code.CodeInvalidToken))
			c.Abort()
			return
		}

		c.Set("user_id", claimsId)
		c.Next()
	}
}
