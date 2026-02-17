package matchinghandler

import (
	"game_app/param"
	"game_app/pkg/claim"
	"game_app/pkg/httpmsg"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) addToWaitingList(c echo.Context) error {
	var req param.AddToWaitingListRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err, filedErrors := h.matchingValidator.ValidateAddToWaitingList(req)

	if err != nil {
		code, msg := httpmsg.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"error":   filedErrors,
		})

	}
	claims := claim.GetClaims(c)

	req.UserID = claims.UserID

	res, err := h.matchingSvc.AddToWaitingList(req)

	if err != nil {
		code, msg := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusCreated, res)
}
