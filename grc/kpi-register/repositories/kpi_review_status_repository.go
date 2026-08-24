package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type KPIReviewStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIReviewStatus, error)
	GetByID(id uint) (*models.KPIReviewStatus, error)
	GetByType(status string) (*models.KPIReviewStatus, error)
}

// implementation
type KPIReviewStatusRepository struct {
	db *gorm.DB
}

func NewKPIReviewStatusRepository(db *gorm.DB) *KPIReviewStatusRepository {
	return &KPIReviewStatusRepository{
		db: db,
	}
}

func (repo *KPIReviewStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIReviewStatus, error) {

	var kpiReviewStatuses []*models.KPIReviewStatus

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
	result := query.Find(&kpiReviewStatuses)

	return kpiReviewStatuses, result.Error
}

func (repo *KPIReviewStatusRepository) GetByID(id uint) (*models.KPIReviewStatus, error) {
	var kpiReviewStatus *models.KPIReviewStatus
	result := repo.db.First(&kpiReviewStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return kpiReviewStatus, result.Error
}

func (repo *KPIReviewStatusRepository) GetByType(status string) (*models.KPIReviewStatus, error) {
	var kpiReviewStatus *models.KPIReviewStatus

	result := repo.db.Where("type = ?", status).First(&kpiReviewStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "KPI Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return kpiReviewStatus, nil
}
