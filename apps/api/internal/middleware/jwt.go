package middleware

import (
	"time"

	"monji/internal/models"

	"github.com/golang-jwt/jwt/v4"
)

// JWTSecret is used to sign tokens.
// In production, load this from configuration.
var JWTSecret = "supersecretkey"

// GenerateJWT creates a new JWT token for the given user.
func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(JWTSecret))
}
