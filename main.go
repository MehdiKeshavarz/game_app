package main

import (
	"game_app/config"
	"game_app/delivery/httpserver"
	"game_app/repository/mysql"
	"game_app/service/authservice"
	"game_app/service/userservice"
	"time"
)

const (
	JwtSignKey = "jwt_secret"
)

func main() {

	cfg := config.Config{
		HttpServer: config.HTTPServer{Port: 8080},
		AuthConfig: authservice.Config{
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

	//	mgr := migrator.New(cfg.Mysql, "mysql")
	//	mgr.Up()

	authSvc, userSvc := setupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc)

	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.AuthConfig)
	mysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(mysqlRepo, authSvc)

	return authSvc, userSvc
}
