package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/MAVIKE/yad-backend/pkg/download"
	"github.com/MAVIKE/yad-backend/pkg/random"
	"github.com/labstack/echo/v4/middleware"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initRestaurantRoutes(api *echo.Group) {
	restaurants := api.Group("/restaurants")
	{
		restaurants.POST("/sign-in", h.restaurantsSignIn)
		restaurants.Use(h.identity)
		restaurants.POST("/sign-up", h.restaurantsSignUp)
		restaurants.GET("/", h.getRestaurants)
		restaurants.GET("/:rid", h.getRestaurantById)
		restaurants.GET("/image", h.getRestaurantImage)
		restaurants.PUT("/:rid/image", h.updateRestaurantImage, middleware.BodyLimit("10M"))
		restaurants.PUT("/:rid", h.updateRestaurant)
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
// @Param input body restaurantsSignInInput true "sign up info"
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

type restaurantSignUpInput struct {
	Name          string        `json:"name"`
	Phone         string        `json:"phone" valid:"required,numeric,length(11|11)"`
	Password      string        `json:"password" valid:"required,length(8|50)"`
	Address       locationInput `json:"address" valid:"required"`
	WorkingStatus int           `json:"working_status"`
	Image         string        `json:"image" valid:"required,length(1|200)"`
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
			Latitude:  input.Address.Latitude,
			Longitude: input.Address.Longitude,
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

// @Summary Get Restaurant Image
// @Security UserAuth
// @Security RestaurantAuth
// @Tags restaurants
// @Description get restaurant image
// @ModuleID getRestaurantImage
// @Accept json
// @Produce json
// @Success 200 {object} string "binary file"
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/image [get]
func (h *Handler) getRestaurantImage(ctx echo.Context) error {
	_, _, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	var image imageInput

	if err := ctx.Bind(&image); err != nil {
		return newResponse(ctx, http.StatusBadRequest, "Bad request")
	}

	return ctx.File(image.Path)
}

// @Summary Update Restaurant Image
// @Security RestaurantAuth
// @Tags restaurants
// @Description update restaurant image
// @ModuleID updateRestaurantImage
// @Accept json
// @Produce json
// @Success 200 {object} domain.Restaurant
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid}/image [put]
func (h *Handler) updateRestaurantImage(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurantId, err := strconv.Atoi(ctx.Param("rid"))
	if err != nil || restaurantId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid restaurantId")
	}

	// TODO: move checking to service
	if clientType != "restaurant" && restaurantId != clientId {
		return newResponse(ctx, http.StatusBadRequest, "Forbidden")
	}

	_, err = h.services.Restaurant.GetById(clientId, clientType, restaurantId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if !isImage(file.Filename) {
		return newResponse(ctx, http.StatusBadRequest, errors.New("not image").Error())
	}

	fileName := imageDir + "restaurant_" + ctx.Param("rid") + "_" + random.GetString(5) + "_" + file.Filename
	if err := download.Download(file, fileName); err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	restaurant, err := h.services.Restaurant.UpdateImage(clientId, clientType, restaurantId, fileName)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, restaurant)
}

type restaurantUpdateInput struct {
	Name          string        `json:"name"`
	Password      string        `json:"password" valid:"length(8|50)"`
	Address       locationInput `json:"address"`
	WorkingStatus int           `json:"working_status"`
}

// @Summary Update Restaurant
// @Security RestaurantAuth
// @Tags restaurants
// @Description update restaurant
// @ModuleID updateRestaurant
// @Accept  json
// @Produce  json
// @Param rid path string true "Restaurant id"
// @Param input body restaurantUpdateInput true "restaurant update info"
// @Success 200 {object} response
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /restaurants/{rid} [put]
func (h *Handler) updateRestaurant(ctx echo.Context) error {
	var input restaurantUpdateInput

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

	update := &domain.Restaurant{
		Name:     input.Name,
		Password: input.Password,
		Address: &domain.Location{
			Latitude:  input.Address.Latitude,
			Longitude: input.Address.Longitude,
		},
		WorkingStatus: input.WorkingStatus,
	}

	err = h.services.Restaurant.Update(clientId, clientType, restaurantId, update)

	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
