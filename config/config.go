package config

import (
	"game_app/repository/mysql"
	"game_app/service/authservice"
)

type HTTPServer struct {
	Port int
}
type Config struct {
	HttpServer HTTPServer
	AuthConfig authservice.Config
	Mysql      mysql.Config
}
