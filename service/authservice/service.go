package authservice

import (
	"fmt"
	"game_app/entity"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	SignKey               string        `koanf:"sign_key"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`
	AccessSubject         string        `koanf:"access_subject"`
	RefreshSubject        string        `koanf:"refresh_subject"`
}

type Service struct {
	config Config
}

func New(config Config) Service {
	return Service{
		config: config,
	}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.creatToken(user.ID, s.config.AccessSubject, time.Hour*24)
}
func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.creatToken(user.ID, s.config.RefreshSubject, time.Hour*7*24)
}

func (s Service) ParseToken(tokenStr string) (*Claims, error) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		fmt.Printf("%v %v", claims.UserID, claims.ExpiresAt)
		return claims, nil
	} else {
		return nil, err
	}
}

func (s Service) creatToken(userID uint, subject string, expireDuration time.Duration) (string, error) {
	// set our claims
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
		},
		UserID: userID,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(s.config.SignKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
