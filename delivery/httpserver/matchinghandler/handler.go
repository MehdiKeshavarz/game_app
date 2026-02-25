package matchinghandler

import (
	"game_app/service/authservice"
	"game_app/service/matchingservice"
	"game_app/service/presenceservice"
	"game_app/validator/matchingvalidator"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	matchingSvc       matchingservice.Service
	matchingValidator matchingvalidator.Validator
	presenceSvc       presenceservice.Service
}

func New(
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator,
	authConfig authservice.Config,
	authSvc authservice.Service,
	presenceSvc presenceservice.Service) Handler {
	return Handler{
		authSvc:           authSvc,
		authConfig:        authConfig,
		matchingSvc:       matchingSvc,
		matchingValidator: matchingValidator,
		presenceSvc:       presenceSvc,
	}
}
