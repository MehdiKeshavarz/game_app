package httpserver

import (
	"fmt"
	"game_app/config"
	"game_app/delivery/httpserver/backofficeuserhandler"
	"game_app/delivery/httpserver/userhandler"
	"game_app/service/authorizationservice"
	"game_app/service/authservice"
	"game_app/service/backofficeuserservice"
	"game_app/service/userservice"
	"game_app/validator/uservalidator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config                config.Config
	userHandler           userhandler.Handler
	backofficeUserHandler backofficeuserhandler.Handler
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, authorizationSvc authorizationservice.Service, userValidator uservalidator.Validator, backofficeService backofficeuserservice.Service) Server {
	return Server{
		config:                config,
		userHandler:           userhandler.New(authSvc, userSvc, userValidator, config.Auth),
		backofficeUserHandler: backofficeuserhandler.New(backofficeService, authSvc, config.Auth, authorizationSvc),
	}
}

func (s Server) Serve() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health-check", s.healthCheck)
	s.userHandler.SetRoutes(e)
	s.backofficeUserHandler.SetRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HttpServer.Port)))
}
