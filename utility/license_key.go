package utility

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// DecryptJWTToken decrypts a JWT token and returns its claims if valid
func DecryptJWTToken(tokenString string) (jwt.MapClaims, error) {
	// Parse and validate JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("f0a3e3b61e7082367820944efeb9d81c7e242ec4d22dc6fa941aaf76d333f158"), nil
	})

	if err != nil {
		return nil, NewValidationError(fmt.Errorf("invalid token"), "not_exist", "license key")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, NewValidationError(fmt.Errorf("invalid token claims"), "not_exist", "token claims")
	}

	return claims, nil
}
