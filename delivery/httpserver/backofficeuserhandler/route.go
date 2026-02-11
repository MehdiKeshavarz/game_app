package backofficeuserhandler

import (
	"game_app/delivery/httpserver/middleware"
	"game_app/entity"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {

	userGroup := e.Group("/backoffice/users")

	userGroup.GET("/get-all-users", h.getAllUsers, middleware.Auth(h.authSvc, h.authConfig),
		middleware.AccessCheck(h.authorizationSvc, entity.UserListPermission))

}
