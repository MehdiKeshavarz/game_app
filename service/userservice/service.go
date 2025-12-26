package userservice

import (
	"fmt"
	"game_app/entity"
	"game_app/pkg/phonenumber"
	"game_app/pkg/richerror"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
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

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
	} `json:"user"`
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO - we verify phone number by verification code

	// validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// check uniqueness of phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
		}

		return RegisterResponse{}, fmt.Errorf("phone number is not unique")
	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name should be greater than 3 characters")
	}

	// 	TODO -check the password with regex pattern
	// validate password
	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("password should be greater than 8 characters")
	}

	hashedPassword, hErr := hashPassword(req.Password)
	if hErr != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", hErr)
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
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return RegisterResponse{
		User: struct {
			ID          uint   `json:"id"`
			Name        string `json:"name"`
			PhoneNumber string `json:"phone_number"`
		}{ID: createdUser.ID, Name: createdUser.Name, PhoneNumber: createdUser.PhoneNumber},
	}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	const op = "userservice.Login"
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, richerror.New(op).SetWrappedError(err)
	}
	if !exist {
		return LoginResponse{}, fmt.Errorf("username or password is't correct")
	}

	cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if cErr != nil {
		return LoginResponse{}, fmt.Errorf("username or password is't correct")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	return LoginResponse{AccessToken: accessToken, RefreshToken: refreshToken}, nil

}

type GetProfileRequest struct {
	UserID uint `json:"user_id"`
}
type GetProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) GetProfile(req GetProfileRequest) (GetProfileResponse, error) {
	const op = "userservice.GetProfile"
	// I don't expect the repository call return "record not found " error,
	// because I assume the interactor input is sanitized.

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return GetProfileResponse{}, richerror.New(op).
			SetWrappedError(err).
			SetMeta(map[string]interface{}{"req": req})
	}

	return GetProfileResponse{user.Name}, nil
}

func hashPassword(password string) (string, error) {
	hashedPass, hErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hashedPass), hErr
}
