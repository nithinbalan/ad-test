package services

import (
	"fmt"
	"jethur/config"
	authDTOs "jethur/grc/auth/dtos"
	"jethur/grc/auth/repositories"
	"jethur/utility"
	"strconv"
	"strings"
	"time"
)

type AuthServiceInterface interface {
	Authenticate(email string, password string) (map[string]string, error)
	SendPasswordRecoveryOTP(email string) error
	ReSendPasswordRecoveryOTP(email string) error
	VerifyPasswordRecoveryOTP(email, otp string) (string, error)
	SetPassword(setPasswordFormDTO authDTOs.PasswordRecoverySetPasswordFormDTO) error
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
		return "", utility.NewValidationError(fmt.Errorf("OTP has expired"), "otp_expired", "otp")
	}

	encryptedData, err := utility.EncryptUser(uint(user.ID), currentTime)

	if err != nil {
		return "", err
	}

	return encryptedData, nil
}

func (as AuthService) SetPassword(setPasswordFormDTO authDTOs.PasswordRecoverySetPasswordFormDTO) error {

	if setPasswordFormDTO.Password != setPasswordFormDTO.ConfirmPassword {
		return utility.NewValidationError(fmt.Errorf("password and confirm password do not match"), "password_match", "password")
	}
	if !isValidPassword(setPasswordFormDTO.Password) {
		return utility.NewValidationError(fmt.Errorf("password does not meet the requirements"), "password_complexity", "password")
	}

	userId, expirationTime, err := utility.DecryptUser(setPasswordFormDTO.Token)
	if err != nil {
		return err
	}
	expirationTime = expirationTime.Add(30 * time.Minute)

	currentTime := time.Now().UTC()

	if currentTime.After(expirationTime) {
		return utility.NewValidationError(fmt.Errorf("session has expired"), "session_expired", "session")
	}

	// Call repository to set password
	user, err := as.repo.SetPassword(userId, setPasswordFormDTO.Password)
	if err != nil {
		return err
	}

	// Prepare data for the email template
	emailData := authDTOs.PasswordRecoverySetPasswordEmailDTO{
		FirstNameEn: user.FirstNameEN,
		LastNameEn:  user.LastNameEN,
		Email:       user.Email,
	}

	to := []string{user.Email}
	subject := "Security alert"
	templateRelativePath := "auth/templates/email/reset_password.html"

	// Send email with the template
	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, nil)
	if err != nil {
		return err
	}

	return nil
}

func isValidPassword(password string) bool {
	// Password should contain at least 1 uppercase letter, 1 lowercase letter, 1 number, and 1 special character
	var (
		hasUpperCase = false
		hasLowerCase = false
		hasNumber    = false
		hasSpecial   = false
	)

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpperCase = true
		case 'a' <= char && char <= 'z':
			hasLowerCase = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsAny(string(char), "!@#$%^&*()-_+=~`[]{}\\|;:'\",.<>?"):
			hasSpecial = true
		}
	}

	return hasUpperCase && hasLowerCase && hasNumber && hasSpecial
}
