package backofficeuserhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) getAllUsers(c echo.Context) error {

	res := h.backofficeUserSvc.ListAllUser()

	return c.JSON(http.StatusOK, res)
}
