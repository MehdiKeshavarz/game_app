package middleware

import (
	"game_app/pkg/constant"
	"game_app/service/authservice"

	mw "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func Auth(service authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey: constant.AuthMiddlewareContextKey,
		SigningKey: config.SignKey,
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			claims, err := service.ParseToken(auth)
			if err != nil {
				return nil, err
			}

			return claims, nil
		},
	})
}
