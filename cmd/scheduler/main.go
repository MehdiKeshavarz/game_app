package main

import (
	"fmt"
	"game_app/scheduler"
	"os"
	"os/signal"
)

func main() {
	//cfg := config.Load()

	done := make(chan bool)

	go func() {
		sch := scheduler.New()
		sch.Start(done)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("received signal interrupt . shutting down gracefully...")
	done <- true

}
