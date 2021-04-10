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
			orderItems.GET("/", h.getOrderItems)
			orderItems.GET("/:id", h.getOrderItemById)
			orderItems.DELETE("/:id", h.deleteOrderItem)
			orderItems.PUT("/:id", h.updateOrderItem)
		}
	}
	restaurants := api.Group("/restaurants")
	{
		restaurants.Use(h.identity)
		restaurants.GET("/:rid/orders/", h.getActiveRestaurantOrders)
	}

	users := api.Group("/users")
	{
		users.Use(h.identity)
		users.GET("/:id/orders", h.usersGetAllOrders)
	}

	couriers := api.Group("/couriers")
	{
		couriers.Use(h.identity)
		couriers.GET("/:cid/orders", h.getActiveCourierOrder)
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
// @Param input body orderInput true "order input info"
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

// @Summary Get All Order Items By Order Id
// @Security UserAuth
// @Security RestaurantAuth
// @Security CourierAuth
// @Tags orders
// @Description get all order items by order id
// @ModuleID getOrderItems
// @Accept  json
// @Produce  json
// @Param oid path string true "Order id"
// @Success 200 {array} domain.OrderItem
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /orders/{oid}/items/ [get]
func (h *Handler) getOrderItems(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	orderId, err := strconv.Atoi(ctx.Param("oid"))
	if err != nil || orderId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderId")
	}

	orderItems, err := h.services.Order.GetAllItems(clientId, clientType, orderId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, orderItems)
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
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	orderId, err := strconv.Atoi(ctx.Param("oid"))
	if err != nil || orderId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderId")
	}

	orderItem, err := h.services.Order.GetById(clientId, clientType, orderId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, orderItem)
}

type orderItemInput struct {
	MenuItemId int `json:"menu_item_id"`
	Count      int `json:"count" valid:"range(1|99)"`
}

// @Summary Create Order Item
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
// @Router /orders/{oid}/items/{id} [get]
func (h *Handler) getOrderItemById(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	orderId, err := strconv.Atoi(ctx.Param("oid"))
	if err != nil || orderId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderId")
	}

	orderItemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || orderItemId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderItemId")
	}

	orderItem, err := h.services.Order.GetItemById(clientId, clientType, orderId, orderItemId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, orderItem)
}

// @Summary Delete Order Item
// @Security UserAuth
// @Tags orders
// @Description delete order item
// @ModuleID deleteOrderItem
// @Accept  json
// @Produce  json
// @Param oid path string true "Order id"
// @Param id path string true "Order item id"
// @Success 200 {object} response
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /orders/{oid}/items/{id} [delete]
func (h *Handler) deleteOrderItem(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	orderId, err := strconv.Atoi(ctx.Param("oid"))
	if err != nil || orderId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderId")
	}

	orderItemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || orderItemId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderItemId")
	}

	err = h.services.Order.DeleteItem(clientId, clientType, orderId, orderItemId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

type orderItemUpdate struct {
	Count int `json:"count" valid:"range(1|99)"`
}

// @Summary Update Order Item
// @Security UserAuth
// @Tags orders
// @Description update order item
// @ModuleID updateOrderItem
// @Accept  json
// @Produce  json
// @Param oid path string true "Order id"
// @Param id path string true "Order item id"
// @Param input body orderItemUpdate true "order item update info"
// @Success 200 {object} response
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /orders/{oid}/items/{id} [put]
func (h *Handler) updateOrderItem(ctx echo.Context) error {
	var input orderItemUpdate
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	orderId, err := strconv.Atoi(ctx.Param("oid"))
	if err != nil || orderId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderId")
	}

	orderItemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || orderItemId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid orderItemId")
	}

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err = h.services.Order.UpdateItem(clientId, clientType, orderId, orderItemId, input.Count)

	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}

// @Summary Get All Active Orders For Restaurant
// @Security RestaurantAuth
// @Tags orders
// @Description get all active orders for restaurant
// @ModuleID getActiveRestaurantOrders
// @Accept  json
// @Produce  json
// @Param rid path string true "Restaurant id"
// @Success 200 {array} domain.Order
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid}/orders/ [get]
func (h *Handler) getActiveRestaurantOrders(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurantId, err := strconv.Atoi(ctx.Param("rid"))
	if err != nil || restaurantId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	orders, err := h.services.Order.GetActiveRestaurantOrders(clientId, clientType, restaurantId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, orders)
}

// @Summary Get All Orders
// @Security UserAuth
// @Security RestaurantAuth
// @Tags orders
// @Description get all orders for user
// @ModuleID getAllOrders
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Order
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/:id/orders [get]
func (h *Handler) usersGetAllOrders(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || userId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid userId")
	}

	orderStatus := ctx.QueryParam("status")
	activeOrdersFlag := false
	if orderStatus == "active" {
		activeOrdersFlag = true
	}

	orders, err := h.services.User.GetAllOrders(clientId, clientType, userId, activeOrdersFlag)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, orders)
}

// @Summary Get Active Order
// @Security CourierAuth
// @Tags couriers
// @Description get active order for courier
// @ModuleID getActiveCourierOrder
// @Accept  json
// @Produce  json
// @Param cid path string true "Courier id"
// @Success 200 {object} response
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /couriers/{cid}/orders [get]
func (h *Handler) getActiveCourierOrder(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	courierId, err := strconv.Atoi(ctx.Param("cid"))
	if err != nil || courierId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	order, err := h.services.Order.GetActiveCourierOrder(clientId, clientType, courierId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, order)
}
