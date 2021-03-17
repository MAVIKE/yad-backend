package v1

import (
	"net/http"
	"strconv"

	_ "github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initRestaurantRoutes(api *echo.Group) {
	restaurants := api.Group("/restaurants")
	{
		restaurants.POST("/sign-in", h.restaurantsSignIn)
		restaurants.Use(h.identity)
		restaurants.GET("", h.getRestaurants)
		restaurants.GET("/:rid", h.getRestaurantById)
	}
}

type restaurantsSignInInput struct {
	Phone    string `json:"phone" valid:"numeric"`
	Password string `json:"password" valid:"length(4|32)"`
}

// @Summary Courier SignIn
// @Tags couriers-auth
// @Description courier sign in
// @ModuleID courierSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /couriers/sign-in [post]
func (h *Handler) restaurantsSignIn(ctx echo.Context) error {
	var input restaurantsSignInInput

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	token, err := h.services.Restaurant.SignIn(input.Phone, input.Password)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, tokenResponse{
		AccessToken: token.AccessToken,
	})
}

// @Summary Get All Restaurant
// @Tags restaurant
// @Description restaurant get all
// @ModuleID restaurantGetAll
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Restaurant
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants [get]
func (h *Handler) getRestaurants(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurants, err := h.services.Restaurant.GetAll(clientId, clientType)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, restaurants)
}

func (h *Handler) getRestaurantById(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurantId, err := strconv.Atoi(ctx.Param("rid"))
	if err != nil || restaurantId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	restaurant, err := h.services.Restaurant.GetById(clientId, clientType, restaurantId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, restaurant)
}
