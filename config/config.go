package config

import (
	"game_app/repository/mysql"
	"game_app/service/authservice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}
type Config struct {
	HttpServer HTTPServer         `koanf:"http_server"`
	Auth       authservice.Config `koanf:"auth"`
	Mysql      mysql.Config       `koanf:"mysql"`
}
