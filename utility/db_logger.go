package utility

import (
	"jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

func LogErrorToDB(errorSvc services.ErrorLogServiceInterface, c *gin.Context, errorMessage string) {
	errorLogEntry := models.ErrorLog{
		HTTPStatusCode:    http.StatusInternalServerError,
		HTTPRequestMethod: c.Request.Method,
		URL:               c.Request.URL.String(),
		ErrorMessage:      errorMessage,
		StackTrace:        string(debug.Stack()),
		CreatedAt:         time.Now(),
	}

	// Attempt to obtain the current user ID
	if currentUserId, err := GetCurrentUserID(c); err == nil {
		errorLogEntry.UserID = &currentUserId
	}

	if err := errorSvc.LogError(errorLogEntry); err != nil {
		LogError("Failed to log error to database", map[string]interface{}{
			"error":   err.Error(),
			"details": errorLogEntry,
		})
	}
}
