package utility

import (
	models "jethur/grc/common/models"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LogException(db *gorm.DB, c *gin.Context, level, errorMessage string) {
	// Check if the context is nil
	if c == nil {
		log.Printf("Context is nil, unable to log error: %s", errorMessage)
		return
	}

	httpRequestMethod := c.Request.Method
	URL := c.Request.URL.String()
	stackTrace := string(debug.Stack())

	errorLogEntry := models.ErrorLog{
		HTTPStatusCode:    http.StatusInternalServerError,
		HTTPRequestMethod: httpRequestMethod,
		URL:               URL,
		ErrorMessage:      errorMessage,
		StackTrace:        stackTrace,
		UserID:            nil, // @todo
		CreatedAt:         time.Now(),
	}

	// Check if the db object is nil before trying to use it
	if db != nil {
		err := db.Create(&errorLogEntry).Error
		if err != nil {
			// If the database is not available, fall back to standard logging
			log.Printf("Unable to log to database: %v, original message: %s", err, errorMessage)
		}
	} else {
		// Fallback to standard logging if db is nil
		log.Printf("DB object is nil, fallback to standard logging: %s", errorMessage)
	}
}
