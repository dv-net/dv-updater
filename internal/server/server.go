package server

import (
	"github.com/dv-net/dv-updater/internal/config"
	"github.com/dv-net/dv-updater/internal/router"
	"github.com/dv-net/dv-updater/internal/service"
	"github.com/dv-net/dv-updater/pkg/logger"

	"github.com/gofiber/fiber/v3"
)

type Server struct {
	app    *fiber.App
	cfg    config.HTTPConfig
	logger logger.Logger
}

func NewServer(cfg config.HTTPConfig, services *service.Services, logger logger.Logger) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	router.NewRouter(cfg, services, logger).Init(app)

	return &Server{
		app:    app,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *Server) Run() error {
	return s.app.Listen(":"+s.cfg.Port, fiber.ListenConfig{
		DisableStartupMessage: true,
	})
}

func (s *Server) Stop() error {
	return s.app.Shutdown()
}
