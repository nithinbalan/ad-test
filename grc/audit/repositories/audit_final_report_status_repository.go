package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditFinalReportStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditFinalReportStatus, error)
	GetByID(id uint) (*models.AuditFinalReportStatus, error)
	GetByType(status string) (*models.AuditFinalReportStatus, error)
}

// Implementation
type AuditFinalReportStatusRepository struct {
	db *gorm.DB
}

func NewAuditFinalReportStatusRepository(db *gorm.DB) *AuditFinalReportStatusRepository {
	return &AuditFinalReportStatusRepository{
		db: db,
	}
}

func (repo *AuditFinalReportStatusRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditFinalReportStatus, error) {

	var auditFinalReportStatuses []*models.AuditFinalReportStatus

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
	result := query.Find(&auditFinalReportStatuses)

	return auditFinalReportStatuses, result.Error
}

func (repo *AuditFinalReportStatusRepository) GetByID(id uint) (*models.AuditFinalReportStatus, error) {
	var status *models.AuditFinalReportStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *AuditFinalReportStatusRepository) GetByType(status string) (*models.AuditFinalReportStatus, error) {
	var auditFinalReportStatus *models.AuditFinalReportStatus

	result := repo.db.Where("type = ?", status).First(&auditFinalReportStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Final Report Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditFinalReportStatus, nil
}
