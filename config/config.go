package config

import (
	"game_app/adapter/redis"
	"game_app/repository/mysql"
	"game_app/service/authservice"
	"game_app/service/matchingservice"
	"time"
)

type Application struct {
	GracefullyShutdownTimeout time.Duration `koanf:"gracefully_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}
type Config struct {
	HttpServer      HTTPServer             `koanf:"http_server"`
	Application     Application            `koanf:"application"`
	Auth            authservice.Config     `koanf:"auth"`
	Mysql           mysql.Config           `koanf:"mysql"`
	MatchingService matchingservice.Config `koanf:"matching_service"`
	Redis           redis.Config           `koanf:"redis"`
}
