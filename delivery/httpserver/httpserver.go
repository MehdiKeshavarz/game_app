package httpserver

import (
	"fmt"
	"game_app/config"
	"game_app/service/authservice"
	"game_app/service/userservice"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config  config.Config
	authSvc authservice.Service
	userSvc userservice.Service
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service) Server {
	return Server{
		config:  config,
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (s Server) Serve() {
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/health-check", s.healthCheck)

	groups := e.Group("/users")

	groups.POST("/register", s.userRegister)
	groups.POST("/login", s.userLogin)
	groups.GET("/profile", s.userProfile)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HttpServer.Port)))
}
