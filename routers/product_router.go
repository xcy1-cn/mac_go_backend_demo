package routers

import (
	"demo/day6-9/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine) {
	r.GET("/products", handlers.GetProducts)
	r.GET("/products/:id", handlers.GetProductByID)
	r.POST("/products", handlers.AddProduct)
	r.PUT("/products/:id", handlers.UpdateProductByID)
	r.DELETE("/products/:id", handlers.DeleteProductByID)
}
