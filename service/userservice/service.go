package userservice

import (
	"fmt"
	"game_app/entity"
	"game_app/pkg/phonenumber"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/golang-jwt/jwt/v5"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(user entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserByID(userID uint) (entity.User, error)
}

type Service struct {
	repo    Repository
	signKey string
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func NewService(repo Repository, signKey string) Service {
	return Service{
		repo:    repo,
		signKey: signKey,
	}
}

type RegisterResponse struct {
	entity.User
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
		Password:    string(hashedPassword),
	}

	createdUser, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	// return created user
	return RegisterResponse{createdUser}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	if !exist {
		return LoginResponse{}, fmt.Errorf("username or password is't correct")
	}

	// compare  user.password with the req.password

	cErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if cErr != nil {
		return LoginResponse{}, fmt.Errorf("username or password is't correct")
	}

	// jwt token generate
	token, err := creatToken(user.ID, s.signKey)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return LoginResponse{AccessToken: token}, nil

}

type GetProfileRequest struct {
	UserID uint `json:"user_id"`
}
type GetProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) GetProfile(req GetProfileRequest) (GetProfileResponse, error) {
	// I don't expect the repository call return "record not found " error,
	// because I assume the interactor input is sanitized.

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return GetProfileResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return GetProfileResponse{user.Name}, nil
}

func hashPassword(password string) (string, error) {
	hashedPass, hErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	return string(hashedPass), hErr
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint
}

func (c Claims) Valid() error {
	return nil
}

func creatToken(userID uint, signKey string) (string, error) {
	// set our claims
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// set the expire time
			// see https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		UserID: userID,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(signKey))

	// Creat token string
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
