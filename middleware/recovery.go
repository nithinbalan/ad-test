package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RecoveryWithDatabaseLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Capture the stack trace
				stack := string(debug.Stack())

				// Attempt to get the database connection from the Gin context
				db, exists := c.MustGet("db").(*gorm.DB)
				if !exists {
					log.Println("Database connection not found in context")
					// Handle the absence of a DB connection, or continue without logging to the database
				} else {
					// Log to database using the extracted DB connection
					go utility.LogException(db, c.Request, "PANIC", fmt.Sprintf("%v\n%s", err, stack))
				}

				// Use utility function to get the localized message
				localizedError := utility.LocalizeMessage(c, "internal_server_error", "Internal server error!")

				// Respond with a 500 Internal Server Error status code and localized message
				utility.SendErrorResponse(c, http.StatusInternalServerError, localizedError)
			}
		}()
		c.Next()
	}
}
