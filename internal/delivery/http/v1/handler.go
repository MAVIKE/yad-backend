package v1

import (
	"github.com/MAVIKE/yad-backend/internal/service"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init(api *echo.Group) {
	v1 := api.Group("/v1")
	{
		v1.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})

		h.initAdminRoutes(v1)
		h.initUserRoutes(v1)
		h.initCourierRoutes(v1)
	}
}
