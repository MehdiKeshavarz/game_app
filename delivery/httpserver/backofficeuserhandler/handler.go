package backofficeuserhandler

import (
	"game_app/service/authorizationservice"
	"game_app/service/authservice"
	"game_app/service/backofficeuserservice"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	authorizationSvc  authorizationservice.Service
	backofficeUserSvc backofficeuserservice.Service
}

func New(svc backofficeuserservice.Service,
	authSvc authservice.Service,
	authConfig authservice.Config,
	authorizationSvc authorizationservice.Service) Handler {
	return Handler{
		backofficeUserSvc: svc,
		authConfig:        authConfig,
		authSvc:           authSvc,
		authorizationSvc:  authorizationSvc,
	}
}
