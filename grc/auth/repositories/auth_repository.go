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
	SetPassword(userID uint, newPassword string) (*models.User, error)
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
	var userSecurity models.UserSecurity

	result := repo.db.First(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("invalid credentials")
	} else if result.Error != nil {
		return nil, utility.NewInternalServerError(result.Error)
	}

	// Next, find the user's security information by user ID
	result = repo.db.Where("user_id = ?", user.ID).First(&userSecurity)
	if result.Error != nil {
		// Consider how you want to handle errors here, this could be a missing record or a database error
		return nil, utility.NewInternalServerError(result.Error)
	}

	// Compare the provided password with the stored hashed password
	err := bcrypt.CompareHashAndPassword([]byte(userSecurity.PasswordHash), []byte(password))
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
		return nil, utility.NewResourceNotFoundError(result.Error, email)
	} else if result.Error != nil {
		return nil, utility.NewInternalServerError(result.Error)
	}
	otpLength := 6
	otp := utility.GenerateOTP(otpLength)
	currentTime := time.Now().UTC()
	user.Otp = &otp
	user.SentAt = &currentTime
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
	user.Otp = &otp
	user.SentAt = &currentTime

	if err := repo.db.Save(&user).Error; err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	return user, nil
}

func (repo *AuthRepository) VerifyPasswordRecoveryOTP(email, otp string) (*models.User, error) {
	var user models.User

	// Get the user with the given email and OTP
	result := repo.db.First(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, email)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (repo *AuthRepository) SetPassword(userId uint, newPassword string) (*models.User, error) {

	var user *models.User

	result := repo.db.First(&user, userId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, userId)
	} else if result.Error != nil {
		return nil, utility.NewInternalServerError(result.Error)
	}
	// Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}
	// Set the hashed password
	user.PasswordHash = string(hashedPassword)

	if err := repo.db.Save(&user).Error; err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	return user, nil

}
