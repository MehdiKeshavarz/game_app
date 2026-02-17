package matchinghandler

import (
	"game_app/service/authservice"
	"game_app/service/matchingservice"
	"game_app/validator/matchingvalidator"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	matchingSvc       matchingservice.Service
	matchingValidator matchingvalidator.Validator
}

func New(
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator,
	authConfig authservice.Config,
	authSvc authservice.Service) Handler {
	return Handler{
		authSvc:           authSvc,
		authConfig:        authConfig,
		matchingSvc:       matchingSvc,
		matchingValidator: matchingValidator,
	}
}
