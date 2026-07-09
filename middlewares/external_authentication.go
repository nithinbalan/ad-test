package middlewares

import (
	"strings"

	"jethur/utility"

	"github.com/gin-gonic/gin"
)

func ExternalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if err := utility.ValidateExternalClientIP(c.ClientIP()); err != nil {
			utility.SendForbiddenErrResp(c, err.Error())
			c.Abort()
			return
		}

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utility.SendUnauthorizedErrResp(c, "authorization header is required")
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			utility.SendUnauthorizedErrResp(c, "invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utility.VerifyExternalToken(tokenString)
		if err != nil {
			utility.SendUnauthorizedErrResp(c, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("clientID", claims.ClientID)

		c.Next()
	}
}
