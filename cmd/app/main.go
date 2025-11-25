package main

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/presentation/routes"
	"hona/backend/wire"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()

	gin.SetMode(gin.ReleaseMode)

	ginEngine := gin.Default()

	app, err := wire.InitializeApplication(bootstrap.Run())
	if err != nil {
		panic(err)
	}

	app.Seeder.DatabaseSeeder.SeedAll()

	ginEngine.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetUpRoutes(ginEngine, app)

	ginEngine.Run()

	// to push
}
