package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditDocumentRequestStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditDocumentRequestStatus, error)
	GetByID(id uint) (*models.AuditDocumentRequestStatus, error)
	GetByType(status string) (*models.AuditDocumentRequestStatus, error)
}

// implementation
type AuditDocumentRequestStatusRepository struct {
	db *gorm.DB
}

func NewAuditDocumentRequestStatusRepository(db *gorm.DB) *AuditDocumentRequestStatusRepository {
	return &AuditDocumentRequestStatusRepository{
		db: db,
	}
}

func (repo *AuditDocumentRequestStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditDocumentRequestStatus, error) {

	var auditDocumentRequestStatus []*models.AuditDocumentRequestStatus

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
	result := query.Find(&auditDocumentRequestStatus)

	return auditDocumentRequestStatus, result.Error
}

func (repo *AuditDocumentRequestStatusRepository) GetByID(id uint) (*models.AuditDocumentRequestStatus, error) {
	var status *models.AuditDocumentRequestStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *AuditDocumentRequestStatusRepository) GetByType(status string) (*models.AuditDocumentRequestStatus, error) {
	var auditDocumentRequestStatus *models.AuditDocumentRequestStatus

	result := repo.db.Where("type = ?", status).First(&auditDocumentRequestStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Document Request Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditDocumentRequestStatus, nil
}
