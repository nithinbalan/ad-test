package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type KRIRiskToleranceStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KRIRiskToleranceStatus, error)
	GetByID(id uint) (*models.KRIRiskToleranceStatus, error)
	GetByType(statusType string) (*models.KRIRiskToleranceStatus, error)
}

type KRIRiskToleranceStatusRepository struct {
	db *gorm.DB
}

func NewKRIRiskToleranceStatusRepository(db *gorm.DB) *KRIRiskToleranceStatusRepository {
	return &KRIRiskToleranceStatusRepository{
		db: db,
	}
}

func (repo *KRIRiskToleranceStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KRIRiskToleranceStatus, error) {
	var statuses []*models.KRIRiskToleranceStatus
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

func (repo *KRIRiskToleranceStatusRepository) GetByID(id uint) (*models.KRIRiskToleranceStatus, error) {
	var status *models.KRIRiskToleranceStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *KRIRiskToleranceStatusRepository) GetByType(statusType string) (*models.KRIRiskToleranceStatus, error) {
	var status *models.KRIRiskToleranceStatus

	result := repo.db.Where("type = ?", statusType).First(&status)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "KRI Risk Tolerance Status with type "+statusType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}
