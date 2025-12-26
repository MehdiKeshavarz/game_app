package userservice

import (
	"fmt"
	"game_app/dto"
	"game_app/entity"
	"game_app/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
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

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// TODO - we verify phone number by verification code

	hashedPassword, hErr := hashPassword(req.Password)
	if hErr != nil {
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %w", hErr)
	}

	// create new user in storage
	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return dto.RegisterResponse{
		User: struct {
			ID          uint   `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		}{ID: createdUser.ID, Name: createdUser.Name, PhoneNumber: createdUser.PhoneNumber},
	}, nil
}

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

func (s Service) GetProfile(req dto.GetProfileRequest) (dto.GetProfileResponse, error) {
	const op = "userservice.GetProfile"
	// I don't expect the repository call return "record not found " error,
	// because I assume the interactor input is sanitized.

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return dto.GetProfileResponse{}, richerror.New(op).
			SetWrappedError(err).
			SetMeta(map[string]interface{}{"req": req})
	}

	return dto.GetProfileResponse{Name: user.Name}, nil
}

func hashPassword(password string) (string, error) {
	hashedPass, hErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hashedPass), hErr
}
