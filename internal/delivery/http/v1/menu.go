package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *Handler) initMenuRoutes(api *echo.Group) {
	menu := api.Group("/restaurants")
	{
		menu.Use(h.identity)
		menu.GET("/:rid/menu/", h.getRestaurantMenu)
		menu.GET("/:rid/menu/:id", h.getMenuItemById)
	}
}

// @Summary Get Restaurant Menu
// @Security UserAuth
// @Security RestaurantAuth
// @Tags restaurants
// @Description get restaurant menu
// @ModuleID getRestaurantMenu
// @Accept  json
// @Produce  json
// @Param rid path string true "Restaurant id"
// @Success 200 {array} domain.MenuItem
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid}/menu/ [get]
func (h *Handler) getRestaurantMenu(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurantId, err := strconv.Atoi(ctx.Param("rid"))
	if err != nil || restaurantId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	menu, err := h.services.Restaurant.GetMenu(clientId, clientType, restaurantId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, menu)
}

// @Summary Get Menu Item By Id
// @Tags restaurants
// @Description get menu item by id
// @ModuleID getMenuItemById
// @Accept  json
// @Produce  json
// @Param rid path string true "Restaurant id"
// @Param id path string true "MenuItem id"
// @Success 200 {object} domain.MenuItem
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid}/menu/{id} [get]
func (h *Handler) getMenuItemById(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurantId, err := strconv.Atoi(ctx.Param("rid"))
	if err != nil || restaurantId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	menuItemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || menuItemId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid menuItemId")
	}

	menuItem, err := h.services.MenuItem.GetById(clientId, clientType, menuItemId, restaurantId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, menuItem)
}
