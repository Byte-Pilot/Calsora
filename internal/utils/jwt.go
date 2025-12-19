package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strconv"
	"time"
)

type UserJWTClaims struct {
	UserID       int       `json:"user_id"`
	Subscription string    `json:"subscription"`
	ExpiresAt    time.Time `json:"expires_at"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID int, subscription string, expiresAt time.Time) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET_ACCESS"))

	claims := UserJWTClaims{
		UserID:       userID,
		Subscription: subscription,
		ExpiresAt:    expiresAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func GenerateRefreshToken(userID int) (string, time.Time, error) {
	secret := []byte(os.Getenv("JWT_SECRET_REFRESH"))
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString(secret)
	return refreshToken, expiresAt, err
}
