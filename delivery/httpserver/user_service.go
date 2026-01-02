package httpserver

import (
	"game_app/dto"
	"game_app/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) userRegister(c echo.Context) error {
	var req dto.RegisterRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err, filedErrors := s.userValidator.ValidateRegisterRequest(req)

	if err != nil {
		code, msg := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"error":   filedErrors,
		})

	}

	res, err := s.userSvc.Register(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, res)
}

func (s Server) userLogin(c echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	res, err := s.userSvc.Login(req)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, res)
}

func (s Server) userProfile(c echo.Context) error {
	auth := c.Request().Header.Get("Authorization")
	claims, err := s.authSvc.ParseToken(auth)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	res, rErr := s.userSvc.Profile(dto.GetProfileRequest{UserID: claims.UserID})

	if rErr != nil {
		code, msg := httpmsg.Error(rErr)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, res)
}
