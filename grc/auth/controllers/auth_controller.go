package controllers

import (
	"jethur/grc/auth/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type AuthController struct {
	service services.AuthServiceInterface
	db      *gorm.DB
}

func NewAuthController(service services.AuthServiceInterface) *AuthController {
	return &AuthController{
		service: service,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginReq LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	tokens, err := ac.service.Authenticate(loginReq.Email, loginReq.Password)
	if err != nil {
		if _, ok := err.(*utility.InternalServerError); ok {
			utility.SendInternalErrResp(ac.db, c, err.Error())
		} else {
			utility.SendUnauthorizedErrResp(c, "Authentication failed")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens["access_token"],
		"refresh_token": tokens["refresh_token"],
	})
}

func (ac *AuthController) GetOtpForgotPasswordRecovery(c *gin.Context) {
	var forgotPasswordReq ForgotPasswordRequest

	if err := c.ShouldBindJSON(&forgotPasswordReq); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Retrieve email from the request body
	email := forgotPasswordReq.Email

	// Get OTP for password recovery
	_, err := ac.service.GetOtpForgotPasswordRecovery(email)
	if err != nil {
		if _, ok := err.(*utility.InternalServerError); ok {
			utility.SendInternalErrResp(ac.db, c, err.Error())
		} else {
			utility.SendUnauthorizedErrResp(c, "Authentication failed")
		}
		return
	}

	// Return response with status OK
	c.JSON(http.StatusOK, gin.H{"message": "Password recovery otp sent to your email successfully"})
}
