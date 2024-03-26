package services

import (
	"fmt"
	"jethur/config"
	authDTOs "jethur/grc/auth/dtos"
	"jethur/grc/auth/repositories"
	"jethur/utility"
	"strconv"
	"time"
)

type AuthServiceInterface interface {
	Authenticate(email string, password string) (map[string]string, error)
	SendPasswordRecoveryOTP(email string) error
	ReSendPasswordRecoveryOTP(email string) error
	VerifyPasswordRecoveryOTP(email, otp string) (string, error)
}

// Service and implementations
type AuthService struct {
	repo repositories.AuthRepositoryInterface
}

func NewAuthService(repo repositories.AuthRepositoryInterface) AuthServiceInterface {
	return &AuthService{
		repo: repo,
	}
}

func (as *AuthService) Authenticate(email string, password string) (map[string]string, error) {
	return as.repo.Authenticate(email, password)
}

func (as AuthService) SendPasswordRecoveryOTP(email string) error {
	user, err := as.repo.SendPasswordRecoveryOTP(email)

	if err != nil {
		return err
	}

	emailData := authDTOs.PasswordRecoveryDataDTO{
		FirstNameEn: user.FirstNameEN,
		LastNameEn:  user.LastNameEN,
		Otp:         user.Otp,
	}

	to := []string{user.Email}
	subject := "Your Password Reset Verification Code"
	templateRelativePath := "auth/templates/email/forgot_password.html"

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, nil)
	fmt.Println("err", err)
	if err != nil {
		return err
	}

	return nil

}
func (as AuthService) ReSendPasswordRecoveryOTP(email string) error {

	user, err := as.repo.ReSendPasswordRecoveryOTP(email)

	if err != nil {
		return err
	}

	emailData := authDTOs.PasswordRecoveryDataDTO{
		FirstNameEn: user.FirstNameEN,
		LastNameEn:  user.LastNameEN,
		Otp:         user.Otp,
	}

	to := []string{user.Email}
	subject := "Your Password Reset Verification Code"
	templateRelativePath := "auth/templates/email/forgot_password.html"

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, nil)
	if err != nil {
		return err
	}

	return nil
}
func (as *AuthService) VerifyPasswordRecoveryOTP(email, otp string) (string, error) {

	user, err := as.repo.VerifyPasswordRecoveryOTP(email, otp)
	if err != nil {
		return "", err
	}

	// Calculate expiration time based on configuration
	verifyOTPDurationStr := config.GetEnv("VERIFY_OTP_DURATION", "10")
	verifyOTPDuration, _ := strconv.Atoi(verifyOTPDurationStr)

	expirationTime := user.SentAt.Add(time.Duration(verifyOTPDuration) * time.Minute)
	currentTime := time.Now().UTC()

	if expirationTime.Before(currentTime) {
		return "", utility.NewValidationError(fmt.Errorf("OTP has expired"), "expired", "otp")
	}

	encryptedData, err := utility.EncryptUser(int(user.ID), currentTime)

	if err != nil {
		return "", err
	}

	return encryptedData, nil

}
