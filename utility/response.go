package utility

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendErrorResponse(c *gin.Context, statusCode int, errMsg string) {

	var errorType string = ""
	var localizedErrorMsg string = ""

	if statusCode == http.StatusInternalServerError {
		errorType = "500"
		localizedErrorMsg = LocalizeMessage(c, "internal_server_error", "Internal server error!")
	}

	//log error on database
	go LogException(c, errorType, errMsg)

	resp := ErrorResponse{
		Error: localizedErrorMsg,
	}
	c.JSON(statusCode, resp)
}
