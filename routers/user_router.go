package routers

import (
	"demo/day6-9/handlers"
	"demo/day6-9/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/")
	{
		userGroup.POST("/register", handlers.Register)
		userGroup.POST("/login", handlers.Login)
	}

	authGroup := r.Group("/")
	authGroup.Use(middlewares.AuthMiddleware())
	{
		authGroup.GET("/user/info", handlers.GetUserInfo)
	}
}
