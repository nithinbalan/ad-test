package services

import (
	"jethur/grc/meeting/repositories"
	"jethur/grc/models"
)

type MeetingFrequencyServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingFrequency, error)
	GetByID(id uint) (*models.MeetingFrequency, error)
}

// Service and implementations
type MeetingFrequencyService struct {
	repo repositories.MeetingFrequencyRepositoryInterface
}

func NewMeetingFrequencyService(repo repositories.MeetingFrequencyRepositoryInterface) MeetingFrequencyServiceInterface {
	return &MeetingFrequencyService{
		repo: repo,
	}
}

func (rfs *MeetingFrequencyService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingFrequency, error) {
	return rfs.repo.GetAll(sort, sortBy, keyword, status)
}

func (rfs MeetingFrequencyService) GetByID(id uint) (*models.MeetingFrequency, error) {
	return rfs.repo.GetByID(id)
}
