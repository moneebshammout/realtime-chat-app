package types
import (
	"github.com/golang-jwt/jwt/v5"
)


type JwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.RegisteredClaims
}

