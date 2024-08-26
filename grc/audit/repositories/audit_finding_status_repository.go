package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditFindingStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditFindingStatus, error)
	GetByID(id uint) (*models.AuditFindingStatus, error)
	GetByType(status string) (*models.AuditFindingStatus, error)
}

// Implementation
type AuditFindingStatusRepository struct {
	db *gorm.DB
}

func NewAuditFindingStatusRepository(db *gorm.DB) *AuditFindingStatusRepository {
	return &AuditFindingStatusRepository{
		db: db,
	}
}

func (repo *AuditFindingStatusRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditFindingStatus, error) {

	var auditFindingStatuses []*models.AuditFindingStatus

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

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&auditFindingStatuses)

	return auditFindingStatuses, result.Error
}

func (repo *AuditFindingStatusRepository) GetByID(id uint) (*models.AuditFindingStatus, error) {
	var status *models.AuditFindingStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *AuditFindingStatusRepository) GetByType(status string) (*models.AuditFindingStatus, error) {
	var auditFindingStatus *models.AuditFindingStatus

	result := repo.db.Where("type = ?", status).First(&auditFindingStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Finding Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditFindingStatus, nil
}
