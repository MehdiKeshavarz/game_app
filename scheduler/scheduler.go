package scheduler

import (
	"fmt"
	"time"
)

type Scheduler struct {
}

func New() Scheduler {
	return Scheduler{}
}

func (s Scheduler) Start(done <-chan bool) {

	for {
		select {
		case <-done:
			fmt.Println("exiting...")
			return
		default:
			now := time.Now()
			fmt.Println("start Schedule this time ", now)
			time.Sleep(500 * time.Millisecond)
		}

	}
}
