package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type MeetingStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingStatus, error)
	GetByID(id uint) (*models.MeetingStatus, error)
}

// implementation
type MeetingStatusRepository struct {
	db *gorm.DB
}

func NewMeetingStatusRepository(db *gorm.DB) *MeetingStatusRepository {
	return &MeetingStatusRepository{
		db: db,
	}
}

func (repo *MeetingStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingStatus, error) {

	var meetingStatuses []*models.MeetingStatus

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
	result := query.Find(&meetingStatuses)

	return meetingStatuses, result.Error
}

func (repo *MeetingStatusRepository) GetByID(id uint) (*models.MeetingStatus, error) {
	var status *models.MeetingStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}
