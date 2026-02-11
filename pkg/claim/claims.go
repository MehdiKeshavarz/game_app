package claim

import (
	cfg "game_app/config"
	"game_app/service/authservice"

	"github.com/labstack/echo/v4"
)

func GetClaims(c echo.Context) *authservice.Claims {
	// let it crash
	return c.Get(cfg.AuthMiddlewareContextKey).(*authservice.Claims)
}
