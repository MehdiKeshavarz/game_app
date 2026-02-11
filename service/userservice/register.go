package userservice

import (
	"fmt"
	"game_app/entity"
	"game_app/param"

	"golang.org/x/crypto/bcrypt"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {
	// TODO - we verify phone number by verification code

	hashedPassword, hErr := hashPassword(req.Password)
	if hErr != nil {
		return param.RegisterResponse{}, fmt.Errorf("unexpected error: %w", hErr)
	}

	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
		Role:        entity.UserRole,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return param.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return param.RegisterResponse{
		User: struct {
			ID          uint   `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		}{ID: createdUser.ID, Name: createdUser.Name, PhoneNumber: createdUser.PhoneNumber},
	}, nil
}

func hashPassword(password string) (string, error) {
	hashedPass, hErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hashedPass), hErr
}
