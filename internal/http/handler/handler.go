package handler

import (
	"errors"

	"github.com/dv-net/dv-updater/internal/http/request"
	"github.com/dv-net/dv-updater/internal/http/response"
	"github.com/dv-net/dv-updater/internal/service"
	"github.com/dv-net/dv-updater/internal/service/package_manager"
	"github.com/dv-net/dv-updater/pkg/logger"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	services *service.Services
	logger   logger.Logger
}

func NewHandler(services *service.Services, logger logger.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) Init(api *fiber.App) {
	v1 := api.Group("api/v1")

	v1.Post("/update", h.updatePackage)
	v1.Get("/version/:name", h.getLastVersionPackage)
	v1.Get("/version", h.getUpdaterVersion)
}

func (h *Handler) updatePackage(c fiber.Ctx) error {
	req := new(request.UpdatePackageRequest)
	if err := c.Bind().Body(req); err != nil {
		return err
	}

	err := h.services.PackageManager.UpgradePackage(c.Context(), req.Name)
	if err != nil {
		return c.JSON(response.Fail(fiber.StatusInternalServerError, err.Error()))
	}

	return c.JSON(response.OkByMessage("Success update package"))
}

func (h *Handler) getLastVersionPackage(c fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return c.JSON(response.Fail(fiber.StatusBadRequest, "name is empty"))
	}

	if err := service.ValidateServiceName(name); err != nil {
		return c.JSON(response.Fail(fiber.StatusBadRequest, "name is invalid"))
	}

	pkg, err := h.services.PackageManager.CheckForUpdates(c.Context(), name)
	if err != nil && !errors.Is(err, package_manager.ErrNothingToUpdate) {
		return c.JSON(response.Fail(fiber.StatusInternalServerError, err.Error()))
	}
	return c.JSON(response.OkByData(pkg))
}

func (h *Handler) getUpdaterVersion(c fiber.Ctx) error {
	return c.JSON(response.OkByData(h.services.SystemInfoService.GetSystemInfo()))
}
