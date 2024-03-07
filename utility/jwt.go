package utility

import (
	"errors"
	"jethur/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour // 7 days, adjust as needed
)

func getSecretKey() string {
	return config.GetEnv("ENCRYPTION_KEY", "")
}

type UserClaims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}

func GenerateToken(userID uint) (string, error) {
	return generateTokenWithExpiry(userID, AccessTokenDuration)
}

func GenerateRefreshToken(userID uint) (string, error) {
	return generateTokenWithExpiry(userID, RefreshTokenDuration)
}

func generateTokenWithExpiry(userID uint, duration time.Duration) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := getSecretKey() // Retrieve the secret key
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(tokenStr string) (*UserClaims, error) {
	secretKey := getSecretKey() // Retrieve the secret key for verification
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method in token")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or claims")
	}

	return claims, nil
}
