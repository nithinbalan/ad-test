package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditDraftFindingStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.AuditDraftFindingStatus, error)
	GetByID(id uint) (*models.AuditDraftFindingStatus, error)
	GetByType(status string) (*models.AuditDraftFindingStatus, error)
}

// Implementation
type AuditDraftFindingStatusRepository struct {
	db *gorm.DB
}

func NewAuditDraftFindingStatusRepository(db *gorm.DB) *AuditDraftFindingStatusRepository {
	return &AuditDraftFindingStatusRepository{
		db: db,
	}
}

func (repo *AuditDraftFindingStatusRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.AuditDraftFindingStatus, error) {

	var auditDraftFindingStatuses []*models.AuditDraftFindingStatus

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
	result := query.Find(&auditDraftFindingStatuses)

	return auditDraftFindingStatuses, result.Error
}

func (repo *AuditDraftFindingStatusRepository) GetByID(id uint) (*models.AuditDraftFindingStatus, error) {
	var status *models.AuditDraftFindingStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *AuditDraftFindingStatusRepository) GetByType(status string) (*models.AuditDraftFindingStatus, error) {
	var auditDraftFindingStatus *models.AuditDraftFindingStatus

	result := repo.db.Where("type = ?", status).First(&auditDraftFindingStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Draft Finding Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditDraftFindingStatus, nil
}
