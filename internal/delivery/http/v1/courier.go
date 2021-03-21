package v1

import (
	"net/http"
	"strconv"

	"github.com/MAVIKE/yad-backend/internal/domain"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initCourierRoutes(api *echo.Group) {
	couriers := api.Group("/couriers")
	{
		couriers.POST("/sign-in", h.couriersSignIn)
		couriers.Use(h.identity)
		couriers.POST("/sign-up", h.couriersSignUp)
		couriers.GET("/:id", h.getCourierById)
	}
}

type courierSignUpInput struct {
	Name          string  `json:"name"`
	Phone         string  `json:"phone" valid:"required,numeric,length(11|11)"`
	Password      string  `json:"password" valid:"required,length(8|50)"`
	Email         string  `json:"email" valid:"email"`
	Latitude      float64 `json:"latitude" valid:"required,latitude"`
	Longitude     float64 `json:"longitude" valid:"required,longitude"`
	WorkingStatus int     `json:"working_status"`
}

// @Summary Courier SignUp
// @Tags couriers
// @Description courier sign up
// @ModuleID courierSignUp
// @Accept  json
// @Produce  json
// @Param input body courierSignUpInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /couriers/sign-up [post]
func (h *Handler) couriersSignUp(ctx echo.Context) error {
	var input courierSignUpInput
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

	courier := &domain.Courier{
		Name:     input.Name,
		Phone:    input.Phone,
		Password: input.Password,
		Email:    input.Email,
		Address: &domain.Location{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
		WorkingStatus: input.WorkingStatus,
	}

	id, err := h.services.Courier.SignUp(courier, clientType)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type courierSignInInput struct {
	Phone    string `json:"phone" valid:"numeric,length(11|11)"`
	Password string `json:"password" valid:"length(8|50)"`
}

// @Summary Courier SignIn
// @Tags couriers
// @Description courier sign in
// @ModuleID courierSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /couriers/sign-in [post]
func (h *Handler) couriersSignIn(ctx echo.Context) error {
	var input courierSignInInput

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	token, err := h.services.Courier.SignIn(input.Phone, input.Password)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, tokenResponse{
		AccessToken: token.AccessToken,
	})
}

// @Summary Get Courier By Id
// @Security UserAuth
// @Security RestaurantAuth
// @Tags categories
// @Description get courier by id
// @ModuleID getCourierById
// @Accept  json
// @Produce  json
// @Param id path string true "Courier id"
// @Success 200 {object} domain.Category
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /couriers/{id} [get]
func (h *Handler) getCourierById(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	courierId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || courierId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid courier")
	}

	courier, err := h.services.Courier.GetById(clientId, clientType, courierId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, courier)
}
