package main

import (
	"context"
	"effective-mobile-task/internal/config"
	"effective-mobile-task/internal/database/postgres"
	"effective-mobile-task/internal/repository"
	"effective-mobile-task/internal/server"
	"effective-mobile-task/internal/service"
	"effective-mobile-task/internal/transport/http"
	"effective-mobile-task/internal/transport/http/middleware"
	"effective-mobile-task/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	serviceName = "effective-mobile-task"
)

func main() {
	mainLogger, err := logger.New(serviceName)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %v", err))
	}

	ctx := context.WithValue(context.Background(), logger.LoggerKey, mainLogger)

	cfg, err := config.Load()
	if err != nil {
		mainLogger.Fatal(ctx, "failed to read config", zap.String("err", err.Error()))
	}

	db := postgres.New(ctx, cfg.Config)

	repo := repository.New(db)
	service := service.New(repo)

	r := gin.Default()

	r.Use(middleware.LogMiddleware(mainLogger))

	c := http.NewClient(cfg.ServerHost, cfg.ServerPort)
	http.NewHandler(r, c, service, mainLogger)

	srv := server.NewServer(cfg, r.Handler())

	go func() {
		if err := srv.Run(ctx); err != nil {
			mainLogger.Fatal(ctx, "failed to start server", zap.String("err", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	if err := srv.Stop(); err != nil {
		mainLogger.Error(ctx, "failed to stop server", zap.String("err", err.Error()))
	}

	mainLogger.Info(ctx, "server stopped")
}
