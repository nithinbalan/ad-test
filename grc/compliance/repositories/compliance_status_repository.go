package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ComplianceStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceStatus, error)
	GetByID(id uint) (*models.ComplianceStatus, error)

	GetMaxBaseScore(tx ...utility.DatabaseTransaction) (int, error)
	GetByType(status string) (*models.ComplianceStatus, error)
}

// implementation
type ComplianceStatusRepository struct {
	db *gorm.DB
}

func NewComplianceStatusRepository(db *gorm.DB) *ComplianceStatusRepository {
	return &ComplianceStatusRepository{
		db: db,
	}
}

func (repo *ComplianceStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceStatus, error) {

	var complianceStatuses []*models.ComplianceStatus

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
	result := query.Find(&complianceStatuses)

	return complianceStatuses, result.Error
}

func (repo *ComplianceStatusRepository) GetByID(id uint) (*models.ComplianceStatus, error) {
	var status *models.ComplianceStatus
	result := repo.db.First(&status, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *ComplianceStatusRepository) GetMaxBaseScore(tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var maxBaseScore int
	err := db.Table("compliance_statuses").Select("MAX(base_score)").Scan(&maxBaseScore).Error
	if err != nil {
		return 0, err
	}
	return maxBaseScore, nil
}

func (repo *ComplianceStatusRepository) GetByType(status string) (*models.ComplianceStatus, error) {
	var complianceStatus *models.ComplianceStatus

	result := repo.db.Where("type = ?", status).First(&complianceStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Compliance Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return complianceStatus, nil
}
