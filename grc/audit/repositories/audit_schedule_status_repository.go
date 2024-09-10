package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditScheduleStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditScheduleStatus, error)
	GetByID(id uint) (*models.AuditScheduleStatus, error)
	GetByType(status string) (*models.AuditScheduleStatus, error)
}

// implementation
type AuditScheduleStatusRepository struct {
	db *gorm.DB
}

func NewAuditScheduleStatusRepository(db *gorm.DB) *AuditScheduleStatusRepository {
	return &AuditScheduleStatusRepository{
		db: db,
	}
}

func (repo *AuditScheduleStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditScheduleStatus, error) {

	var auditScheduleStatuses []*models.AuditScheduleStatus

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
	result := query.Find(&auditScheduleStatuses)

	return auditScheduleStatuses, result.Error
}

func (repo *AuditScheduleStatusRepository) GetByID(id uint) (*models.AuditScheduleStatus, error) {
	var status *models.AuditScheduleStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *AuditScheduleStatusRepository) GetByType(status string) (*models.AuditScheduleStatus, error) {
	var auditScheduleStatus *models.AuditScheduleStatus

	result := repo.db.Where("type = ?", status).First(&auditScheduleStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditScheduleStatus, nil
}
