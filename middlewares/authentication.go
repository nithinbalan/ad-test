package middlewares

import (
	"jethur/utility"

	"github.com/gin-gonic/gin"
)

var jwtManager *utility.JWTManager

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// tokenString := c.GetHeader("Authorization")
		// log.Fatalln(tokenString)

		// tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDg4MzQ2MzAsImlhdCI6MTcwODc0ODIzMCwidXNlcl9pZCI6MX0.6dzUvEEa0eKcT8mmT8KYJlQykX01pqMliZ_vNAoiGQs"
		// Verify the token using your JWTManager
		// claims, err := jwtManager.VerifyToken(tokenString)
		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		// 	c.Abort()
		// 	return
		// }

		// If the token is valid, set the user ID in the Gin context for later use
		// c.Set("currentUserID", claims.UserID)

		// Continue down the chain to handler
		c.Next()
	}
}
