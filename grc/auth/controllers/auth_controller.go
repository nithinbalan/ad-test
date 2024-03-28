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

func NewAuthController(db *gorm.DB, service services.AuthServiceInterface) *AuthController {
	return &AuthController{
		db:      db,
		service: service,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var authDTO authDTOs.LoginFormDTO

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

func (ac *AuthController) SendPasswordRecoveryOTP(c *gin.Context) {
	var passwordRecoveryFormDTO authDTOs.PasswordRecoveryFormDTO

	if err := c.ShouldBindJSON(&passwordRecoveryFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	email := passwordRecoveryFormDTO.Email

	err := ac.service.SendPasswordRecoveryOTP(email)
	if err != nil {
		if _, ok := err.(*utility.InternalServerError); ok {
			utility.SendInternalErrResp(ac.db, c, err.Error())
		}
		return
	}

	c.Status(http.StatusOK)
}

func (ac *AuthController) ReSendPasswordRecoveryOTP(c *gin.Context) {
	var passwordRecoveryFormDTO authDTOs.PasswordRecoveryFormDTO

	if err := c.ShouldBindJSON(&passwordRecoveryFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	email := passwordRecoveryFormDTO.Email

	err := ac.service.ReSendPasswordRecoveryOTP(email)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ac *AuthController) VerifyPasswordRecoveryOTP(c *gin.Context) {
	var passwordRecoveryOTPVerifyFormDTO authDTOs.PasswordRecoveryOTPVerifyFormDTO

	if err := c.ShouldBindJSON(&passwordRecoveryOTPVerifyFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	email := passwordRecoveryOTPVerifyFormDTO.Email

	otp := passwordRecoveryOTPVerifyFormDTO.Otp

	encryptedData, err := ac.service.VerifyPasswordRecoveryOTP(email, otp)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"token": encryptedData})
}

func (ac *AuthController) SetPassword(c *gin.Context) {
	var setPasswordFormDTO authDTOs.PasswordRecoverySetPasswordFormDTO

	if err := c.ShouldBindJSON(&setPasswordFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	err := ac.service.SetPassword(setPasswordFormDTO)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}
