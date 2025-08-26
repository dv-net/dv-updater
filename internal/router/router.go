package router

import (
	"github.com/dv-net/dv-updater/internal/config"
	"github.com/dv-net/dv-updater/internal/http/handler"
	"github.com/dv-net/dv-updater/internal/service"
	"github.com/dv-net/dv-updater/pkg/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
)

type Router struct {
	config   config.HTTPConfig
	services *service.Services
	logger   logger.Logger
}

func NewRouter(conf config.HTTPConfig, services *service.Services, logger logger.Logger) *Router {
	return &Router{
		config:   conf,
		services: services,
		logger:   logger,
	}
}

func (r *Router) Init(app *fiber.App) {
	app.Use(etag.New())

	if r.config.Cors.Enabled {
		corsConfig := cors.ConfigDefault
		corsConfig.AllowMethods = []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		}

		if len(r.config.Cors.AllowedOrigins) > 0 {
			corsConfig.AllowOrigins = r.config.Cors.AllowedOrigins
		}

		app.Use(cors.New(corsConfig))
	}

	app.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("pong")
	})
	r.initAPI(app)
}

func (r *Router) initAPI(app *fiber.App) {
	handlerV1 := handler.NewHandler(r.services, r.logger)
	handlerV1.Init(app)
}
