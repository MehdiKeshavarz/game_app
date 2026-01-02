package userservice

import (
	"game_app/entity"
)

type Repository interface {
	Register(user entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserByID(userID uint) (entity.User, error)
}

type AuthService interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	auth AuthService
	repo Repository
}

func New(repo Repository, auth AuthService) Service {
	return Service{
		repo: repo,
		auth: auth,
	}
}
