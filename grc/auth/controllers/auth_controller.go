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
	var loginFormDTO authDTOs.LoginFormDTO

	if err := c.ShouldBindJSON(&loginFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	responseDTO, err := ac.service.Authenticate(loginFormDTO.Email, loginFormDTO.Password)
	if err != nil {
		switch e := err.(type) {
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		case *utility.InternalServerError:
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		default:
			utility.SendUnauthorizedErrResp(c, "Authentication failed")
			return
		}
	}

	c.JSON(http.StatusOK, responseDTO)
}

func (ac *AuthController) GetTokenFromRefreshToken(c *gin.Context) {
	var refreshTokenFormDTO authDTOs.RefreshTokenFormDTO

	if err := c.ShouldBindJSON(&refreshTokenFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	responseDTO, err := ac.service.GetTokenFromRefreshToken(refreshTokenFormDTO.RefreshToken)
	if err != nil {
		switch e := err.(type) {
		case *utility.InternalServerError:
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		default:
			utility.SendUnauthorizedErrResp(c, "Authentication failed")
			return
		}
	}

	c.JSON(http.StatusOK, responseDTO)

}

func (ac *AuthController) LoginResendOTP(c *gin.Context) {
	var sendOTPFormDTO authDTOs.SendOTPFormDTO

	if err := c.ShouldBindJSON(&sendOTPFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	email := sendOTPFormDTO.Email

	err := ac.service.LoginResendOTP(email)
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

func (ac *AuthController) LoginVerifyOTP(c *gin.Context) {
	var verifyOTPFormDTO authDTOs.VerifyOTPFormDTO

	if err := c.ShouldBindJSON(&verifyOTPFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	email := verifyOTPFormDTO.Email

	otp := verifyOTPFormDTO.Otp

	responseDTO, err := ac.service.LoginVerifyOTP(email, otp)
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

	c.JSON(http.StatusOK, responseDTO)
}

func (ac *AuthController) SendPasswordRecoveryOTP(c *gin.Context) {
	var passwordRecoveryFormDTO authDTOs.SendOTPFormDTO

	if err := c.ShouldBindJSON(&passwordRecoveryFormDTO); err != nil {
		utility.SendInvalidErrResp(c, err)
		return
	}

	email := passwordRecoveryFormDTO.Email

	err := ac.service.SendPasswordRecoveryOTP(email)
	if err != nil {
		switch e := err.(type) {
		case *utility.ResourceNotFoundError:
			utility.SendNotFoundErrResp(c)
			return
		case *utility.ValidationError:
			utility.SendInvalidErrResp(c, e)
			return
		default:
			fmt.Println(err)
			utility.SendInternalErrResp(ac.db, c, e.Error())
			return
		}
	}

	c.Status(http.StatusOK)
}

func (ac *AuthController) ReSendPasswordRecoveryOTP(c *gin.Context) {
	var passwordRecoveryFormDTO authDTOs.SendOTPFormDTO

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
	var passwordRecoveryOTPVerifyFormDTO authDTOs.VerifyOTPFormDTO

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
