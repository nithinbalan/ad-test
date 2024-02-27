package controllers

import (
	"jethur/grc/auth/repositories/interfaces"
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthController struct {
	repo interfaces.AuthRepositoryInterface
	db   *gorm.DB
}

func NewAuthController(repo interfaces.AuthRepositoryInterface, db *gorm.DB) *AuthController {
	return &AuthController{
		repo: repo,
		db:   db,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (ac *AuthController) Login(c *gin.Context) {
	var loginReq LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		utility.SendErrorResponse(ac.db, c, http.StatusBadRequest, "Invalid Request", err)
		return
	}

	tokens, err := ac.repo.Authenticate(c, loginReq.Email, loginReq.Password)
	if err != nil {
		if _, ok := err.(*utility.InternalServerError); ok {
			utility.SendErrorResponse(ac.db, c, http.StatusInternalServerError, err.Error())
		} else {
			utility.SendErrorResponse(ac.db, c, http.StatusUnauthorized, "Authentication failed")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  tokens["access_token"],
		"refresh_token": tokens["refresh_token"],
	})
}
