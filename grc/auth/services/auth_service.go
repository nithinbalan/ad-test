package services

import (
	authDTOs "jethur/grc/auth/dtos"
	"jethur/grc/auth/repositories"
	"jethur/utility"
)

type AuthServiceInterface interface {
	Authenticate(email string, password string) (map[string]string, error)
	SendPasswordRecoveryOTP(email string) error
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

	emailData := authDTOs.ForgotPasswordDataDTO{
		FirstNameEN: user.FirstNameEN,
		LastNameEN:  user.LastNameEN,
		Otp:         user.Otp,
	}

	to := []string{user.Email}
	subject := "Password Reset Verification Code"
	templateRelativePath := "auth/templates/email/forgot_password.html"
	// Attachments, if any
	attachments := []string{}

	err = utility.SendTemplatedEmail(to, subject, templateRelativePath, emailData, attachments)
	if err != nil {
		return err
	}

	return nil

}
