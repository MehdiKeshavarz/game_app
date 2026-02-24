package scheduler

import (
	"fmt"
	"game_app/param"
	"game_app/service/matchingservice"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	sch         gocron.Scheduler
	matchingSvc matchingservice.Service
}

func New(matchingSvc matchingservice.Service) Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	return Scheduler{
		sch:         s,
		matchingSvc: matchingSvc,
	}
}

func (s Scheduler) Start(done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := s.sch.NewJob(
		gocron.DurationJob(5*time.Second),
		gocron.NewTask(
			s.MatchWaitedUser,
		),
	)

	if err != nil {
		fmt.Println("cron job err:", err)
		return

	}

	s.sch.Start()

	<-done
	fmt.Println("exiting...")
	sErr := s.sch.StopJobs()

	if sErr != nil {
		fmt.Println("shutdown err:", sErr)
		return
	}

}

func (s Scheduler) MatchWaitedUser() {
	fmt.Println("schedule match waited user")
	_, _ = s.matchingSvc.MatchWaitedUser(param.MatchWaitedUserRequest{})

}
