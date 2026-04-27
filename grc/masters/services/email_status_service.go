package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type EmailStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EmailStatus, error)
	GetByID(id uint) (*models.EmailStatus, error)
	GetByType(status string) (*models.EmailStatus, error)
}

type EmailStatusService struct {
	repo repositories.EmailStatusRepositoryInterface
}

func NewEmailStatusService(repo repositories.EmailStatusRepositoryInterface) EmailStatusServiceInterface {
	return &EmailStatusService{
		repo: repo,
	}
}

func (s *EmailStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EmailStatus, error) {
	return s.repo.GetAll(sort, sortBy, keyword, status)
}

func (s *EmailStatusService) GetByID(id uint) (*models.EmailStatus, error) {
	return s.repo.GetByID(id)
}

func (s *EmailStatusService) GetByType(status string) (*models.EmailStatus, error) {
	return s.repo.GetByType(status)
}
