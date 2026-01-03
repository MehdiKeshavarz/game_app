package userservice

import (
	"fmt"
	"game_app/param"
	"game_app/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	const op = "userservice.Login"
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).SetWrappedError(err)
	}

	cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if cErr != nil {
		return param.LoginResponse{}, fmt.Errorf("username or password is't correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	return param.LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil

}
