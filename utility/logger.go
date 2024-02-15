package utility

import (
	models "jethur/grc/common/models"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

func LogException(db *gorm.DB, req *http.Request, level, message string) {

	url := "Unknown"
	if req != nil {
		url = req.URL.String()
	}

	errorLogEntry := models.ErrorLog{
		Level:     level,
		URL:       url,
		Message:   message,
		CreatedAt: time.Now(),
	}

	err := db.Create(&errorLogEntry).Error
	if err != nil {
		// If the database is not available, fall back to standard logging
		log.Printf("Unable to log to database: %v, original message: %s", err, message)
	}

	log.Printf(message)
}
