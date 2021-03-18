package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/MAVIKE/yad-backend/internal/domain"

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
	Name      string  `json:"name"`
	Phone     string  `json:"phone" valid:"required,numeric,length(11|11)"`
	Password  string  `json:"password" valid:"required,length(8|50)"`
	Email     string  `json:"email" valid:"email"`
	Latitude  float64 `json:"latitude" valid:"required,latitude"`
	Longitude float64 `json:"longitude" valid:"required,longitude"`
}

// @Summary User SignUp
// @Tags users
// @Description user sign up
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body userSignUpInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/sign-up [post]
func (h *Handler) usersSignUp(ctx echo.Context) error {
	var input userSignUpInput

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	user := &domain.User{
		Name:     input.Name,
		Phone:    input.Phone,
		Password: input.Password,
		Email:    input.Email,
		Address: &domain.Location{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
	}

	id, err := h.services.User.SignUp(user)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type userSignInInput struct {
	Phone    string `json:"phone" valid:"numeric"`
	Password string `json:"password" valid:"length(4|32)"`
}

// @Summary User SignIn
// @Tags users
// @Description user sign in
// @ModuleID userSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign in info"
// @Success 200 {object} tokenResponse
// @Failure 400,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/sign-in [post]
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
