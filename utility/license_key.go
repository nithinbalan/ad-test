package utility

import (
	"encoding/json"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var publicKeyPEM = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwaJLxaJUkd2JoDdA6ls7
XTL82Mpc/MKLw0bXSJe8bTKDWJ9q1SADXAdHlU3aGgMBjmgDdn0/EE6xE9dfA08P
QIw55/yxoNdhi36yT7IbZ3TbJnbw85+hrL3cUvfp/VueXlqta2LckL4m93a/5DSC
s1KJMTaKuPEkv+ahnA1xXjJnrm1PsVoU8iEEGUOQegXKxvxIDFAyiL6ywsf68h96
Sc/RMDbTPe1Nt9obUv1tEZJqbUiecqTVt+PUb4gg6GvedwsL1ePqOynTybHKbjKi
DuYQlhZJzBu1UxtKAne0F1pLieHwpHs1gXOdT6ht3/MjePQ/rBC4mds9EwhZjI1w
XwIDAQAB
-----END PUBLIC KEY-----`)

func DecryptJWTToken(tokenString string) (jwt.MapClaims, error) {
	// Parse RSA public key from PEM
	// Read the PEM file from the root path
	// pemData, err := os.ReadFile("license_public.pem")
	// if err != nil {
	// 	return nil, NewValidationError(fmt.Errorf("failed to read public key file: %v", err), "invalid", "public key file")
	// }

	// Parse the RSA public key
	pubKeyParsed, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyPEM)
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
