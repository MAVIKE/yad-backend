package v1

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/MAVIKE/yad-backend/internal/service"
	"github.com/MAVIKE/yad-backend/pkg/auth"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	services     *service.Service
	tokenManager *auth.Manager
}

func NewHandler(services *service.Service, tokenManager *auth.Manager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *echo.Group) {
	v1 := api.Group("/v1")
	{
		v1.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})

		h.initAdminRoutes(v1)
		h.initUserRoutes(v1)
		h.initCourierRoutes(v1)
		h.initRestaurantRoutes(v1)
		h.initCategoryRoutes(v1)
	}
}

func (h *Handler) identity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token, err := getToken(ctx)
		if err != nil {
			return newResponse(ctx, http.StatusUnauthorized, err.Error())
		}

		userId, clientType, err := h.tokenManager.Parse(token)
		if err != nil {
			return newResponse(ctx, http.StatusUnauthorized, err.Error())
		}

		ctx.Request().Header.Set(idCtx, strconv.Itoa(userId))
		ctx.Request().Header.Set(clientTypeCtx, clientType)
		return next(ctx)
	}
}

func (h *Handler) getClientParams(ctx echo.Context) (int, string, error) {
	id := ctx.Request().Header.Get(idCtx)
	if id == "" {
		return 0, "", errors.New("user id not found")
	}

	intId, err := strconv.Atoi(id)
	if err != nil {
		return 0, "", errors.New("user id is of invalid type")
	}

	clientType := ctx.Request().Header.Get(clientTypeCtx)
	if clientType == "" {
		return 0, "", errors.New("client type not found")
	}

	return intId, clientType, nil
}
