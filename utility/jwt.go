package utility

import (
	"errors"
	"jethur/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenDuration  = 30 * 24 * time.Hour // 1 month
	RefreshTokenDuration = 7 * 24 * time.Hour  // 7 days
)

func getSecretKey() string {
	return config.GetEnv("ENCRYPTION_KEY", "")
}

type UserClaims struct {
	UserID uint `json:"user_id"`
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
