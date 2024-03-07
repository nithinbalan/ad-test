package utility

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendInternalErrResp(db *gorm.DB, c *gin.Context, errMsg string) {
	localizedErrorMsg := LocalizeMessage(c, "internal_server_error", "Internal server error!")

	go LogException(db, c, "internal_error", errMsg)

	resp := ErrorResponse{Error: localizedErrorMsg}
	c.JSON(http.StatusInternalServerError, resp)
}

func SendInvalidErrResp(db *gorm.DB, c *gin.Context, errs ...error) {

	var err error
	if len(errs) > 0 {
		err = errs[0] // Use the first error if any are provided
	}

	localizedErrorMsg := LocalizeMessage(c, "invalid_request", "Invalid Request!")
	errorMessages := ExtractValidationErrors(c, err)
	c.JSON(http.StatusBadRequest, gin.H{
		"error":  localizedErrorMsg,
		"errors": errorMessages,
	})

}

func SendNotFoundErrResp(db *gorm.DB, c *gin.Context) {
	resp := ErrorResponse{Error: LocalizeMessage(c, "resource_not_found", "Resource not found!")}
	c.JSON(http.StatusNotFound, resp)
}

func SendUnauthorizedErrResp(db *gorm.DB, c *gin.Context, errMsg string) {
	resp := ErrorResponse{Error: LocalizeMessage(c, "login_failed", "Login failed")}
	c.JSON(http.StatusUnauthorized, resp)
}
