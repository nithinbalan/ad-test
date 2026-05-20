package utility

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ContextErrorKey is the gin.Context key under which response helpers stash
// the human-readable error message they sent back to the client. The API
// request logging middleware reads it after the request completes so the
// error appears on the api_request_logs row without requiring handlers to
// change.
const ContextErrorKey = "api_log_error"

// recordContextError stores the error message on the gin context. Helpers
// below call this just before writing the response.
func recordContextError(c *gin.Context, msg string) {
	if c == nil || msg == "" {
		return
	}
	c.Set(ContextErrorKey, msg)
}

type InternalErrorResponse struct {
	Code            string `json:"code"`
	Error           string `json:"error"`
	InternalMessage string `json:"internal_message"`
}

type ErrorResponse struct {
	Code  string `json:"code"`
	Error string `json:"error"`
}

func SendInternalErrResp(db *gorm.DB, c *gin.Context, errMsg string) {
	recordContextError(c, errMsg)

	localizedErrorMsg := LocalizeMessage(c, "internal_server_error", "Internal Server Error!")

	log.Printf("Recovered from panic: %v", errMsg)

	// errorLogRepo := repositories.NewErrorLogRepository(db)
	// errorLogSvc := services.NewErrorLoggerService(errorLogRepo)
	// go LogErrorToDB(errorLogSvc, c, errMsg)

	resp := InternalErrorResponse{
		Code:            "INTERNAL_ERROR",
		Error:           localizedErrorMsg,
		InternalMessage: errMsg,
	}
	c.JSON(http.StatusInternalServerError, resp)
}

func SendInvalidErrResp(c *gin.Context, errs ...error) {
	var err error
	if len(errs) > 0 {
		err = errs[0]
	}
	if err != nil {
		recordContextError(c, err.Error())
	}

	localizedErrorMsg := LocalizeMessage(c, "invalid_request", "Invalid Request!")
	errorMessages := ExtractValidationErrors(c, err)

	c.JSON(http.StatusBadRequest, gin.H{
		"code":   "INVALID_REQUEST",
		"error":  localizedErrorMsg,
		"errors": errorMessages,
	})
}

func SendConflictErrResp(c *gin.Context, message string) {
	recordContextError(c, message)

	localizedErrorMsg := LocalizeMessage(c, "conflict", "Conflict!")
	c.JSON(http.StatusConflict, gin.H{
		"code":    "CONFLICT_ERROR",
		"error":   localizedErrorMsg,
		"message": message,
	})
}

func SendNotFoundErrResp(c *gin.Context) {
	recordContextError(c, "resource not found")

	resp := ErrorResponse{
		Code:  "RESOURCE_NOT_FOUND",
		Error: LocalizeMessage(c, "resource_not_found", "Resource Not Found!"),
	}
	c.JSON(http.StatusNotFound, resp)
}

func SendUnauthorizedErrResp(c *gin.Context, errMsg string) {
	if errMsg == "" {
		errMsg = "unauthorized"
	}
	recordContextError(c, errMsg)

	resp := ErrorResponse{
		Code:  "UNAUTHORIZED",
		Error: LocalizeMessage(c, "unauthorized_request", "Unauthorized Request"),
	}
	c.JSON(http.StatusUnauthorized, resp)
}

func SendForbiddenErrResp(c *gin.Context, errorMessage string) {
	recordContextError(c, errorMessage)

	resp := ErrorResponse{
		Code:  "FORBIDDEN",
		Error: errorMessage,
	}
	c.JSON(http.StatusForbidden, resp)
}

func SendCustomErrResp(c *gin.Context, statusCode int, message string) {
	if statusCode >= 400 {
		recordContextError(c, fmt.Sprintf("HTTP %d: %s", statusCode, message))
	}
	c.JSON(statusCode, gin.H{
		"code":    "CUSTOM_ERROR",
		"message": message,
	})
}

func SendLicenseExpiryErrResp(c *gin.Context, message string) {
	recordContextError(c, message)

	c.JSON(http.StatusForbidden, gin.H{
		"code":    "LICENSE_EXPIRED",
		"message": message,
	})
}
