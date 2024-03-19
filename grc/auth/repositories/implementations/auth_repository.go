package implementations

import (
	"errors"
	"fmt"
	"jethur/config"
	"jethur/grc/organization/models"
	"jethur/utility"

	"crypto/tls"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repo *AuthRepository) Authenticate(c *gin.Context, email, password string) (map[string]string, error) {

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

func (repo *AuthRepository) GetOtpForgotPasswordRecovery(c *gin.Context, email string) (map[string]string, error) {

	var user models.User

	result := repo.db.First(&user, "email = ?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("invalid email")
	} else if result.Error != nil {
		return nil, utility.NewInternalServerError(result.Error)
	}
	otpLength := 6
	otp := utility.GenerateOTP(otpLength)
	user.Otp = otp
	user.SentAt = time.Now()

	// Update the user model
	if err := repo.db.Save(&user).Error; err != nil {
		return nil, utility.NewInternalServerError(err)
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
			Best regards,`, user.FirstNameEN, user.LastNameEN, otp)

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

	otpMap := map[string]string{"otp": otp} // Create a map containing the OTP
	return otpMap, nil

}
