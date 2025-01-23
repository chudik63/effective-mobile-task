package server

import (
	"context"
	"effective-mobile-task/internal/config"
	"effective-mobile-task/pkg/logger"
	"fmt"
	"net/http"
	"time"
)

const shutdownTimeout = 5 * time.Second

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    ":" + cfg.ServerPort,
			Handler: handler,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	logs := logger.GetLoggerFromCtx(ctx)

	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	logs.Info(ctx, fmt.Sprintf("Server listening on %s", s.httpServer.Addr))

	return nil
}

func (s *Server) Stop() error {
	ctx, shutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
