package utils

import (
	"user-service/pkg/types"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secret string, claims types.JwtCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
