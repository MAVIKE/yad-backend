package v1

import (
	"net/http"
	"strconv"

	_ "github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initCategoryRoutes(api *echo.Group) {
	restaurants := api.Group("/restaurants/:rid/categories")
	{
		restaurants.Use(h.identity)
		restaurants.GET("", h.getCategories)
	}
}

func (h *Handler) getCategories(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurantId, err := strconv.Atoi(ctx.Param("rid"))
	if err != nil || restaurantId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	categories, err := h.services.Category.GetAll(clientId, clientType, restaurantId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, categories)
}
