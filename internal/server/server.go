package server

import (
	"context"
	"effective-mobile-task/internal/config"
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

func (s *Server) Run() error {
	err := s.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop() error {
	ctx, shutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdown()

	return s.httpServer.Shutdown(ctx)
}
