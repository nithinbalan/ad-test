package services

import (
	"errors"
	"jethur/grc/meeting/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type MeetingMinuteStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingMinuteStatus, error)
	GetByID(id uint) (*models.MeetingMinuteStatus, error)
	GetByType(status string) (*models.MeetingMinuteStatus, error)
}

// Service and implementations
type MeetingMinuteStatusService struct {
	repo repositories.MeetingMinuteStatusRepositoryInterface
}

func NewMeetingMinuteStatusService(repo repositories.MeetingMinuteStatusRepositoryInterface) MeetingMinuteStatusServiceInterface {
	return &MeetingMinuteStatusService{
		repo: repo,
	}
}

func (mmss *MeetingMinuteStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingMinuteStatus, error) {
	return mmss.repo.GetAll(sort, sortBy, keyword, status)
}

func (mmss MeetingMinuteStatusService) GetByID(id uint) (*models.MeetingMinuteStatus, error) {
	meetingMinuteStatus, err := mmss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meetingMinuteStatus, nil
}

func (mmss *MeetingMinuteStatusService) GetByType(status string) (*models.MeetingMinuteStatus, error) {
	statusResult, err := mmss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
