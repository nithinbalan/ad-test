package middleware

import (
	"jethur/config"

	"github.com/gin-gonic/gin"
)

// DatabaseMiddleware sets the database instance into the Gin context
func DatabaseMiddleware() gin.HandlerFunc {
	db := config.GetDatabaseInstance() // Obtain the database instance

	return func(c *gin.Context) {
		c.Set("db", db) // Set the database instance in the context
		c.Next()        // Continue to the next middleware/handler
	}
}
