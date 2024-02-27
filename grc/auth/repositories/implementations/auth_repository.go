package implementations

import (
	"errors"
	"jethur/config"
	"jethur/grc/organization/models"
	"jethur/utility"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repo *AuthRepository) Authenticate(c *gin.Context, email, password string) (map[string]string, error) {

	var user models.User

	result := repo.db.First(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("invalid credentials")
	} else if result.Error != nil {
		return nil, utility.NewInternalServerError(result.Error)
	}

	// Compare the provided password with the stored hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	jwtKey := config.GetEnv("ENCRYPTION_KEY", "")
	jwtManager := utility.NewJWTManager(jwtKey, 24*time.Hour)

	//generate access token
	accessToken, err := jwtManager.GenerateToken(user.ID)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	refreshToken, err := jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	// Return both access and refresh tokens
	authTokens := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return authTokens, nil

}
