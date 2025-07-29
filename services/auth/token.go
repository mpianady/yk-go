package auth

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go-blog/models/auth"
	"go-blog/services/config"
	"time"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}


func (ts *TokenService) GenerateAccessToken(email string) (string, error) {
	return GenerateToken(email)
}

func (ts *TokenService) GenerateRefreshToken(userID uint) (string, error) {
	tokenString, err := ts.generateRandomString(32)
	if err != nil {
		return "", err
	}

	refreshToken := auth.RefreshToken{
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7 days
	}

	if err := config.Db.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

func (ts *TokenService) ValidateRefreshToken(tokenString string) (*auth.RefreshToken, error) {
	var refreshToken auth.RefreshToken
	err := config.Db.Where("token = ?", tokenString).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	return &refreshToken, nil
}

func (ts *TokenService) DeleteRefreshToken(tokenString string) error {
	return config.Db.Where("token = ?", tokenString).Delete(&auth.RefreshToken{}).Error
}

func (ts *TokenService) ParseAndValidateAccessToken(tokenString string) (*auth.Claims, error) {
	claims := &auth.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func (ts *TokenService) generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GenerateToken(email string) (string, error) {
	expirationTime := time.Now().Add(config.TokenExpiration)

	claims := &auth.Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JWTSecretKey))
}
