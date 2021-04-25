package v1

import (
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

func (h *Handler) initAdminRoutes(api *echo.Group) {
	admins := api.Group("/admins")
	{
		admins.POST("/sign-in", h.adminsSignIn)
	}
}

type adminSignInInput struct {
	Name     string `json:"name" valid:"length(4|32)"`
	Password string `json:"password" valid:"length(4|50)"`
}

// @Summary Admin SignIn
// @Tags admins
// @Description admin sign in
// @ModuleID adminSignIn
// @Accept  json
// @Produce  json
// @Param input body adminSignInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /admins/sign-in [post]
func (h *Handler) adminsSignIn(ctx echo.Context) error {
	var input adminSignInInput

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
