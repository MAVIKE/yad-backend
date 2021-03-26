package v1

import (
	"net/http"
	"strconv"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initOrderRoutes(api *echo.Group) {
	orders := api.Group("/orders")
	{
		orders.Use(h.identity)
		orders.POST("/", h.createOrder)
		orders.GET("/:oid", h.getOrderById)

		orderItems := orders.Group("/:oid/items")
		{
			orderItems.POST("/", h.createOrderItem)
			orderItems.GET("/:id", h.getOrderItemById)
		}
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

// @Summary Get Order By Id
// @Security UserAuth
// @Security RestaurantAuth
// @Security CourierAuth
// @Tags orders
// @Description get order by id
// @ModuleID getOrderById
// @Accept  json
// @Produce  json
// @Param oid path string true "Order id"
// @Success 200 {object} domain.Order
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /orders/{oid}/ [get]
func (h *Handler) getOrderById(ctx echo.Context) error {
	_, _, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	orderId, err := strconv.Atoi(ctx.Param("oid"))
	if err != nil || orderId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderId")
	}

	return ctx.JSON(http.StatusOK, 1)
}

type orderItemInput struct {
	MenuItemId int `json:"menu_item_id"`
	Count      int `json:"count" valid:"range(1|99)"`
}

// @Summary Add menu item to order
// @Security UserAuth
// @Tags orders
// @Description create order item
// @ModuleID createOrderItem
// @Accept  json
// @Produce  json
// @Param oid path string true "Order id"
// @Param input body orderItemInput true "order item create info"
// @Success 200 {object} idResponse
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /orders/{oid}/items/ [post]
func (h *Handler) createOrderItem(ctx echo.Context) error {
	var input orderItemInput
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	orderId, err := strconv.Atoi(ctx.Param("oid"))
	if err != nil || orderId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderId")
	}

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	orderItem := &domain.OrderItem{
		OrderId:    orderId,
		MenuItemId: input.MenuItemId,
		Count:      input.Count,
	}

	orderItemId, err := h.services.Order.CreateItem(clientId, clientType, orderItem)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, idResponse{
		Id: orderItemId,
	})
}

// @Summary Get Order Item By Id
// @Security UserAuth
// @Security RestaurantAuth
// @Security CourierAuth
// @Tags orders
// @Description get order item by id
// @ModuleID getOrderItemById
// @Accept  json
// @Produce  json
// @Param oid path string true "Order id"
// @Param id path string true "Order item id"
// @Success 200 {object} domain.OrderItem
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /orders/{oid}/items{id} [get]
func (h *Handler) getOrderItemById(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	orderId, err := strconv.Atoi(ctx.Param("oid"))
	if err != nil || orderId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	orderItem, err := h.services.Order.GetItemById(clientId, clientType, orderId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, orderItem)
}
