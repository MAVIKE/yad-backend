package v1

import (
	"net/http"
	"strconv"

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
		users.Use(h.identity)
		users.PUT("/:uid", h.updateUser)
		users.GET("/:id", h.getUserById)
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
	Phone    string `json:"phone" valid:"numeric,length(11|11)"`
	Password string `json:"password" valid:"length(8|50)"`
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

type userUpdate struct {
	Name      string  `json:"name"`
	Password  string  `json:"password" valid:"required,length(8|50)"`
	Email     string  `json:"email" valid:"email"`
	Latitude  float64 `json:"latitude" valid:"required,latitude"`
	Longitude float64 `json:"longitude" valid:"required,longitude"`
}

// @Summary Update User
// @Security UserAuth
// @Tags users
// @Description update user
// @ModuleID updateUser
// @Accept  json
// @Produce  json
// @Param uid path string true "User id"
// @Param input body userUpdate true "user update info"
// @Success 200 {object} response
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/{uid} [put]
func (h *Handler) updateUser(ctx echo.Context) error {
	var input userUpdate
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	userId, err := strconv.Atoi(ctx.Param("uid"))
	if err != nil || userId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid userId")
	}

	if err := ctx.Bind(&input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	if _, err := govalidator.ValidateStruct(input); err != nil {
		return newResponse(ctx, http.StatusBadRequest, err.Error())
	}

	update := &domain.User{
		Name:     input.Name,
		Password: input.Password,
		Email:    input.Email,
		Address: &domain.Location{
			Latitude:  input.Latitude,
			Longitude: input.Longitude,
		},
	}

	err = h.services.User.Update(clientId, clientType, userId, update)

	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}


// @Summary Get User By Id
// @Security UserAuth
// @Security RestaurantAuth
// @Security CourierAuth
// @Tags users
// @Description get user by id
// @ModuleID getUserById
// @Accept  json
// @Produce  json
// @Param id path string true "user id"
// @Success 200 {object} domain.User
// @Failure 400,403,404 {object} response
// @Failure 500 {object} response
// @Failure default {object} response
// @Router /users/{id} [get]
func (h *Handler) getUserById(ctx echo.Context) error {
	clientId, clientType, err := h.getClientParams(ctx)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || userId == 0 {
		return newResponse(ctx, http.StatusBadRequest, "Invalid courier")
	}

	user, err := h.services.User.GetById(clientId, clientType, userId)
	if err != nil {
		return newResponse(ctx, http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, user)
}
