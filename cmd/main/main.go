package main

import (
	"context"
	"effective-mobile-task/internal/config"
	"effective-mobile-task/internal/database/postgres"
	"effective-mobile-task/pkg/logger"
	"fmt"

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

	_ = db
}
