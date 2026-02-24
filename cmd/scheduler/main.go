package main

import (
	"fmt"
	"game_app/adapter/redis"
	"game_app/config"
	"game_app/repository/redis/matching"
	"game_app/scheduler"
	"game_app/service/matchingservice"
	"os"
	"os/signal"
)

func main() {
	cfg := config.Load()

	done := make(chan bool)

	matchingSvc := setupServices(cfg)

	go func() {
		sch := scheduler.New(matchingSvc)
		sch.Start(done)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("received signal interrupt . shutting down gracefully...")
	done <- true

}

func setupServices(cfg config.Config) matchingservice.Service {
	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := matching.New(redisAdapter)
	matchingSvc := matchingservice.New(matchingRepo, cfg.MatchingService)

	return matchingSvc
}
