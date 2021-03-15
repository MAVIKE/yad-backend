package v1

import (
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handler) initUserRoutes(api *echo.Group) {
	admins := api.Group("/users")
	{
		admins.POST("/sign-in", h.usersSignIn)
	}
}

type userSignInInput struct {
	Phone    string `json:"phone" valid:"numeric"`
	Password string `json:"password" valid:"length(4|32)"`
}

func (h *Handler) usersSignIn(ctx echo.Context) error {
	var input userSignInInput

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	token, err := h.services.User.SignIn(input.Phone, input.Password)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, tokenResponse{
		AccessToken: token.AccessToken,
	})
}
