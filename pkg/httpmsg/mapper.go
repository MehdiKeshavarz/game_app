package httpmsg

import (
	"game_app/pkg/richerror"
	"net/http"
)

func Error(err error) (code int, message string) {
	switch err.(type) {
	case richerror.RichError:
		re := err.(richerror.RichError)
		msg := re.Message()
		code := MapKindToHTTPStatusCode(re.Kind())

		if code >= 500 {
			msg = "something went wrong "
		}

		return code, msg
	default:
		return http.StatusBadRequest, err.Error()
	}
}

func MapKindToHTTPStatusCode(k richerror.Kind) int {
	switch k {
	case richerror.KindInvalid:
		return http.StatusUnprocessableEntity
	case richerror.KindNotFound:
		return http.StatusNotFound
	case richerror.KindUnexpected:
		return http.StatusInternalServerError
	case richerror.KindForbidden:
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}
