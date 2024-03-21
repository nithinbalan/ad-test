package repositories

import (
	"errors"
	"jethur/grc/organization/models"
	"jethur/utility"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// Interfaces
type AuthRepositoryInterface interface {
	Authenticate(email string, password string) (map[string]string, error)
	GetOtpForgotPasswordRecovery(email string) (*models.User, error)
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repo *AuthRepository) Authenticate(email, password string) (map[string]string, error) {

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

	//generate access token
	accessToken, err := utility.GenerateToken(user.ID)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	refreshToken, err := utility.GenerateRefreshToken(user.ID)
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

func (repo *AuthRepository) GetOtpForgotPasswordRecovery(email string) (*models.User, error) {

	var user *models.User

	result := repo.db.First(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, utility.NewInternalServerError(result.Error)
	}
	otpLength := 6
	otp := utility.GenerateOTP(otpLength)
	user.Otp = otp
	user.SentAt = time.Now()

	// Update the user model
	if err := repo.db.Save(&user).Error; err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	return user, nil

}
