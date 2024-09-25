package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditWorkflowStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditWorkflowStatus, error)
	GetByID(id uint) (*models.AuditWorkflowStatus, error)
	GetByType(status string) (*models.AuditWorkflowStatus, error)
}

// Implementation
type AuditWorkflowStatusRepository struct {
	db *gorm.DB
}

func NewAuditWorkflowStatusRepository(db *gorm.DB) *AuditWorkflowStatusRepository {
	return &AuditWorkflowStatusRepository{
		db: db,
	}
}

func (repo *AuditWorkflowStatusRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditWorkflowStatus, error) {

	var auditWorkflowStatuses []*models.AuditWorkflowStatus

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
	result := query.Find(&auditWorkflowStatuses)

	return auditWorkflowStatuses, result.Error
}

func (repo *AuditWorkflowStatusRepository) GetByID(id uint) (*models.AuditWorkflowStatus, error) {
	var status *models.AuditWorkflowStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *AuditWorkflowStatusRepository) GetByType(status string) (*models.AuditWorkflowStatus, error) {
	var auditWorkflowStatus *models.AuditWorkflowStatus

	result := repo.db.Where("type = ?", status).First(&auditWorkflowStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Draft Finding Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditWorkflowStatus, nil
}
