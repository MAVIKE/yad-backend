package v1

import (
	"net/http"
	"strconv"

	"github.com/MAVIKE/yad-backend/internal/domain"
	_ "github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initCategoryRoutes(api *echo.Group) {
	restaurants := api.Group("/restaurants/:rid/categories")
	{
		restaurants.Use(h.identity)
		restaurants.POST("/", h.createCategory)
		restaurants.GET("/", h.getCategories)
	}
}

type categoryInput struct {
	Title string `json:"title" valid:"length(1|50)"`
}

// @Summary Create Category
// @Security RestaurantAuth
// @Tags categories
// @Description create category
// @ModuleID createCategory
// @Accept  json
// @Produce  json
// @Param rid path string true "Restaurant id"
// @Success 200 {object} idResponse
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid}/categories/ [post]
func (h *Handler) createCategory(ctx echo.Context) error {
	var input categoryInput
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurantId, err := strconv.Atoi(ctx.Param("rid"))
	if err != nil || restaurantId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	category := &domain.Category{
		RestaurantId: restaurantId,
		Title:        input.Title,
	}

	categoryId, err := h.services.Category.Create(clientId, clientType, category)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, idResponse{
		Id: categoryId,
	})
}

// @Summary Get All Categories
// @Security UserAuth
// @Security RestaurantAuth
// @Tags categories
// @Description get all categories for user
// @ModuleID getAllCategories
// @Accept  json
// @Produce  json
// @Param rid path string true "Restaurant id"
// @Success 200 {array} domain.Category
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid}/categories/ [get]
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
