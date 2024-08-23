package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type MeetingMinuteStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingMinuteStatus, error)
	GetByID(id uint) (*models.MeetingMinuteStatus, error)
	GetByType(status string) (*models.MeetingMinuteStatus, error)
}

// implementation
type MeetingMinuteStatusRepository struct {
	db *gorm.DB
}

func NewMeetingMinuteStatusRepository(db *gorm.DB) *MeetingMinuteStatusRepository {
	return &MeetingMinuteStatusRepository{
		db: db,
	}
}

func (repo *MeetingMinuteStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingMinuteStatus, error) {

	var meetingMinuteStatuses []*models.MeetingMinuteStatus

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&meetingMinuteStatuses)

	return meetingMinuteStatuses, result.Error
}

func (repo *MeetingMinuteStatusRepository) GetByID(id uint) (*models.MeetingMinuteStatus, error) {
	var status *models.MeetingMinuteStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *MeetingMinuteStatusRepository) GetByType(status string) (*models.MeetingMinuteStatus, error) {
	var meetingMinuteStatus *models.MeetingMinuteStatus

	result := repo.db.Where("type = ?", status).First(&meetingMinuteStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Meeting Minute Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return meetingMinuteStatus, nil
}
