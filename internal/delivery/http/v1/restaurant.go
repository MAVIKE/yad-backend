package v1

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func (h *Handler) initRestaurantRoutes(api *echo.Group) {
	restaurants := api.Group("/restaurants")
	{
		restaurants.Use(h.identity)
		restaurants.GET("", h.getRestaurants)
	}
}

func (h *Handler) getRestaurants(ctx echo.Context) error {
	userId, clientType, _ := getClientParams(ctx)
	fmt.Println("params:", userId, clientType)
	return nil
}
