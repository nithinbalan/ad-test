package utility

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// DecryptJWTToken decrypts a JWT token and returns its claims if valid
func DecryptJWTToken(tokenString string) (jwt.MapClaims, error) {
	// Parse and validate JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("f0a3e3b61e7082367820944efeb9d81c7e242ec4d22dc6fa941aaf76d333f158"), nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
