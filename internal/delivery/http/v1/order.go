package v1

import (
	"net/http"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initOrderRoutes(api *echo.Group) {
	orders := api.Group("/orders")
	{
		orders.Use(h.identity)
		orders.POST("/", h.createOrder)
	}
}

type orderInput struct {
	RestaurantId int `json:"restaurant_id"`
}

// @Summary Create Order
// @Security UserAuth
// @Tags orders
// @Description create order
// @ModuleID createOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} idResponse
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /orders/ [post]
func (h *Handler) createOrder(ctx echo.Context) error {
	var input orderInput
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	order := &domain.Order{
		RestaurantId: input.RestaurantId,
	}

	orderId, err := h.services.Order.Create(clientId, clientType, order)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, idResponse{
		Id: orderId,
	})
}
