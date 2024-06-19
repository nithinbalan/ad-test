package utility

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func SendInternalErrResp(db *gorm.DB, c *gin.Context, errMsg string) {
	localizedErrorMsg := LocalizeMessage(c, "internal_server_error", "Internal Server Error!")

	log.Printf("Recovered from panic: %v", errMsg)

	// log.Printf("Recovered from panic: %v\nStack trace: %s", errMsg)

	// errorLogRepo := repositories.NewErrorLogRepository(db)
	// errorLogSvc := services.NewErrorLoggerService(errorLogRepo)

	// go LogErrorToDB(errorLogSvc, c, errMsg)

	resp := ErrorResponse{Error: localizedErrorMsg}
	c.JSON(http.StatusInternalServerError, resp)
}

func SendInvalidErrResp(c *gin.Context, errs ...error) {

	var err error
	if len(errs) > 0 {
		err = errs[0] // Use the first error if any are provided
	}

	localizedErrorMsg := LocalizeMessage(c, "invalid_request", "Invalid Request!")
	errorMessages := ExtractValidationErrors(c, err)
	c.JSON(http.StatusBadRequest, gin.H{
		"error":  localizedErrorMsg,
		"errors": errorMessages,
	})

}

func SendNotFoundErrResp(c *gin.Context) {
	resp := ErrorResponse{Error: LocalizeMessage(c, "resource_not_found", "Resource Not Found!")}
	c.JSON(http.StatusNotFound, resp)
}

func SendUnauthorizedErrResp(c *gin.Context, errMsg string) {
	resp := ErrorResponse{Error: LocalizeMessage(c, "unauthorized_request", "Unauthorized Request")}
	c.JSON(http.StatusUnauthorized, resp)
}

func SendForbiddenErrResp(c *gin.Context, errorMessage string) {
	resp := ErrorResponse{Error: errorMessage}
	c.JSON(http.StatusForbidden, resp)
}
