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
				errorMessage := fmt.Sprintf("Recovered from panic: %v\nStack trace: %s", err, stack)

				utility.LogError(errorMessage, map[string]interface{}{
					"error": err,
					"stack": stack,
				})

				utility.SendInternalErrResp(db, c, fmt.Sprintf("%v\n%s", err, stack))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
