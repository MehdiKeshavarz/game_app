package authservice

import (
	"fmt"
	"game_app/entity"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	signKey               string
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	accessSubject         string
	refreshSubject        string
}

func New(signKey, accessSubject, refreshSubject string,
	accessExpirationTime, refreshExpirationTime time.Duration) Service {
	return Service{
		signKey:               signKey,
		accessExpirationTime:  accessExpirationTime,
		refreshExpirationTime: refreshExpirationTime,
		accessSubject:         accessSubject,
		refreshSubject:        refreshSubject,
	}
}

func (s Service) CreateAccessToken(user entity.User) (string, error) {
	return s.creatToken(user.ID, s.accessSubject, time.Hour*24)
}
func (s Service) CreateRefreshToken(user entity.User) (string, error) {
	return s.creatToken(user.ID, s.refreshSubject, time.Hour*7*24)
}

func (s Service) ParseToken(tokenStr string) (*Claims, error) {
	tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.signKey), nil
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
	tokenString, err := accessToken.SignedString([]byte(s.signKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
