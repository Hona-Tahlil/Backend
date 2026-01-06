package main

import (
	"context"
	"errors"
	"hona/backend/bootstrap"
	"hona/backend/internal/infrastructure/logger"
	"hona/backend/internal/infrastructure/websocket"
	"hona/backend/internal/presentation/routes"
	"hona/backend/wire"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	config := bootstrap.Run()
	logger.Init()
	defer logger.Close()

	gin.DisableConsoleColor()

	gin.SetMode(gin.ReleaseMode)

	ginEngine := gin.New()
	ginEngine.RedirectTrailingSlash = true
	ginEngine.RemoveExtraSlash = true

	app, err := wire.InitializeApplication(config)
	if err != nil {
		panic(err)
	}
	hub := websocket.NewHub()
	go hub.Run()

	routes.SetUpRoutes(ginEngine, app)

	if err := app.Consumers.EmailConsumer.Start(); err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    ":8080",
		Handler: app.Middlewares.CORSMiddleware.Handler(ginEngine),
	}

	go func() {
		slog.Info("server running", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown failed", "error", err)
	}
	slog.Info("server exiting")
}
