package services

import (
	"crypto/tls"
	"errors"
	"fmt"
	authDTOs "jethur/grc/auth/dtos"
	commonDTOs "jethur/grc/common/dtos"
	commonService "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/grc/organization/services"
	settingService "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/go-ldap/ldap/v3"
	"golang.org/x/crypto/bcrypt"

	"strings"
	"time"
)

type AuthServiceInterface interface {
	Authenticate(email string, password string, ipAddress string, location string, device string, browser string) (*authDTOs.LoginResponseDTO, error)
	GetTokenFromRefreshToken(refreshToken string) (*authDTOs.LoginResponseDTO, error)
	LoginResendOTP(email string) error
	LoginVerifyOTP(email string, otp string, ipAddress string, location string, device string, browser string) (*authDTOs.LoginResponseDTO, error)
	SendPasswordRecoveryOTP(email string) error
	ReSendPasswordRecoveryOTP(email string) error
	VerifyPasswordRecoveryOTP(email, otp string) (string, error)
	SetPassword(setPasswordFormDTO authDTOs.PasswordRecoverySetPasswordFormDTO) error
}

// Service and implementations
type AuthService struct {
	userService               services.UserServiceInterface
	userSettingService        services.UserSettingServiceInterface
	notificationSvc           commonService.NotificationServiceInterface
	emailTemplateSettingsSvc  settingService.EmailTemplateSettingsServiceInterface
	activeDirectorySettingSvc settingService.ActiveDirectorySettingServiceInterface
	integrationSvc            settingService.IntegrationServiceInterface
}

func NewAuthService(
	userService services.UserServiceInterface,
	userSettingService services.UserSettingServiceInterface,
	notificationSvc commonService.NotificationServiceInterface,
	emailTemplateSettingsSvc settingService.EmailTemplateSettingsServiceInterface,
	activeDirectorySettingSvc settingService.ActiveDirectorySettingServiceInterface,
	integrationSvc settingService.IntegrationServiceInterface) AuthServiceInterface {
	return &AuthService{
		userService:               userService,
		userSettingService:        userSettingService,
		notificationSvc:           notificationSvc,
		emailTemplateSettingsSvc:  emailTemplateSettingsSvc,
		activeDirectorySettingSvc: activeDirectorySettingSvc,
		integrationSvc:            integrationSvc,
	}
}

func (as *AuthService) Authenticate(email string, password string, ipAddress string, location string, device string, browser string) (*authDTOs.LoginResponseDTO, error) {
	// ✅ Check if AD login is enabled/configured
	integration, err := as.integrationSvc.GetByType("active-directory")
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}
	// ✅ Check if integration is enabled
	if integration.IsEnable {
		adSetting, err := as.activeDirectorySettingSvc.GetFirst()
		if err != nil {
			return nil, utility.NewInternalServerError(err)
		}
		if adSetting != nil {
			if response, err := as.LoginWithActiveDirectory(email, password, ipAddress, location, device, browser); err != nil {
				return nil, err
			} else if response != nil {
				return response, nil
			}
		}
	}

	// Get user
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Get userSecurity
	userSecurity, err := as.userSettingService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(userSecurity.PasswordHash), []byte(password))
	if err != nil {
		loginStatus := "failed"
		logErr := as.userService.CreateUserLoginHistory(user, loginStatus, ipAddress, location, device, browser)
		if logErr != nil {
			return nil, utility.NewInternalServerError(logErr)
		}
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

		loginStatus := "success"
		err = as.userService.CreateUserLoginHistory(user, loginStatus, ipAddress, location, device, browser)
		if err != nil {
			return nil, utility.NewInternalServerError(err)
		}

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

	userSecurity, err := as.userSettingService.GetUserSecurityByUserID(user.ID)
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

func (as *AuthService) LoginVerifyOTP(email, otp string, ipAddress string, location string, device string, browser string) (*authDTOs.LoginResponseDTO, error) {
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	userSecurity, err := as.userSettingService.GetUserSecurityByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	if *userSecurity.OTP != otp {
		loginStatus := "failed"
		logErr := as.userService.CreateUserLoginHistory(user, loginStatus, ipAddress, location, device, browser)
		if logErr != nil {
			return nil, utility.NewInternalServerError(logErr)
		}
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

	loginStatus := "success"
	err = as.userService.CreateUserLoginHistory(user, loginStatus, ipAddress, location, device, browser)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	return responseDTO, nil
}

func (as AuthService) SendPasswordRecoveryOTP(email string) error {
	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return err
	}

	userSecurity, err := as.userSettingService.GetUserSecurityByUserID(user.ID)
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

	userSecurity, err := as.userSettingService.GetUserSecurityByUserID(user.ID)
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

	userSecurity, err := as.userSettingService.GetUserSecurityByUserID(user.ID)
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

	userSecurity, err := as.userSettingService.GetUserSecurityByUserID(userId)
	if err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(setPasswordFormDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userSecurity.PasswordHash = string(hashedPassword)

	err = as.userSettingService.UpdateUserSecurity(userSecurity)

	if err != nil {
		return err
	}

	user, err := as.userService.GetByID(userId)
	if err != nil {
		return err
	}

	emailTemplateSettings, err := as.emailTemplateSettingsSvc.GetEmailTemplateSettings()
	if err != nil {
		return err
	}

	socialLinks, err := as.notificationSvc.GetSocialMediaLinks()
	if err != nil {
		return err
	}

	emailTemplateAndSocialDTO := &commonDTOs.EmailTemplateAndSocialDTO{
		HeaderImageURL: emailTemplateSettings.HeaderImageURL,
		SupportText:    emailTemplateSettings.SupportText,
		ClosingText:    emailTemplateSettings.ClosingText,
		FooterLogoURL:  emailTemplateSettings.FooterLogoURL,
		FooterAddress:  emailTemplateSettings.FooterAddress,
		FooterEmail:    emailTemplateSettings.FooterEmail,
		Linkedin:       &socialLinks.Linkedin,
		X:              &socialLinks.X,
		Facebook:       &socialLinks.Facebook,
	}

	emailData := authDTOs.PasswordRecoverySetPasswordEmailDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Email:       user.Email,
	}

	to := []string{user.Email}
	subject := "Important: Password Update Confirmation for Your Jethur Account"
	templateRelativePath := "auth/templates/email/password_change_confirmation.html"

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, emailTemplateAndSocialDTO, nil)
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

	err := as.userSettingService.UpdateUserSecurity(userSecurity)
	if err != nil {
		return err
	}

	emailTemplateSettings, err := as.emailTemplateSettingsSvc.GetEmailTemplateSettings()
	if err != nil {
		return err
	}

	socialLinks, err := as.notificationSvc.GetSocialMediaLinks()
	if err != nil {
		return err
	}

	emailTemplateAndSocialDTO := &commonDTOs.EmailTemplateAndSocialDTO{
		HeaderImageURL: emailTemplateSettings.HeaderImageURL,
		SupportText:    emailTemplateSettings.SupportText,
		ClosingText:    emailTemplateSettings.ClosingText,
		FooterLogoURL:  emailTemplateSettings.FooterLogoURL,
		FooterAddress:  emailTemplateSettings.FooterAddress,
		FooterEmail:    emailTemplateSettings.FooterEmail,
		Linkedin:       &socialLinks.Linkedin,
		X:              &socialLinks.X,
		Facebook:       &socialLinks.Facebook,
	}

	emailData := authDTOs.PasswordRecoveryDataDTO{
		FirstNameEn: user.FirstNameEn,
		LastNameEn:  user.LastNameEn,
		Otp:         otp,
	}

	to := []string{user.Email}
	subject := mailSubject
	templateRelativePath := mailRelativePath

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, emailTemplateAndSocialDTO, nil)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthService) LoginWithActiveDirectory(email, password, ip, location, device, browser string) (*authDTOs.LoginResponseDTO, error) {
	adSetting, err := as.activeDirectorySettingSvc.GetFirst()
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}
	if adSetting == nil {
		return nil, nil // no AD config, skip
	}

	bindPsw := []byte(adSetting.BindPassword)
	decryptedBindPwd, err := utility.DecryptData(string(bindPsw))
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	address := fmt.Sprintf("%s:%d", adSetting.Server, adSetting.Port)
	var l *ldap.Conn
	if adSetting.UseSSL {
		l, err = ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true})
	} else {
		l, err = ldap.Dial("tcp", address)
	}
	if err != nil {
		return nil, fmt.Errorf("LDAP connection failed: %w", err)
	}
	defer l.Close()

	if err := l.Bind(adSetting.BindDN, string(decryptedBindPwd)); err != nil {
		return nil, fmt.Errorf("Bind failed: %w", err)
	}

	searchReq := ldap.NewSearchRequest(
		adSetting.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(mail=%s)", email),
		[]string{"dn", "cn", "displayName", "mail"},
		nil,
	)

	sr, err := l.Search(searchReq)
	if err != nil || len(sr.Entries) == 0 {
		return nil, fmt.Errorf("AD user not found: %w", err)
	}

	userDN := sr.Entries[0].DN

	if err := l.Bind(userDN, password); err != nil {
		// _ = as.userService.CreateUserLoginHistory(nil, "failed", ip, location, device, browser)
		return nil, fmt.Errorf("AD login failed: %w", err)
	}

	user, err := as.userService.GetByEmail(email)
	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, email)
	}

	accessToken, err := utility.GenerateToken(user.ID)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}
	refreshToken, err := utility.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	_ = as.userService.CreateUserLoginHistory(user, "success", ip, location, device, browser)

	return &authDTOs.LoginResponseDTO{
		AccessToken:  &accessToken,
		RefreshToken: &refreshToken,
		IsMFAEnabled: false,
	}, nil
}
