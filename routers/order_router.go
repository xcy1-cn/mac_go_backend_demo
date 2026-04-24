package routers

import (
	"demo/day6-9/handlers"
	"demo/day6-9/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(r *gin.Engine) {
	authGroup := r.Group("/")
	authGroup.Use(middlewares.AuthMiddleware())
	{
		authGroup.POST("/order/create", handlers.CreateOrder)
		authGroup.GET("/orders", handlers.GetOrders)
		authGroup.GET("/orders/:id", handlers.GetOrderByID)
	}
}
