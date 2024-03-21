package services

import (
	"jethur/grc/auth/repositories"
)

type AuthServiceInterface interface {
	Authenticate(email string, password string) (map[string]string, error)
	GetOtpForgotPasswordRecovery(email string) error
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

func (as AuthService) GetOtpForgotPasswordRecovery(email string) error {
	_, err := as.repo.GetOtpForgotPasswordRecovery(email)

	if err != nil {
		return err
	}

	//@todo Mail otp to respecticve user

	return nil

}
