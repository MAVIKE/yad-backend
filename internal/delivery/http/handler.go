package handler

import (
	"net/http"

	v1 "github.com/MAVIKE/yad-backend/internal/delivery/http/v1"
	"github.com/MAVIKE/yad-backend/internal/service"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/MAVIKE/yad-backend/docs/swagger"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	services     *service.Service
	tokenManager *auth.Manager
}

func NewHandler(services *service.Service, tokenManager *auth.Manager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

// @title Yet Another Delivery API
// @version 1.0
// @description API Server for Yet Another Delivery App

// @host localhost:9000
// @BasePath /api/v1/

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey UserAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey CourierAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey RestaurantAuth
// @in header
// @name Authorization

func (h *Handler) Init(router *echo.Echo) {
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)
}

func (h *Handler) initAPI(router *echo.Echo) {
	handlerV1 := v1.NewHandler(h.services, h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
