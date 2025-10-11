package internalhttp

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
	logger Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

type ServerDeps struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewServer(logger Logger, app Application, cfg ServerDeps) *Server {
	addr := net.JoinHostPort(cfg.Host, cfg.Port)
	handler := NewHandler(app, logger)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler.Routes(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return &Server{
		server: server,
		logger: logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error("failed server stop")
		return err
	}
	s.logger.Info("server stopped")
	return nil
}

// TODO
