package middlewares

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserStatusMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		currentUserID, err := utility.GetCurrentUserID(c)
		if err != nil {
			utility.SendUnauthorizedErrResp(c, "Unauthorized")
			c.Abort()
			return
		}

		var user models.User

		err = db.
			Joins("JOIN user_statuses ON user_statuses.id = users.user_status_id").
			Where("users.id = ? AND user_statuses.is_active = ?", currentUserID, true).
			First(&user).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utility.SendForbiddenErrResp(c, "User account is inactive")
			} else {
				utility.SendInternalErrResp(db, c, "Unable to verify user status")
			}
			c.Abort()
			return
		}

		c.Next()
	}
}
