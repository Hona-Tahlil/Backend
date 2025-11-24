package main

import (
	"hona/backend/bootstrap"
	"hona/backend/internal/presentation/routes"
	"hona/backend/wire"

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

	// app.Seeder.DatabaseSeeder.SeedAll()

	routes.SetUpRoutes(ginEngine, app)

	ginEngine.Run()
}

// TODO: use len and index instead of append everywhere!
