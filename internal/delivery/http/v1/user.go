package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

const (
	authorizationHeader = "Authorization"
	idCtx               = "id"
	clientTypeCtx       = "client_type"
)

func (h *Handler) initUserRoutes(api *echo.Group) {
	users := api.Group("/users")
	{
		users.POST("/sign-up", h.usersSignUp)
		users.POST("/sign-in", h.usersSignIn)
	}
}

type userSignUpInput struct {
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

func (h *Handler) usersSignUp(ctx echo.Context) error {
	var input userSignUpInput

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]int{"id": 0})
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

func getToken(ctx echo.Context) (string, error) {
	header := ctx.Request().Header.Get(authorizationHeader)
	if header == "" {
		return "", errors.New("Empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("Invalid token")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("Empty token")
	}

	return headerParts[1], nil
}
