package presenceservice

import (
	"context"
	"fmt"
	"game_app/param"
	"game_app/pkg/richerror"
	"time"
)

type Config struct {
	PresencePrefix string        `koanf:"prefix"`
	ExpirationTime time.Duration `koanf:"expiration_time"`
}

type Repo interface {
	Upsert(ctx context.Context, key string, timestamp int64, expirationTime time.Duration) error
}

type Service struct {
	repo Repo
	cfg  Config
}

func New(repo Repo, cfg Config) Service {
	return Service{
		repo: repo,
		cfg:  cfg,
	}
}

func (s Service) UpsertPresence(ctx context.Context, req param.UpsertPresenceRequest) (param.UpsertPresenceResponse, error) {
	const op = richerror.Op("presenceservice.UpsertPresence")
	key := fmt.Sprintf("%s:%d", s.cfg.PresencePrefix, req.UserID)
	err := s.repo.Upsert(ctx, key, req.Timestamp, s.cfg.ExpirationTime)
	if err != nil {
		return param.UpsertPresenceResponse{}, richerror.New(op).
			SetWrappedError(err).
			SetKind(richerror.KindUnexpected)
	}
	return param.UpsertPresenceResponse{}, nil
}
