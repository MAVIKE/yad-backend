package handler

import (
	v1 "github.com/MAVIKE/yad-backend/internal/delivery/http/v1"
	"github.com/MAVIKE/yad-backend/internal/service"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"

	_ "github.com/MAVIKE/yad-backend/docs/swagger"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// @title Yet Another Delivery API
// @version 1.0
// @description API Server for Yet Another Delivery App

// @host localhost:9000
// @BasePath /api/v1/

func (h *Handler) Init(router *echo.Echo) {
	router.Use(middleware.Logger())
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)
}

func (h *Handler) initAPI(router *echo.Echo) {
	handlerV1 := v1.NewHandler(h.services)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
