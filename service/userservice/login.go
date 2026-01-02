package userservice

import (
	"fmt"
	"game_app/dto"
	"game_app/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(req dto.LoginRequest) (dto.LoginResponse, error) {
	const op = "userservice.Login"
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return dto.LoginResponse{}, richerror.New(op).SetWrappedError(err)
	}
	if !exist {
		return dto.LoginResponse{}, fmt.Errorf("username or password is't correct")
	}

	cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if cErr != nil {
		return dto.LoginResponse{}, fmt.Errorf("username or password is't correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return dto.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	return dto.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil

}
