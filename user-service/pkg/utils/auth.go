package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"user-service/pkg/types"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

func GenerateJWT(secret string, claims types.JwtCustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func hashHelper(password string, saltBytes []byte) []byte {
	const (
		saltSize    = 16
		hashSize    = 32
		iterations  = 3
		memory      = 64 * 1024
		parallelism = 4
	)

	return argon2.IDKey([]byte(password), saltBytes, iterations, memory, parallelism, hashSize)
}

// HashPassword generates a hashed password and a salt
func HashPassword(password string) (string, string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", "", err
	}

	// Hash the password with Argon2
	hashedPassword := hashHelper(password, salt)

	return fmt.Sprintf("%x", hashedPassword), fmt.Sprintf("%x", salt), nil
}

// CheckPassword verifies if the provided password matches the hashed password and salt
func CheckPassword(password string, hashedPassword string, salt string) bool {
	// Decode the salt from hex
	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		return false
	}

	// Hash the provided password with the stored salt
	newHashedPassword := hashHelper(password, saltBytes)
	passwordBytes, _ := hex.DecodeString(hashedPassword)
	// Compare the generated hash with the stored hashed password
	return bytes.Equal(newHashedPassword, passwordBytes)
}
