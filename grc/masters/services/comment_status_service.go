package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type CommentStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommentStatus, error)
	GetByID(id uint) (*models.CommentStatus, error)
	GetByType(status string) (*models.CommentStatus, error)
}

// Service and implementations
type CommentStatusService struct {
	repo repositories.CommentStatusRepositoryInterface
}

func NewCommentStatusService(repo repositories.CommentStatusRepositoryInterface) CommentStatusServiceInterface {
	return &CommentStatusService{
		repo: repo,
	}
}

func (apss *CommentStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommentStatus, error) {
	return apss.repo.GetAll(sort, sortBy, keyword, status)
}

func (apss CommentStatusService) GetByID(id uint) (*models.CommentStatus, error) {
	return apss.repo.GetByID(id)
}

func (apss *CommentStatusService) GetByType(status string) (*models.CommentStatus, error) {
	// Fetch the action plan status by its status
	statusResult, err := apss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
