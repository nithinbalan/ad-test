package utility

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendErrorResponse(c *gin.Context, statusCode int, errMsg string) {
	resp := ErrorResponse{
		Error: errMsg,
	}
	c.JSON(statusCode, resp)
}
