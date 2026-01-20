package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ConsultationRequestStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ConsultationRequestStatus, error)
	GetByID(id uint) (*models.ConsultationRequestStatus, error)
	GetByType(status string) (*models.ConsultationRequestStatus, error)
}

// Implementation
type ConsultationRequestStatusRepository struct {
	db *gorm.DB
}

func NewConsultationRequestStatusRepository(db *gorm.DB) *ConsultationRequestStatusRepository {
	return &ConsultationRequestStatusRepository{
		db: db,
	}
}

func (repo *ConsultationRequestStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ConsultationRequestStatus, error) {
	var statuses []*models.ConsultationRequestStatus

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ? OR description_en LIKE ? OR description_ar LIKE ?)",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
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

	result := query.Order(sortBy + " " + sort).Find(&statuses)
	return statuses, result.Error
}

func (repo *ConsultationRequestStatusRepository) GetByID(id uint) (*models.ConsultationRequestStatus, error) {
	var status *models.ConsultationRequestStatus
	result := repo.db.First(&status, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *ConsultationRequestStatusRepository) GetByType(status string) (*models.ConsultationRequestStatus, error) {
	var consultationStatus *models.ConsultationRequestStatus
	result := repo.db.Where("type = ?", status).First(&consultationStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Consultation Request Status with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return consultationStatus, nil
}
