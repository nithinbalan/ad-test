package services

import (
	"errors"
	"jethur/grc/meeting/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type MeetingStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingStatus, error)
	GetByID(id uint) (*models.MeetingStatus, error)
}

// Service and implementations
type MeetingStatusService struct {
	repo repositories.MeetingStatusRepositoryInterface
}

func NewMeetingStatusService(repo repositories.MeetingStatusRepositoryInterface) MeetingStatusServiceInterface {
	return &MeetingStatusService{
		repo: repo,
	}
}

func (mss *MeetingStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingStatus, error) {
	return mss.repo.GetAll(sort, sortBy, keyword, status)
}

func (mss MeetingStatusService) GetByID(id uint) (*models.MeetingStatus, error) {
	meetingStatus, err := mss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meetingStatus, nil
}
