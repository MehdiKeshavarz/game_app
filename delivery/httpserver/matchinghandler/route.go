package matchinghandler

import (
	"game_app/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	groups := e.Group("/matching")

	groups.POST("/add-to-waiting-list", h.addToWaitingList,
		middleware.Auth(h.authSvc, h.authConfig),
		middleware.UpsertPresence(h.presenceSvc))

}
