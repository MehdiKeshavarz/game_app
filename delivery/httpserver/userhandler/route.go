package userhandler

import (
	"game_app/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	groups := e.Group("/users")

	groups.POST("/register", h.userRegister)
	groups.POST("/login", h.userLogin)
	groups.GET("/profile", h.userProfile, middleware.Auth(h.authSvc, h.authConfig))

}
