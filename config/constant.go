package config

import "time"

const (
	AuthMiddlewareContextKey = "claims"
	JwtSignKey               = "jwt_secret"
	AccessExpirationTime     = time.Hour * 24
	RefreshExpirationTime    = time.Hour * 24 * 7
	AccessSubject            = "at"
	RefreshSubject           = "rt"
)
