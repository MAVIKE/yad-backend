package v1

import (
	"net/http"
	"strconv"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initRestaurantRoutes(api *echo.Group) {
	restaurants := api.Group("/restaurants")
	{
		restaurants.POST("/sign-in", h.restaurantsSignIn)
		restaurants.Use(h.identity)
		restaurants.GET("/", h.getRestaurants)
		restaurants.GET("/:rid", h.getRestaurantById)
		restaurants.GET("/:rid/menu/", h.getRestaurantMenu)
		restaurants.POST("/sign-up", h.restaurantsSignUp)
		restaurants.GET("/:rid/menu/:id", h.getMenuItemById)
	}
}

type restaurantsSignInInput struct {
	Phone    string `json:"phone" valid:"numeric,length(11|11)"`
	Password string `json:"password" valid:"length(8|50)"`
}

// @Summary Restaurant SignIn
// @Tags restaurants
// @Description restaurant sign in
// @ModuleID restaurantSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/sign-in [post]
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

// @Summary Get All Restaurants
// @Security UserAuth
// @Security RestaurantAuth
// @Tags restaurants
// @Description get all restaurants for user
// @ModuleID getAllRestaurants
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Restaurant
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/ [get]
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

// @Summary Get Restaurant By Id
// @Security UserAuth
// @Security RestaurantAuth
// @Tags restaurants
// @Description get restaurant by id for user
// @ModuleID getRestaurantById
// @Accept  json
// @Produce  json
// @Param rid path string true "Restaurant id"
// @Success 200 {object} domain.Restaurant
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid} [get]
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

type restaurantSignUpInput struct {
	Name          string  `json:"name"`
	Phone         string  `json:"phone" valid:"required,numeric,length(11|11)"`
	Password      string  `json:"password" valid:"required,length(8|50)"`
	Latitude      float64 `json:"latitude" valid:"required,latitude"`
	Longitude     float64 `json:"longitude" valid:"required,longitude"`
	WorkingStatus int     `json:"working_status"`
	Image         string  `json:"image" valid:"required,length(1|200)"`
}

// @Summary Restaurant SignUp
// @Tags restaurants
// @Description restaurant sign up
// @ModuleID restaurantSignUp
// @Accept  json
// @Produce  json
// @Param input body restaurantSignUpInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/sign-up [post]
func (h *Handler) restaurantsSignUp(ctx echo.Context) error {
	var input restaurantSignUpInput
	_, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	restaurant := &domain.Restaurant{
		Name:     input.Name,
		Phone:    input.Phone,
		Password: input.Password,
		Address: &domain.Location{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
		WorkingStatus: input.WorkingStatus,
		Image:         input.Image,
	}

	id, err := h.services.Restaurant.SignUp(restaurant, clientType)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
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
