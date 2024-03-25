package repositories

import (
	"errors"
	"fmt"
	"jethur/grc/organization/models"
	"jethur/utility"
	"time"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/gorm"
)

// Interfaces
type AuthRepositoryInterface interface {
	Authenticate(email string, password string) (map[string]string, error)
	SendPasswordRecoveryOTP(email string) (*models.User, error)
	ReSendPasswordRecoveryOTP(email string) (*models.User, error)
	VerifyPasswordRecoveryOTP(email, otp string) (*models.User, error)
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

func (repo *AuthRepository) SendPasswordRecoveryOTP(email string) (*models.User, error) {

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
	user.SentAt = time.Now().UTC()

	if err := repo.db.Save(&user).Error; err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	return user, nil

}
func (repo *AuthRepository) ReSendPasswordRecoveryOTP(email string) (*models.User, error) {

	var user *models.User

	result := repo.db.First(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if result.Error != nil {
		return nil, utility.NewInternalServerError(result.Error)
	}

	resendTime := user.SentAt.Add(1 * time.Minute)
	currentTime := time.Now().UTC()
	if currentTime.Before(resendTime) {
		return nil, utility.NewValidationError(fmt.Errorf("resend time has not reached yet"), "resend_time_not_reached", "resend_time")
	}

	otpLength := 6
	otp := utility.GenerateOTP(otpLength)
	user.Otp = otp
	user.SentAt = time.Now().UTC()

	if err := repo.db.Save(&user).Error; err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	return user, nil
}

func (repo *AuthRepository) VerifyPasswordRecoveryOTP(email, otp string) (*models.User, error) {
	var user models.User

	// Get the user with the given email and OTP
	result := repo.db.Where("email = ? AND otp = ?", email, otp).First(&user)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, email)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
