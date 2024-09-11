package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditDraftReportStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditDraftReportStatus, error)
	GetByID(id uint) (*models.AuditDraftReportStatus, error)
	GetByType(status string) (*models.AuditDraftReportStatus, error)
}

// Implementation
type AuditDraftReportStatusRepository struct {
	db *gorm.DB
}

func NewAuditDraftReportStatusRepository(db *gorm.DB) *AuditDraftReportStatusRepository {
	return &AuditDraftReportStatusRepository{
		db: db,
	}
}

func (repo *AuditDraftReportStatusRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditDraftReportStatus, error) {

	var auditDraftReportStatuses []*models.AuditDraftReportStatus

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
	result := query.Find(&auditDraftReportStatuses)

	return auditDraftReportStatuses, result.Error
}

func (repo *AuditDraftReportStatusRepository) GetByID(id uint) (*models.AuditDraftReportStatus, error) {
	var status *models.AuditDraftReportStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *AuditDraftReportStatusRepository) GetByType(status string) (*models.AuditDraftReportStatus, error) {
	var auditDraftReportStatus *models.AuditDraftReportStatus

	result := repo.db.Where("type = ?", status).First(&auditDraftReportStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Draft Finding Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftReportStatus, nil
}
