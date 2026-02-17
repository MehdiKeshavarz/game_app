package config

import "time"

const (
	AuthMiddlewareContextKey = "claims"
	AccessExpirationTime     = time.Hour * 24
	RefreshExpirationTime    = time.Hour * 24 * 7
	AccessSubject            = "at"
	RefreshSubject           = "rt"
)
