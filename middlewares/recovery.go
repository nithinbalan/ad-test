package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RecoveryWithDatabaseLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Capture the stack trace
				stack := string(debug.Stack())

				// Respond with a 500 Internal Server Error status code and localized message
				utility.SendErrorResponse(db, c, http.StatusInternalServerError, fmt.Sprintf("%v\n%s", err, stack))
			}
		}()
		c.Next()
	}
}
