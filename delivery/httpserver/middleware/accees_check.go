package middleware

import (
	"game_app/entity"
	"game_app/pkg/claim"
	"game_app/pkg/errmsg"
	"game_app/service/authorizationservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AccessCheck(service authorizationservice.Service, permissions ...entity.PermissionTitle) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := claim.GetClaims(c)
			isAllowed, err := service.CheckAccess(claims.UserID, claims.Role, permissions...)

			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorSomethingWentWrong})
			}

			if !isAllowed {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errmsg.ErrorAccessDenied})
			}

			return next(c)
		}
	}
}
