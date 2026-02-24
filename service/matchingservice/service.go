package matchingservice

import (
	"game_app/entity"
	"game_app/param"
	"game_app/pkg/richerror"
	"time"
)

type Repo interface {
	AddToWaitingList(userID uint, category entity.Category) error
}

type Config struct {
	WaitingTimeOut time.Duration `koanf:"waitingTimeOut"`
}

type Service struct {
	repo   Repo
	config Config
}

func New(repo Repo, config Config) Service {
	return Service{
		repo:   repo,
		config: config,
	}
}

func (s Service) AddToWaitingList(req param.AddToWaitingListRequest) (param.AddToWaitingListResponse, error) {
	const op = "matchingservice.AddToWaitingList"

	err := s.repo.AddToWaitingList(req.UserID, req.Category)

	if err != nil {
		return param.AddToWaitingListResponse{}, richerror.New(op).SetWrappedError(err).
			SetKind(richerror.KindUnexpected)
	}

	return param.AddToWaitingListResponse{
		TimeOut: s.config.WaitingTimeOut,
	}, nil
}

func (s Service) MatchWaitedUser(req param.MatchWaitedUserRequest) (param.MatchWaitedUserResponse, error) {

	return param.MatchWaitedUserResponse{}, nil
}
