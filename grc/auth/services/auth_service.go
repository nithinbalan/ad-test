package services

import (
	"errors"
	"fmt"
	authDTOs "jethur/grc/auth/dtos"
	"jethur/grc/models"
	"jethur/grc/organization/services"
	"jethur/utility"

	"golang.org/x/crypto/bcrypt"

	"strings"
	"time"
)

type AuthServiceInterface interface {
	Authenticate(email string, password string) (*authDTOs.LoginResponseDTO, error)
	GetTokenFromRefreshToken(refreshToken string) (*authDTOs.LoginResponseDTO, error)
	LoginResendOTP(email string) error
	LoginVerifyOTP(email string, otp string) (*authDTOs.LoginResponseDTO, error)
	SendPasswordRecoveryOTP(email string) error
	ReSendPasswordRecoveryOTP(email string) error
	VerifyPasswordRecoveryOTP(email, otp string) (string, error)
	SetPassword(setPasswordFormDTO authDTOs.PasswordRecoverySetPasswordFormDTO) error
}

// Service and implementations
type AuthService struct {
	userService services.UserServiceInterface
}

func NewAuthService(userService services.UserServiceInterface) AuthServiceInterface {
	return &AuthService{
		userService: userService,
	}
}

func (as *AuthService) Authenticate(email string, password string) (*authDTOs.LoginResponseDTO, error) {
	// Get user
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Get userSecurity
	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(userSecurity.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	responseDTO := &authDTOs.LoginResponseDTO{}

	if userSecurity.IsMFA {
		subject := "Verify Your Login: Your One-Time Password"
		templateRelativePath := "auth/templates/email/send_mfa_login_otp.html"
		as.sendOTP(userSecurity, user, subject, templateRelativePath)
		responseDTO.IsMFAEnabled = true
	} else {
		accessToken, err := utility.GenerateToken(user.ID)
		if err != nil {
			return nil, utility.NewInternalServerError(err)
		}
		responseDTO.AccessToken = &accessToken

		refreshToken, err := utility.GenerateRefreshToken(user.ID)
		if err != nil {
			return nil, utility.NewInternalServerError(err)
		}
		responseDTO.RefreshToken = &refreshToken

		responseDTO.IsMFAEnabled = false

	}
	return responseDTO, nil
}

func (as *AuthService) GetTokenFromRefreshToken(refreshToken string) (*authDTOs.LoginResponseDTO, error) {
	//generate access token
	userClaims, err := utility.VerifyToken(refreshToken)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	//generate access token
	accessToken, err := utility.GenerateToken(userClaims.UserID)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	newRefreshToken, err := utility.GenerateRefreshToken(userClaims.UserID)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	responseDTO := &authDTOs.LoginResponseDTO{
		AccessToken:  &accessToken,
		RefreshToken: &newRefreshToken,
		IsMFAEnabled: false,
	}

	return responseDTO, nil
}

func (as AuthService) LoginResendOTP(email string) error {
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return err
	}

	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return err
	}

	resendTime := userSecurity.OTPExpiryAt.Add(-9 * time.Minute)
	currentTime := time.Now()
	if currentTime.Before(resendTime) {
		timeRemaining := resendTime.Sub(currentTime).Round(time.Minute)
		errorMessage := fmt.Sprintf("Resend request cannot be processed yet. Please wait %v more before trying again.", timeRemaining)
		return utility.NewValidationError(errors.New(errorMessage), "resend_not_yet_allowed", "resend_time")
	}

	subject := "Verify Your Login: Your One-Time Password"
	templateRelativePath := "auth/templates/email/send_mfa_login_otp.html"

	return as.sendOTP(userSecurity, user, subject, templateRelativePath)
}

func (as *AuthService) LoginVerifyOTP(email, otp string) (*authDTOs.LoginResponseDTO, error) {
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	if *userSecurity.OTP != otp {
		return nil, utility.NewValidationError(fmt.Errorf("invalid OTP"), "invalid_otp", "otp")
	}

	expirationTime := userSecurity.OTPExpiryAt
	currentTime := time.Now()

	if currentTime.After(*expirationTime) {
		return nil, utility.NewValidationError(fmt.Errorf("OTP has expired"), "otp_expired", "otp")
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

	responseDTO := &authDTOs.LoginResponseDTO{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
		IsMFAEnabled: true,
	}

	return responseDTO, nil
}

func (as AuthService) SendPasswordRecoveryOTP(email string) error {
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return err
	}

	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return err
	}

	subject := "Your Password Reset Verification Code"
	templateRelativePath := "auth/templates/email/send_password_recovery_otp.html"

	return as.sendOTP(userSecurity, user, subject, templateRelativePath)
}

func (as AuthService) ReSendPasswordRecoveryOTP(email string) error {
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return err
	}

	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return err
	}

	resendTime := userSecurity.OTPExpiryAt.Add(-9 * time.Minute)
	currentTime := time.Now()
	if currentTime.Before(resendTime) {
		timeRemaining := resendTime.Sub(currentTime).Round(time.Minute)
		errorMessage := fmt.Sprintf("Resend request cannot be processed yet. Please wait %v more before trying again.", timeRemaining)
		return utility.NewValidationError(errors.New(errorMessage), "resend_not_yet_allowed", "resend_time")
	}

	subject := "Your Password Reset Verification Code"
	templateRelativePath := "auth/templates/email/send_password_recovery_otp.html"

	return as.sendOTP(userSecurity, user, subject, templateRelativePath)
}

func (as *AuthService) VerifyPasswordRecoveryOTP(email, otp string) (string, error) {
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return "", err
	}

	userSecurity, err := as.userService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return "", err
	}

	if *userSecurity.OTP != otp {
		return "", utility.NewValidationError(fmt.Errorf("invalid OTP"), "invalid_otp", "otp")
	}

	expirationTime := userSecurity.OTPExpiryAt
	currentTime := time.Now()

	if currentTime.After(*expirationTime) {
		return "", utility.NewValidationError(fmt.Errorf("OTP has expired"), "otp_expired", "otp")
	}

	plaintext := fmt.Sprintf("%d:%d", user.ID, currentTime.Unix())
	plaintextBytes := []byte(plaintext)
	encryptedData, err := utility.EncryptData(plaintextBytes)
	if err != nil {
		return "", err
	}

	return encryptedData, nil
}

func (as AuthService) SetPassword(setPasswordFormDTO authDTOs.PasswordRecoverySetPasswordFormDTO) error {
	if !isValidPassword(setPasswordFormDTO.Password) {
		return utility.NewValidationError(fmt.Errorf("password does not meet the requirements"), "password_complexity", "password")
	}

	decryptedBytes, err := utility.DecryptData(setPasswordFormDTO.Token)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("failed to decrypt data: %w", err), "decryption_failed", "token")
	}

	decryptedText := string(decryptedBytes)

	var userId uint
	var timestamp int64

	_, err = fmt.Sscanf(decryptedText, "%d:%d", &userId, &timestamp)
	if err != nil {
		return err
	}

	timeFromTimestamp := time.Unix(timestamp, 0)

	currentTime := time.Now()

	fifteenMinutesAgo := currentTime.Add(-15 * time.Minute)

	if timeFromTimestamp.Before(fifteenMinutesAgo) {
		return utility.NewValidationError(fmt.Errorf("session has expired"), "session_expired", "session")
	}

	userSecurity, err := as.userService.GetUserSecurityByUserID(userId)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(setPasswordFormDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userSecurity.PasswordHash = string(hashedPassword)

	err = as.userService.UpdateUserSecurity(userSecurity)

	if err != nil {
		return err
	}

	user, err := as.userService.GetByID(userId)
	if err != nil {
		return err
	}

	emailData := authDTOs.PasswordRecoverySetPasswordEmailDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Email:       user.Email,
	}

	to := []string{user.Email}
	subject := "Important: Password Update Confirmation for Your Jethur Account"
	templateRelativePath := "auth/templates/email/password_change_confirmation.html"

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

// Private method to handle the common functionality of sending OTP.
func (as AuthService) sendOTP(userSecurity *models.UserSecurity, user *models.User, mailSubject string, mailRelativePath string) error {
	otpLength := 6
	verifyOTPDuration := 10

	otp := utility.GenerateOTP(otpLength)
	currentTime := time.Now()

	expirationTime := currentTime.Add(time.Duration(verifyOTPDuration) * time.Minute)

	userSecurity.OTP = &otp
	userSecurity.OTPExpiryAt = &expirationTime

	err := as.userService.UpdateUserSecurity(userSecurity)
	if err != nil {
		return err
	}

	emailData := authDTOs.PasswordRecoveryDataDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Otp:         otp,
	}

	to := []string{user.Email}
	subject := mailSubject
	templateRelativePath := mailRelativePath

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, nil)
	if err != nil {
		return err
	}

	return nil
}
