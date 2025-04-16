package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type EvaluationStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EvaluationStatus, error)
	GetByID(id uint) (*models.EvaluationStatus, error)
	GetByType(status string) (*models.EvaluationStatus, error)
}

type EvaluationStatusRepository struct {
	db *gorm.DB
}

func NewEvaluationStatusRepository(db *gorm.DB) *EvaluationStatusRepository {
	return &EvaluationStatusRepository{
		db: db,
	}
}

func (repo *EvaluationStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.EvaluationStatus, error) {
	var statuses []*models.EvaluationStatus
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
	result := query.Find(&statuses)

	return statuses, result.Error
}

func (repo *EvaluationStatusRepository) GetByID(id uint) (*models.EvaluationStatus, error) {
	var status *models.EvaluationStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *EvaluationStatusRepository) GetByType(status string) (*models.EvaluationStatus, error) {
	var requirementStatus *models.EvaluationStatus
	result := repo.db.Where("type = ?", status).First(&requirementStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Evaluation Status with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return requirementStatus, nil
}
