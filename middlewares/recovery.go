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
				// log.Printf("Recovered from panic: %v\nStack trace: %s", err)

				utility.SendInternalErrResp(db, c, fmt.Sprintf("%v\n%s", err, stack))
				c.Abort()
				return
			}
		}()
		c.Next()
	}
}
