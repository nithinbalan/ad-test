package utility

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"jethur/config"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenDuration         = 30 * 24 * time.Hour // 1 month
	RefreshTokenDuration        = 7 * 24 * time.Hour  // 7 days
	ExternalAccessTokenDuration = 7 * 24 * time.Hour  // 7 days
)

func getSecretKey() string {
	return config.GetEnv("ENCRYPTION_KEY", "")
}

func getExternalSecretKey() string {
	key := config.GetEnv("EXTERNAL_ENCRYPTION_KEY", "")
	if key == "" {
		LogFatal("EXTERNAL_ENCRYPTION_KEY not found in environment variables", nil)
	}
	return key
}

type UserClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type ExternalClientClaims struct {
	ClientID string `json:"client_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint) (string, error) {
	return generateTokenWithExpiry(userID, AccessTokenDuration)
}

func GenerateRefreshToken(userID uint) (string, error) {
	return generateTokenWithExpiry(userID, RefreshTokenDuration)
}

func generateTokenWithExpiry(userID uint, duration time.Duration) (string, error) {
	now := time.Now()

	claims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(getSecretKey()))
}

func VerifyToken(tokenStr string) (*UserClaims, error) {
	secretKey := getSecretKey()

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method in token")
		}
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or claims")
	}

	return claims, nil
}

// GenerateRandomString creates a secure random string for OAuth2 state/nonce.
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// fallback if crypto/rand fails
		return base64.URLEncoding.EncodeToString([]byte("fallback-random"))[:length]
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

// DecodeJWTClaims decodes a JWT (ID token or access token) without verifying signature.
// Used to inspect claims from Microsoft Entra ID.
func DecodeJWTClaims(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid JWT token format")
	}

	payload := parts[1]
	// Add missing padding if required
	if l := len(payload) % 4; l > 0 {
		payload += strings.Repeat("=", 4-l)
	}

	decoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		// try standard encoding as fallback
		decoded, err = base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return nil, errors.New("failed to decode JWT payload")
		}
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, errors.New("failed to unmarshal JWT claims: " + err.Error())
	}

	return claims, nil
}

func GenerateExternalToken(clientID string) (string, error) {
	return generateExternalTokenWithExpiry(clientID, ExternalAccessTokenDuration)
}

func generateExternalTokenWithExpiry(clientID string, duration time.Duration) (string, error) {
	now := time.Now()

	claims := ExternalClientClaims{
		ClientID: clientID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(getExternalSecretKey()))
}

func VerifyExternalToken(tokenStr string) (*ExternalClientClaims, error) {
	secretKey := getExternalSecretKey()

	token, err := jwt.ParseWithClaims(tokenStr, &ExternalClientClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method in token")
		}
		return []byte(secretKey), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*ExternalClientClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or claims")
	}

	return claims, nil
}
