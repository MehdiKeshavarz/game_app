package middleware

import (
	"game_app/param"
	"game_app/pkg/claim"
	"game_app/pkg/errmsg"
	"game_app/pkg/timestamp"
	"game_app/service/presenceservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func UpsertPresence(presenceSvc presenceservice.Service) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := claim.GetClaims(c)
			req := param.UpsertPresenceRequest{
				UserID:    claims.UserID,
				Timestamp: timestamp.Now(),
			}
			_, err := presenceSvc.UpsertPresence(c.Request().Context(), req)

			if err != nil {
				// TODO - log expected error
				// we can just log the error and go to next step(middleware , handler)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorSomethingWentWrong})
			}

			return next(c)
		}
	}
}
