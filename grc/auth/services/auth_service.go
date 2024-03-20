package services

import (
	"crypto/tls"
	"fmt"
	"jethur/config"
	"jethur/grc/auth/repositories"
	"jethur/grc/organization/models"

	"gopkg.in/gomail.v2"
)

type AuthServiceInterface interface {
	Authenticate(email string, password string) (map[string]string, error)
	GetOtpForgotPasswordRecovery(email string) (*models.User, error)
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

func (as AuthService) GetOtpForgotPasswordRecovery(email string) (*models.User, error) {
	user, err := as.repo.GetOtpForgotPasswordRecovery(email)

	if err != nil {
		return nil, err
	}

	config.LoadEnv() // Load environment variables

	from := config.GetEnv("SMTP_EMAIL_FROM", "no-reply@jethur.com")
	username := config.GetEnv("SMTP_EMAIL_USERNAME", "username")
	password := config.GetEnv("SMTP_EMAIL_PASSWORD", "password")
	host := config.GetEnv("SMTP_EMAIL_HOST", "smtp.gmail.com")

	// Prepare the email body
	emailBody := fmt.Sprintf(`
			Dear %s %s , <br/>
			<br/>
			To reset your password, please use the following One-Time Password (OTP): <br/>
			<br/>
			OTP: %s <br/>
			<br/>
			If you didn't initiate this password recovery, please ignore this message. <br/>
			<br/>
			Best regards,`, user.FirstNameEN, user.LastNameEN, user.Otp)

	// Create a goroutine to send the email
	go func() {
		// Create the email message
		msg := gomail.NewMessage()
		msg.SetHeader("From", from)
		msg.SetHeader("To", user.Email)
		msg.SetHeader("Subject", "Reset password instructions for Jethur")
		msg.SetBody("text/html", emailBody)

		// Create a dialer and configure it
		n := gomail.NewDialer(host, 587, username, password)
		n.TLSConfig = &tls.Config{InsecureSkipVerify: true}

		// Send the email
		if err := n.DialAndSend(msg); err != nil {
			fmt.Println(err) // Log the error
		}
	}()

	return user, nil

}
