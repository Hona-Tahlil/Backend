package main

import (
	"context"
	"errors"
	"hona/backend/bootstrap"
	"hona/backend/internal/presentation/routes"
	"hona/backend/wire"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.DisableConsoleColor()

	gin.SetMode(gin.ReleaseMode)

	ginEngine := gin.Default()
	ginEngine.RedirectTrailingSlash = true
	ginEngine.RemoveExtraSlash = true

	app, err := wire.InitializeApplication(bootstrap.Run())
	if err != nil {
		panic(err)
	}

	routes.SetUpRoutes(ginEngine, app)

	if err := app.Consumers.EmailConsumer.Start(); err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Middlewares.CORSMiddleware.Handler(ginEngine),
	}

	go func() {
		log.Println("Server Running ...")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
