package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type CommunicationActionPlanStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommunicationActionPlanStatus, error)
	GetByID(id uint) (*models.CommunicationActionPlanStatus, error)
	GetByType(status string) (*models.CommunicationActionPlanStatus, error)
}

// implementation
type CommunicationActionPlanStatusRepository struct {
	db *gorm.DB
}

func NewCommunicationActionPlanStatusRepository(db *gorm.DB) *CommunicationActionPlanStatusRepository {
	return &CommunicationActionPlanStatusRepository{
		db: db,
	}
}

func (repo *CommunicationActionPlanStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommunicationActionPlanStatus, error) {
	var testAndExerciseStatuses []*models.CommunicationActionPlanStatus

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
	result := query.Find(&testAndExerciseStatuses)

	return testAndExerciseStatuses, result.Error
}

func (repo *CommunicationActionPlanStatusRepository) GetByID(id uint) (*models.CommunicationActionPlanStatus, error) {
	var status *models.CommunicationActionPlanStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *CommunicationActionPlanStatusRepository) GetByType(status string) (*models.CommunicationActionPlanStatus, error) {
	var testAndExerciseStatus *models.CommunicationActionPlanStatus

	result := repo.db.Where("type = ?", status).First(&testAndExerciseStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExerciseStatus, nil
}
