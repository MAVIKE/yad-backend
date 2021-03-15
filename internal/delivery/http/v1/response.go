package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type response struct {
	Message string `json:"message"`
}

func newResponse(ctx echo.Context, statusCode int, message string) error {
	log.Error(message)
	return ctx.JSON(statusCode, response{message})
}