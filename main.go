package main

import (
	"demo/day6-9/config"
	"demo/day6-9/db"
	"demo/day6-9/middlewares"
	"demo/day6-9/routers"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("=== CURRENT SERVER STARTED ===")

	config.InitEnv()
	db.InitMySQL()

	r := gin.Default()
	r.Use(middlewares.LoggerMiddleware())

	routers.RegisterProductRoutes(r)
	routers.RegisterCartRoutes(r)
	routers.RegisterOrderRoutes(r)
	routers.RegisterUserRoutes(r)

	port := config.GetEnv("SERVER_PORT")
	if port == "" {
		port = "8077"
	}

	r.Run(":" + port)
}
