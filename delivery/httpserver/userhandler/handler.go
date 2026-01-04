package userhandler

import (
	cfg "game_app/config"
	"game_app/param"
	"game_app/pkg/httpmsg"
	"game_app/service/authservice"
	"game_app/service/userservice"
	"game_app/validator/uservalidator"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	authConfig    authservice.Config
}

func New(authSvc authservice.Service,
	userSvc userservice.Service,
	userValidator uservalidator.Validator,
	authConfig authservice.Config) Handler {
	return Handler{
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: userValidator,
		authConfig:    authConfig,
	}
}

func (h Handler) userRegister(c echo.Context) error {
	var req param.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err, filedErrors := h.userValidator.ValidateRegisterRequest(req)

	if err != nil {
		code, msg := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"error":   filedErrors,
		})

	}

	res, err := h.userSvc.Register(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (h Handler) userLogin(c echo.Context) error {
	var req param.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err, filedErrors := h.userValidator.ValidateLoginRequest(req)

	if err != nil {
		code, msg := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"error":   filedErrors,
		})

	}
	res, err := h.userSvc.Login(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func getClaims(c echo.Context) *authservice.Claims {
	// let it crash
	return c.Get(cfg.AuthMiddlewareContextKey).(*authservice.Claims)
}

func (h Handler) userProfile(c echo.Context) error {
	claims := getClaims(c)

	res, rErr := h.userSvc.Profile(param.GetProfileRequest{UserID: claims.UserID})

	if rErr != nil {
		code, msg := httpmsg.Error(rErr)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, res)
}
