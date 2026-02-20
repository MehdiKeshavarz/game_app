package main

import (
	"fmt"
	"game_app/adapter/redis"
	"game_app/config"
	"game_app/delivery/httpserver"
	"game_app/repository/mysql"
	"game_app/repository/mysql/accesscontrol"
	"game_app/repository/mysql/user"
	"game_app/repository/redis/matching"
	"game_app/scheduler"
	"game_app/service/authorizationservice"
	"game_app/service/authservice"
	"game_app/service/backofficeuserservice"
	"game_app/service/matchingservice"
	"game_app/service/userservice"
	"game_app/validator/matchingvalidator"
	"game_app/validator/uservalidator"
	"os"
	"os/signal"
	"time"
)

const (
	JwtSignKey = "jwt_secret"
)

func main() {
	cfg := config.Load()

	//mgr := migrator.New(cfg.Mysql, "mysql")
	//mgr.Up()

	authSvc, userSvc, userValidator, authorizationSvc, backofficeUserSvc, matchingSvc, matchingValidator := setupServices(cfg)
	go func() {
		server := httpserver.New(cfg,
			authSvc,
			userSvc,
			authorizationSvc,
			userValidator,
			backofficeUserSvc,
			matchingSvc,
			matchingValidator)

		server.Serve()
	}()

	done := make(chan bool)

	go func() {
		sch := scheduler.New()
		sch.Start(done)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("received signal interrupt . shutting down gracefully...")
	done <- true
	time.Sleep(5 * time.Second)
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
