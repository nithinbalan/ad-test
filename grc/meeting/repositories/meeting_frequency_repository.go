package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type MeetingFrequencyRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingFrequency, error)
	GetByID(id uint) (*models.MeetingFrequency, error)
}

// implementation
type MeetingFrequencyRepository struct {
	db *gorm.DB
}

func NewMeetingFrequencyRepository(db *gorm.DB) *MeetingFrequencyRepository {
	return &MeetingFrequencyRepository{
		db: db,
	}
}

func (repo *MeetingFrequencyRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingFrequency, error) {

	var meetingFrequencies []*models.MeetingFrequency

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
	result := query.Find(&meetingFrequencies)

	return meetingFrequencies, result.Error
}

func (repo *MeetingFrequencyRepository) GetByID(id uint) (*models.MeetingFrequency, error) {
	var meetingFrequency *models.MeetingFrequency
	result := repo.db.First(&meetingFrequency, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return meetingFrequency, result.Error
}
