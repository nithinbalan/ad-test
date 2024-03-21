package controllers

import (
	authDTOs "jethur/grc/auth/dtos"
	"jethur/grc/auth/services"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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
	var authDTO authDTOs.AuthFormDTO

	if err := c.ShouldBindJSON(&authDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	tokens, err := ac.service.Authenticate(authDTO.Email, authDTO.Password)
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
	var forgotPasswordFormDTO authDTOs.ForgotPasswordFormDTO

	if err := c.ShouldBindJSON(&forgotPasswordFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	// Retrieve email from the request body
	email := forgotPasswordFormDTO.Email

	// Get OTP for password recovery
	err := ac.service.GetOtpForgotPasswordRecovery(email)
	if err != nil {
		if _, ok := err.(*utility.InternalServerError); ok {
			utility.SendInternalErrResp(ac.db, c, err.Error())
		}
		return
	}

	// Return response with status OK
	c.Status(http.StatusOK)
}
