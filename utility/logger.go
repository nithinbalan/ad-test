package utility

import (
	models "jethur/grc/common/models"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
)

func LogException(c *gin.Context, level, errorMessage string) {

	db, err := GetDBFromContext(c)
	if err != nil {
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

	err = db.Create(&errorLogEntry).Error
	if err != nil {
		// If the database is not available, fall back to standard logging
		log.Printf("Unable to log to database: %v, original message: %s", err, errorMessage)
	}

}
