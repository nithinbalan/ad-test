package services

import (
	"fmt"
	"jethur/config"
	authDTOs "jethur/grc/auth/dtos"
	"jethur/grc/auth/repositories"
	"jethur/grc/organization/services"
	"jethur/utility"

	"golang.org/x/crypto/bcrypt"

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
	repo        repositories.AuthRepositoryInterface
	userService services.UserServiceInterface
}

func NewAuthService(repo repositories.AuthRepositoryInterface, userService services.UserServiceInterface) AuthServiceInterface {
	return &AuthService{
		repo:        repo,
		userService: userService,
	}
}

func (as *AuthService) Authenticate(email string, password string) (map[string]string, error) {
	//get user
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	//get userSecurity
	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(userSecurity.PasswordHash), []byte(password))
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid username or password"), "invalid_credentials", "email")
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

func (as AuthService) SendPasswordRecoveryOTP(email string) error {

	//get user
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return err
	}

	//get userSecurity
	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return err
	}

	otpLength := 6
	otp := utility.GenerateOTP(otpLength)
	currentTime := time.Now().UTC()

	verifyOTPDurationStr := config.GetEnv("VERIFY_OTP_DURATION", "10")
	verifyOTPDuration, _ := strconv.Atoi(verifyOTPDurationStr)

	expirationTime := currentTime.Add(time.Duration(verifyOTPDuration) * time.Minute)

	userSecurity.OTP = &otp
	userSecurity.OTPExpiryAt = &expirationTime

	err = as.userService.UpdateUserSecurity(userSecurity)

	if err != nil {
		return err
	}

	emailData := authDTOs.PasswordRecoveryDataDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Otp:         otp,
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

func (as AuthService) ReSendPasswordRecoveryOTP(email string) error {
	//get user
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return err
	}

	//get userSecurity
	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return err
	}

	resendTime := userSecurity.OTPExpiryAt.Add(-9 * time.Minute)
	currentTime := time.Now().UTC()
	if currentTime.Before(resendTime) {
		return utility.NewValidationError(fmt.Errorf("resend time has not reached yet"), "resend_time_not_reached", "resend_time")
	}

	otpLength := 6
	otp := utility.GenerateOTP(otpLength)

	verifyOTPDurationStr := config.GetEnv("VERIFY_OTP_DURATION", "10")
	verifyOTPDuration, _ := strconv.Atoi(verifyOTPDurationStr)

	expirationTime := currentTime.Add(time.Duration(verifyOTPDuration) * time.Minute)

	userSecurity.OTP = &otp
	userSecurity.OTPExpiryAt = &expirationTime

	err = as.userService.UpdateUserSecurity(userSecurity)

	if err != nil {
		return err
	}

	emailData := authDTOs.PasswordRecoveryDataDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Otp:         otp,
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

	//get user
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return "", err
	}

	//get userSecurity
	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return "", err
	}

	if *userSecurity.OTP != otp {
		return "", utility.NewValidationError(fmt.Errorf("invalid OTP"), "invalid_otp", "otp")
	}

	expirationTime := userSecurity.OTPExpiryAt
	currentTime := time.Now().UTC()

	if currentTime.After(*expirationTime) {
		return "", utility.NewValidationError(fmt.Errorf("OTP has expired"), "otp_expired", "otp")
	}

	encryptedData, err := utility.EncryptUser(uint(user.ID), currentTime)

	if err != nil {
		return "", err
	}

	return encryptedData, nil
}

func (as AuthService) SetPassword(setPasswordFormDTO authDTOs.PasswordRecoverySetPasswordFormDTO) error {

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

	//get userSecurity
	userSecurity, err := as.userService.GetUserSecurityByUserID(userId)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(setPasswordFormDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return utility.NewInternalServerError(err)
	}
	// Set the hashed password
	userSecurity.PasswordHash = string(hashedPassword)

	err = as.userService.UpdateUserSecurity(userSecurity)

	if err != nil {
		return err
	}

	//get user
	user, err := as.userService.GetByID(userId)
	if err != nil {
		return err
	}

	// Prepare data for the email template
	emailData := authDTOs.PasswordRecoverySetPasswordEmailDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
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

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}
