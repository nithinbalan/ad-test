package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditStatus, error)
	GetByID(id uint) (*models.AuditStatus, error)
	GetByType(status string) (*models.AuditStatus, error)
}

// implementation
type AuditStatusRepository struct {
	db *gorm.DB
}

func NewAuditStatusRepository(db *gorm.DB) *AuditStatusRepository {
	return &AuditStatusRepository{
		db: db,
	}
}

func (repo *AuditStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditStatus, error) {

	var auditStatuses []*models.AuditStatus

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
	result := query.Find(&auditStatuses)

	return auditStatuses, result.Error
}

func (repo *AuditStatusRepository) GetByID(id uint) (*models.AuditStatus, error) {
	var status *models.AuditStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *AuditStatusRepository) GetByType(status string) (*models.AuditStatus, error) {
	var auditStatus *models.AuditStatus

	result := repo.db.Where("type = ?", status).First(&auditStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditStatus, nil
}
