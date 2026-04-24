package routers

import (
	"demo/day6-9/handlers"
	"demo/day6-9/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
)

func RegisterCartRoutes(r *gin.Engine) {
	fmt.Println("=== RegisterCartRoutes called ===")
	authGroup := r.Group("/")
	authGroup.Use(middlewares.AuthMiddleware())
	{
		authGroup.POST("/cart/add", handlers.AddCart)
		authGroup.GET("/cart", handlers.GetCartList)
		authGroup.DELETE("/cart/:id", handlers.DeleteCartByID)
		authGroup.PUT("/cart/:id", handlers.UpdateCartByID)
		authGroup.GET("/cart/summary", handlers.GetCartSummary)
	}
}
