package utility

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendErrorResponse(db *gorm.DB, c *gin.Context, statusCode int, errMsg string, errs ...error) {

	var err error
	if len(errs) > 0 {
		err = errs[0] // Use the first error if any are provided
		log.Fatalf(fmt.Sprintf("%+v", err))
	}

	var errorType string = ""
	var localizedErrorMsg string = ""

	if statusCode == http.StatusInternalServerError {
		errorType = "500"
		localizedErrorMsg = LocalizeMessage(c, "internal_server_error", "Internal server error!")
		//log 500 error on database
		go LogException(db, c, errorType, errMsg)
	} else if statusCode == http.StatusBadRequest {
		localizedErrorMsg = LocalizeMessage(c, "invalid_request", "Invalid Request!")
		errorMessages := ExtractValidationErrors(c, err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  localizedErrorMsg, // Single error message
			"errors": errorMessages,     // Array of error messages
		})
		return
	} else if statusCode == http.StatusUnauthorized {
		localizedErrorMsg = LocalizeMessage(c, "login_failed", "Login failed")
	}

	resp := ErrorResponse{
		Error: localizedErrorMsg,
	}
	c.JSON(statusCode, resp)
}
