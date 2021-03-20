package v1

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initOrderRoutes(api *echo.Group) {
	orders := api.Group("/orders")
	{
		orders.Use(h.identity)
		orders.POST("", h.createOrder)
	}
}

type orderInput struct {
	RestaurantId int `json:"restaurant_id"`
}

func (h *Handler) createOrder(ctx echo.Context) error {
	var input orderInput
	_, _, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": 1,
	})
}
