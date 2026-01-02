package userhandler

import (
	"game_app/dto"
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
}

func New(authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator) Handler {
	return Handler{
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: userValidator,
	}
}

func (h Handler) userRegister(c echo.Context) error {
	var req dto.RegisterRequest
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
	var req dto.LoginRequest
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

func (h Handler) userProfile(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	
	claims, err := h.authSvc.ParseToken(auth)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	res, rErr := h.userSvc.Profile(dto.GetProfileRequest{UserID: claims.UserID})

	if rErr != nil {
		code, msg := httpmsg.Error(rErr)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, res)
}
