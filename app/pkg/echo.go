package pkg

import (
	"context"
	"fmt"

	"control-mitsubishi-plc-w-kube/config"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type server struct {
	Server  *echo.Echo
	Context context.Context
	Port    string
}

type Server interface {
	Start(errC chan error)
	Shutdown(ctx context.Context) error
}

func New(ctx context.Context, cfg *config.Config) Server {
	// Echo instance
	e := echo.New()
	// Routes
	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// use echo default logger
	e.Use(middleware.Logger())
	return &server{
		Server:  e,
		Context: ctx,
		Port:    fmt.Sprintf(":%v", cfg.Server.Port),
	}
}

func (s *server) Start(errC chan error) {
	err := s.Server.Start(s.Port)
	errC <- err
}

func (s *server) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
