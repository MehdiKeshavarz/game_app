package main

import (
	"fmt"
	"game_app/adapter/redis"
	"game_app/config"
	"game_app/delivery/httpserver"
	"game_app/repository/migrator"
	"game_app/repository/mysql"
	"game_app/repository/mysql/accesscontrol"
	"game_app/repository/mysql/user"
	"game_app/repository/redis/matching"
	"game_app/service/authorizationservice"
	"game_app/service/authservice"
	"game_app/service/backofficeuserservice"
	"game_app/service/matchingservice"
	"game_app/service/userservice"
	"game_app/validator/matchingvalidator"
	"game_app/validator/uservalidator"
)

const (
	JwtSignKey = "jwt_secret"
)

func main() {
	cfg := config.Load()
	fmt.Printf("cfg: %+v\n", cfg)

	mgr := migrator.New(cfg.Mysql, "mysql")
	mgr.Up()

	authSvc, userSvc, userValidator, authorizationSvc, backofficeUserSvc, matchingSvc, matchingValidator := setupServices(cfg)

	server := httpserver.New(cfg,
		authSvc,
		userSvc,
		authorizationSvc,
		userValidator,
		backofficeUserSvc,
		matchingSvc,
		matchingValidator)

	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator, authorizationservice.Service, backofficeuserservice.Service, matchingservice.Service, matchingvalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.Mysql)

	userMysql := user.New(mysqlRepo)

	userValidator := uservalidator.New(userMysql)
	userSvc := userservice.New(userMysql, authSvc)

	backofficeUserSvc := backofficeuserservice.New()

	accesscontrolMysql := accesscontrol.New(mysqlRepo)

	authorizationSvc := authorizationservice.New(accesscontrolMysql)

	matchingValidator := matchingvalidator.New()
	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := matching.New(redisAdapter)
	matchingSvc := matchingservice.New(matchingRepo, cfg.MatchingService)

	return authSvc, userSvc, userValidator, authorizationSvc, backofficeUserSvc, matchingSvc, matchingValidator
}
