package utility

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

	localizedErrorMsg := LocalizeMessage(c, "invalid_request", "Invalid Request!")
	errorMessages := ExtractValidationErrors(c, err)

	c.JSON(http.StatusBadRequest, gin.H{
		"code":   "INVALID_REQUEST",
		"error":  localizedErrorMsg,
		"errors": errorMessages,
	})
}

func SendConflictErrResp(c *gin.Context, message string) {
	localizedErrorMsg := LocalizeMessage(c, "conflict", "Conflict!")
	c.JSON(http.StatusConflict, gin.H{
		"code":    "CONFLICT_ERROR",
		"error":   localizedErrorMsg,
		"message": message,
	})
}

func SendNotFoundErrResp(c *gin.Context) {
	resp := ErrorResponse{
		Code:  "RESOURCE_NOT_FOUND",
		Error: LocalizeMessage(c, "resource_not_found", "Resource Not Found!"),
	}
	c.JSON(http.StatusNotFound, resp)
}

func SendUnauthorizedErrResp(c *gin.Context, errMsg string) {
	resp := ErrorResponse{
		Code:  "UNAUTHORIZED",
		Error: LocalizeMessage(c, "unauthorized_request", "Unauthorized Request"),
	}
	c.JSON(http.StatusUnauthorized, resp)
}

func SendForbiddenErrResp(c *gin.Context, errorMessage string) {
	resp := ErrorResponse{
		Code:  "FORBIDDEN",
		Error: errorMessage,
	}
	c.JSON(http.StatusForbidden, resp)
}

func SendCustomErrResp(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"code":    "CUSTOM_ERROR",
		"message": message,
	})
}

func SendLicenseExpiryErrResp(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, gin.H{
		"code":    "LICENSE_EXPIRED",
		"message": message,
	})
}
