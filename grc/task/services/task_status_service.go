package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/task/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type TaskStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TaskStatus, error)
	GetByID(id uint) (*models.TaskStatus, error)
	GetByType(status string) (*models.TaskStatus, error)
}

// Service and implementations
type TaskStatusService struct {
	repo repositories.TaskStatusRepositoryInterface
}

func NewTaskStatusService(repo repositories.TaskStatusRepositoryInterface) TaskStatusServiceInterface {
	return &TaskStatusService{
		repo: repo,
	}
}

func (tss *TaskStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TaskStatus, error) {
	return tss.repo.GetAll(sort, sortBy, keyword, status)
}

func (tss TaskStatusService) GetByID(id uint) (*models.TaskStatus, error) {
	taskStatus, err := tss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return taskStatus, nil
}

func (tss *TaskStatusService) GetByType(status string) (*models.TaskStatus, error) {
	// Fetch the task status by its status
	statusResult, err := tss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
