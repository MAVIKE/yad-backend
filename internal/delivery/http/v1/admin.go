package v1

import (
	"github.com/MAVIKE/yad-backend/internal/domain"
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
		Name     string `json:"name" valid:"length(3|32)"`
		Password string `json:"password" valid:"length(6|32)"`
	}
	var input signInInput

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	admin := &domain.Admin{
		Name:     input.Name,
		Password: input.Password,
	}

	return ctx.JSON(http.StatusOK, admin)
}
