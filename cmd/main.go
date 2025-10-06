package main

import (
	"hona/backend/internal/presentation/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()

	ginEngine := gin.Default()

	routes.SetUpRoutes(ginEngine)

	ginEngine.Run()
}
