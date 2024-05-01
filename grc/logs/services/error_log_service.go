package services

import (
	"jethur/grc/logs/repositories"
	"jethur/grc/models"
)

type ErrorLogServiceInterface interface {
	LogError(errorLogEntry models.ErrorLog) error
}

type ErrorLoggerService struct {
	repo repositories.ErrorLogRepositoryInterface
}

func NewErrorLoggerService(repo repositories.ErrorLogRepositoryInterface) ErrorLogServiceInterface {
	return &ErrorLoggerService{repo: repo}
}

func (els *ErrorLoggerService) LogError(errorLogEntry models.ErrorLog) error {
	return els.repo.LogError(errorLogEntry)
}
