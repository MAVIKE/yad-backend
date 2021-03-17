package v1

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) initCourierRoutes(api *echo.Group) {
	couriers := api.Group("/couriers")
	{
		couriers.POST("/sign-in", h.couriersSignIn)
		couriers.POST("/sign-up", h.couriersSignUp)
	}
}

type courierSignUpInput struct {
	Name          string `json:"name"`
	Phone         string `json:"phone"`
	Password      string `json:"password"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	WorkingStatus int    `json:"working_status"`
}

func (h *Handler) couriersSignUp(ctx echo.Context) error {
	var input courierSignUpInput

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]int{"id": 0})
}

type courierSignInInput struct {
	Phone    string `json:"phone" valid:"numeric"`
	Password string `json:"password" valid:"length(4|32)"`
}

// @Summary Courier SignIn
// @Tags couriers-auth
// @Description courier sign in
// @ModuleID courierSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInputPhone true "sign up info"
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
