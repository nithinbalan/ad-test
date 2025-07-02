package utility

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func DecryptJWTToken(tokenString string) (jwt.MapClaims, error) {
	// Parse RSA public key from PEM
	// Read the PEM file from the root path
	pemData, err := os.ReadFile("license_public.pem")
	if err != nil {
		return nil, NewValidationError(fmt.Errorf("failed to read public key file: %v", err), "invalid", "public key file")
	}

	// Parse the RSA public key
	pubKeyParsed, err := jwt.ParseRSAPublicKeyFromPEM(pemData)
	if err != nil {
		return nil, NewValidationError(fmt.Errorf("invalid public key: %v", err), "invalid", "public key")
	}

	// Parse and verify the JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKeyParsed, nil
	})

	if err != nil || !token.Valid {
		return nil, NewValidationError(fmt.Errorf("invalid token"), "invalid", "license key")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, NewValidationError(fmt.Errorf("invalid token claims"), "invalid", "token claims")
	}

	// Parse enabled_modules
	if modules, exists := claims["enabled_modules"]; exists {
		var enabledModules []string
		if str, ok := modules.(string); ok {
			if err := json.Unmarshal([]byte(str), &enabledModules); err == nil {
				claims["enabled_modules"] = enabledModules
			}
		}
	}

	// Parse enabled_addons
	if addons, exists := claims["enabled_addons"]; exists {
		var enabledAddons []string
		if str, ok := addons.(string); ok {
			if err := json.Unmarshal([]byte(str), &enabledAddons); err == nil {
				claims["enabled_addons"] = enabledAddons
			}
		}
	}

	return claims, nil
}
