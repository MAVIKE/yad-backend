package app

import (
	handler "github.com/MAVIKE/yad-backend/internal/delivery/http"
	"github.com/MAVIKE/yad-backend/internal/repository"
	"github.com/MAVIKE/yad-backend/internal/service"
	
	"github.com/labstack/echo/v4"
	"log"
)

func Run(configPath string) {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	app := echo.New()
	handlers.InitRoutes(app)

	if err := app.Start(":8000"); err != nil {
		log.Fatalf("failed to listen: %s", err.Error())
	}
}
