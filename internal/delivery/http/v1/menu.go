package v1

import (
	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/asaskevich/govalidator"
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
		menu.PUT("/:rid/menu/:id", h.updateMenuItem)
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

type menuItemUpdate struct {
	Title       string `json:"title" db:"title"`
	Image       string `json:"image" db:"image"`
	Description string `json:"description" db:"description"`
	Price       int    `json:"price" db:"price"`
	CategoryId  int
}

// @Summary Update Menu Item
// @Security RestaurantAuth
// @Tags restaurants
// @Description update menu item
// @ModuleID updateMenuItem
// @Accept  json
// @Produce  json
// @Param rid path string true "Restaurant id"
// @Param id path string true "MenuItem id"
// @Param input body menuItemUpdate true "menu item update info"
// @Success 200 {object} response
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid}/menu/{id} [put]
func (h *Handler) updateMenuItem(ctx echo.Context) error {
	var input menuItemUpdate
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

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	update := &domain.MenuItem{
		Title:       input.Title,
		Image:       input.Image,
		Description: input.Description,
		Price:       input.Price,
	}

	err = h.services.MenuItem.UpdateMenuItem(clientId, clientType, restaurantId, menuItemId, input.CategoryId, update)

	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
