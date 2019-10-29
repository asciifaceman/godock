package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"

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
	cert     []byte
	key      []byte
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

// DetectTLS checks if the TLS cert/key exists
// and ingests them if they do
func (s *Server) DetectTLS() {
	crt := os.Getenv("GODOCKCRT")
	if crt != "" {
		dat, err := ioutil.ReadFile(crt)
		if err == nil {
			s.cert = dat
		}
	}
	key := os.Getenv("GODOCKKEY")
	if key != "" {
		dat, err := ioutil.ReadFile(key)
		if err == nil {
			s.key = dat
		}
	}

}

// DetectPortOverride detects if a port env var is set to override
func (s *Server) DetectPortOverride() {
	port := os.Getenv("PORT")
	if port != "" {
		portInt, err := strconv.Atoi(port)
		if err != nil {
			s.logger.Error("Failed to set custom listening port",
				zap.String("Detected", port),
				zap.Error(err),
			)
			return
		}
		s.port = portInt
	}
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

	s.DetectPortOverride()

	var withTLS bool
	s.DetectTLS()

	if len(s.cert) > 0 || len(s.key) > 0 {
		withTLS = true
	}

	s.logger.Info("API Starting",
		zap.String("hostname", s.host),
		zap.Int("Port", s.port),
		zap.Bool("WithTLS", withTLS),
	)
	go func() {
		if len(s.cert) == 0 || len(s.key) == 0 {
			panic(s.server.Start(s.getBindString()))
		} else {
			panic(s.server.StartTLS(s.getBindString(), s.cert, s.key))
		}

	}()
	<-stop
	s.logger.Info("shutting down...")
	s.server.Shutdown(context.Background())
	return nil
}
