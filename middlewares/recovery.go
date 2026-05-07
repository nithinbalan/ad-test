package middlewares

import (
	"fmt"
	"jethur/utility"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RecoveryWithDatabaseLog(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := string(debug.Stack())
				rid := CurrentRequestID(c)
				errorMessage := fmt.Sprintf("Recovered from panic: %v\nrequest_id=%s\nStack trace: %s", err, rid, stack)

				utility.LogError(errorMessage, map[string]interface{}{
					"error":      err,
					"stack":      stack,
					"request_id": rid,
				})

				utility.SendInternalErrResp(db, c, fmt.Sprintf("%v\n%s", err, stack))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
