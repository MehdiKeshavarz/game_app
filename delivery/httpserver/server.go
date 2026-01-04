package httpserver

import (
	"fmt"
	"game_app/config"
	"game_app/delivery/httpserver/userhandler"
	"game_app/service/authservice"
	"game_app/service/userservice"
	"game_app/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config      config.Config
	userHandler userhandler.Handler
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator) Server {
	return Server{
		config:      config,
		userHandler: userhandler.New(authSvc, userSvc, userValidator, config.AuthConfig),
	}
}

func (s Server) Serve() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health-check", s.healthCheck)
	s.userHandler.SetUserRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HttpServer.Port)))
}
