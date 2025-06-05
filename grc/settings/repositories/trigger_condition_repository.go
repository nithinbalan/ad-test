package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type TriggerConditionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TriggerCondition, error)
	GetByID(id uint) (*models.TriggerCondition, error)
	GetByType(status string) (*models.TriggerCondition, error)
}

// implementation
type TriggerConditionRepository struct {
	db *gorm.DB
}

func NewTriggerConditionRepository(db *gorm.DB) *TriggerConditionRepository {
	return &TriggerConditionRepository{
		db: db,
	}
}

func (repo *TriggerConditionRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TriggerCondition, error) {
	var triggerConditions []*models.TriggerCondition

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
	result := query.Find(&triggerConditions)

	return triggerConditions, result.Error
}

func (repo *TriggerConditionRepository) GetByID(id uint) (*models.TriggerCondition, error) {
	var status *models.TriggerCondition

	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *TriggerConditionRepository) GetByType(status string) (*models.TriggerCondition, error) {
	var triggerCondition *models.TriggerCondition

	result := repo.db.Where("status = ?", status).First(&triggerCondition)

	if result.Error != nil {
		return nil, result.Error
	}

	return triggerCondition, nil
}
