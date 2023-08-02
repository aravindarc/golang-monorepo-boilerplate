package health

import (
	"github.com/labstack/echo/v4"
	"golang-monorepo-boilerplate/core/config"
	"golang-monorepo-boilerplate/core/log"
)

type (
	handlerDependencies interface {
		log.Provider
		config.Provider
	}
	HandlerProvider interface {
		HealthHandler() *Handler
	}
	Handler struct {
		d handlerDependencies
	}
)

func NewHandler(d handlerDependencies) *Handler {
	return &Handler{d: d}
}

func (h *Handler) RegisterRoutes(r *echo.Group) {
	r.GET("/v1/health", h.health)
}

func (h *Handler) health(c echo.Context) error {
	return c.String(200, "OK")
}
