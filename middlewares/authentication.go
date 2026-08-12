package middlewares

import (
	"errors"
	"fmt"
	"jethur/grc/models"
	"jethur/utility"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
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
			fmt.Println(err)
			utility.SendUnauthorizedErrResp(c, "invalid token")
			c.Abort()
			return
		}

		currentUserID := claims.UserID
		c.Set("currentUserID", currentUserID)

		var user models.User
		err = db.Preload("UserStatus").
			Preload("UserSetting").
			Where("id = ?", currentUserID).
			First(&user).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utility.SendUnauthorizedErrResp(c, "User not found")
			} else {
				utility.SendInternalErrResp(db, c, "Unable to verify user status")
			}
			c.Abort()
			return
		}

		if user.UserStatus == nil || !user.UserStatus.IsActive {
			utility.SendForbiddenErrResp(c, "User account is inactive")
			c.Abort()
			return
		}

		// Override the request localizer with the user's preferred language
		if user.UserSetting != nil && user.UserSetting.Language != "" {
			if bundleInterface, exists := c.Get("i18n_bundle"); exists {
				if bundle, ok := bundleInterface.(*i18n.Bundle); ok {
					localizer := i18n.NewLocalizer(bundle, user.UserSetting.Language)
					c.Set("localizer", localizer)
				}
			}
		}

		c.Next()
	}
}
