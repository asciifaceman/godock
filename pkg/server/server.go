package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/asciifaceman/godock/pkg/middleware"

	"github.com/asciifaceman/godock/pkg/handlers"
	"github.com/asciifaceman/godock/pkg/handlers/health"
	echo "github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultHost = "0.0.0.0"
	defaultPort = 8085
)

// Server ...
type Server struct {
	host     string
	port     int
	server   *echo.Echo
	logger   *zap.Logger
	handlers []handlers.Handler
}

// NewServer ...
func NewServer() (*Server, error) {

	zconf := zap.NewProductionConfig()
	zconf.OutputPaths = []string{"stdout"}
	zconf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Human-reaable time ISO
	conf, err := zconf.Build()
	if err != nil {
		return nil, err
	}

	srv := echo.New()
	srv.HideBanner = true

	srv.Use(middleware.LoggingMiddleware(conf))
	srv.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	srv.Use(echoMiddleware.RequestID())

	return &Server{
		logger:   conf,
		server:   srv,
		host:     defaultHost,
		port:     defaultPort,
		handlers: make([]handlers.Handler, 0),
	}, nil

}

// getBindString returns formatted host:port
func (s *Server) getBindString() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}

// RegisterHandlers registers handlers in bulk
func (s *Server) RegisterHandlers() error {
	hh := health.NewHealthHandler(s.logger)

	middlewares := make([]echo.MiddlewareFunc, 0)

	for _, route := range hh.Routes() {
		s.logger.Info(fmt.Sprintf("Mounting route: %s=>%s", route.GetMethod(), route.GetPath()))
		s.server.Add(route.GetMethod(), route.GetPath(), route.GetHandler(), middlewares...)
	}

	s.handlers = append(s.handlers, hh)

	return nil
}

// Run ...
func (s *Server) Run() error {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	s.logger.Info("API Starting",
		zap.String("hostname", s.host),
		zap.Int("Port", s.port),
	)
	go func() {
		panic(s.server.Start(s.getBindString()))
	}()
	<-stop
	s.logger.Info("shutting down...")
	s.server.Shutdown(context.Background())
	return nil
}
