package config

import (
	"os"
	"strconv"
	"time"
)

var (
	JWTSecretKey    string
	TokenExpiration time.Duration
)

func InitJWTConfig() {
	JWTSecretKey = os.Getenv("JWT_SECRET_KEY")

	expMinutes, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_MINUTES"))
	if err != nil {
		expMinutes = 15 // fallback par défaut
	}
	TokenExpiration = time.Duration(expMinutes) * time.Minute
}