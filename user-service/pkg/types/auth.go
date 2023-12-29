package types
import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthConfig struct {
	SigningKey   string
	TokenLookup  string
	PublicRoutes []string
}

type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

