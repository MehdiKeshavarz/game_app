package main

import (
	"fmt"
	"game_app/config"
	"game_app/delivery/httpserver"
	"game_app/repository/migrator"
	"game_app/repository/mysql"
	"game_app/repository/mysql/accesscontrol"
	"game_app/repository/mysql/user"
	"game_app/service/authorizationservice"
	"game_app/service/authservice"
	"game_app/service/backofficeuserservice"
	"game_app/service/userservice"
	"game_app/validator/uservalidator"
	"time"
)

const (
	JwtSignKey = "jwt_secret"
)

func main() {
	cfg2 := config.Load()
	fmt.Printf("cfg2 : %v\n", cfg2)
	// TODO - marge cfg with cfg2
	cfg := config.Config{
		HttpServer: config.HTTPServer{Port: 8088},
		Auth: authservice.Config{
			SignKey:               JwtSignKey,
			AccessExpirationTime:  time.Hour * 24,
			RefreshExpirationTime: time.Hour * 24 * 7,
			AccessSubject:         "at",
			RefreshSubject:        "rt",
		},
		Mysql: mysql.Config{
			Host:     "localhost",
			Port:     3308,
			Username: "gameapp",
			Password: "gameappt0lk2o20",
			DBName:   "gameapp_db",
		},
	}

	mgr := migrator.New(cfg.Mysql, "mysql")
	mgr.Up()

	authSvc, userSvc, userValidator, authorizationSvc, backofficeUserSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, authorizationSvc, userValidator, backofficeUserSvc)

	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator, authorizationservice.Service, backofficeuserservice.Service) {
	authSvc := authservice.New(cfg.Auth)

	mysqlRepo := mysql.New(cfg.Mysql)

	userMysql := user.New(mysqlRepo)

	userValidator := uservalidator.New(userMysql)
	userSvc := userservice.New(userMysql, authSvc)

	backofficeUserSvc := backofficeuserservice.New()

	accesscontrolMysql := accesscontrol.New(mysqlRepo)

	authorizationSvc := authorizationservice.New(accesscontrolMysql)

	return authSvc, userSvc, userValidator, authorizationSvc, backofficeUserSvc
}
