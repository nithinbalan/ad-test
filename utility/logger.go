package utility

import (
	"fmt"
	"jethur/grc/logs/services"

	"github.com/gin-gonic/gin"
)

func LogErrorToDB(errorSvc services.ErrorLogServiceInterface, c *gin.Context, errorMessage string) {
	currentUserId, err := GetCurrentUserID(c)
	if err != nil {
		// @todo write log to file
		fmt.Println(err.Error())
	}

	httpRequestMethod := c.Request.Method
	URL := c.Request.URL.String()

	errorSvc.LogError(currentUserId, errorMessage, httpRequestMethod, URL)
}
