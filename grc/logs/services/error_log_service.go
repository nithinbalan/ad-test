package services

import (
	"fmt"
	"jethur/grc/logs/models"
	"jethur/grc/logs/repositories"
	"net/http"
	"runtime/debug"
	"time"
)

type ErrorLogServiceInterface interface {
	LogError(currentUserID uint, errorMessage string, httpRequestMethod, URL string)
}

type ErrorLoggerService struct {
	repo repositories.ErrorLogRepositoryInterface
}

func NewErrorLoggerService(repo repositories.ErrorLogRepositoryInterface) ErrorLogServiceInterface {
	return &ErrorLoggerService{repo: repo}
}

func (s *ErrorLoggerService) LogError(currentUserID uint, errorMessage string, httpRequestMethod, URL string) {

	stackTrace := string(debug.Stack())

	errorLogEntry := models.ErrorLog{
		HTTPStatusCode:    http.StatusInternalServerError,
		HTTPRequestMethod: httpRequestMethod,
		URL:               URL,
		ErrorMessage:      errorMessage,
		StackTrace:        stackTrace,
		UserID:            &currentUserID,
		CreatedAt:         time.Now(),
	}

	err := s.repo.LogError(errorLogEntry)
	if err != nil {
		// @todo write to log file
		fmt.Println("write to log file")
	}
}
