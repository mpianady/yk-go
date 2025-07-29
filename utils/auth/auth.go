package auth

import (
	"errors"
	"go-blog/models/user"
	"go-blog/services/config"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func ValidateCredentials(email, password string) (*user.User, error) {
	var userModel user.User
	if err := config.Db.Where("email = ?", email).First(&userModel).Error; err != nil {
		return nil, err
	}

	if userModel.Status != string(user.StatusActive) {
		return nil, errors.New("user is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(password)); err != nil {
		return nil, err
	}

	return &userModel, nil
}

func ExtractBearerToken(authHeader string) (string, bool) {
	const bearerPrefix = "Bearer "
	if authHeader == "" {
		return "", false
	}
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", false
	}
	token := strings.TrimPrefix(authHeader, bearerPrefix)
	if token == "" {
		return "", false
	}
	return token, true
}
