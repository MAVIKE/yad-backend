package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

const (
	authorizationHeader = "Authorization"
	IdCtx               = "id"
	ClientTypeCtx       = "client_type"
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

func (h *Handler) identity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := getToken(ctx)
		if err != nil {
			return newResponse(ctx, http.StatusUnauthorized, err.Error())
		}

		userId, clientType, err := h.services.User.ParseToken(token)
		if err != nil {
			return newResponse(ctx, http.StatusUnauthorized, err.Error())
		}

		ctx.Request().Header.Set(IdCtx, strconv.Itoa(userId))
		ctx.Request().Header.Set(ClientTypeCtx, clientType)
		return next(ctx)
	}
}

func getClientParams(ctx echo.Context) (int, string, error) {
	id := ctx.Request().Header.Get(IdCtx)
	if id == "" {
		return 0, "", errors.New("user id not found")
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		return 0, "", errors.New("user id is of invalid type")
	}

	clientType := ctx.Request().Header.Get(ClientTypeCtx)
	if clientType == "" {
		return 0, "", errors.New("client type not found")
	}

	return intId, clientType, nil
}
