package v1

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) initAdminRoutes(api *echo.Group) {
	admins := api.Group("/admins")
	{
		admins.POST("/sign-in", h.adminsSignIn)
	}
}

func (h *Handler) adminsSignIn(ctx echo.Context) error {
	type signInInput struct {
		Name     string `json:"name" valid:"length(4|32)"`
		Password string `json:"password" valid:"length(4|32)"`
	}
	var input signInInput

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	token, err := h.services.Admin.SignIn(input.Name, input.Password)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, tokenResponse{
		AccessToken: token.AccessToken,
	})
}
