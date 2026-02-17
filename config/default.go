package config

var defaultConfig = map[string]interface{}{
	"auth.access_subject":     AccessSubject,
	"auth.refresh_subject":    RefreshSubject,
	"access_expiration_time":  AccessExpirationTime,
	"refresh_expiration_time": RefreshExpirationTime,
}
