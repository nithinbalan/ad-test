package middlewares

import (
	"jethur/utility"
	"net/http"

	"github.com/gin-gonic/gin"
)

var jwtManager *utility.JWTManager

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		// Verify the token using your JWTManager
		claims, err := jwtManager.VerifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// If the token is valid, set the user ID in the Gin context for later use
		c.Set("currentUserID", claims.UserID)

		// Continue down the chain to handler
		c.Next()
	}
}
