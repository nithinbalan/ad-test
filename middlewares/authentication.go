package middlewares

import (
	"jethur/utility"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utility.SendUnauthorizedErrResp(c, "authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utility.SendUnauthorizedErrResp(c, "invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := utility.VerifyToken(tokenString)
		if err != nil {
			utility.SendUnauthorizedErrResp(c, "invalid token")
			c.Abort()
			return
		}
		c.Set("currentUserID", claims.UserID)

		c.Next()
	}
}
