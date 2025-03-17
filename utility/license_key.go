package utility

import (
	"encoding/json"
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
		return nil, NewValidationError(fmt.Errorf("invalid token"), "invalid", "license key")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, NewValidationError(fmt.Errorf("invalid token claims"), "invalid", "token claims")
	}

	// Extract and decode enabled_modules
	if modules, exists := claims["enabled_modules"]; exists {
		var enabledModules []string
		if str, ok := modules.(string); ok {
			if err := json.Unmarshal([]byte(str), &enabledModules); err != nil {
				return nil, NewValidationError(fmt.Errorf("failed to parse enabled_modules"), "invalid", "enabled modules")
			}
		} else {
			return nil, NewValidationError(fmt.Errorf("invalid enabled_modules format"), "invalid", "enabled modules")
		}

		// Update claims with parsed enabled_modules
		claims["enabled_modules"] = enabledModules
	}

	return claims, nil
}
