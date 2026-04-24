package middlewares

import (
	"fmt"
	"strings"

	"demo/day6-9/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 取请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, 401, "authorization header is required")
			c.Abort()
			return
		}

		// 拆分 Bearer 和 token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Error(c, 401, "invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			utils.Error(c, 401, "invalid or expired token")
			c.Abort()
			return
		}

		fmt.Println("parsed userId:", claims.UserID, "username:", claims.Username)

		// 把用户信息放进上下文
		c.Set("userId", claims.UserID)
		c.Set("username", claims.Username)

		// 后续任何 handler 都可以直接取：userId, _ := c.Get("userId")
		c.Next()
	}
}
