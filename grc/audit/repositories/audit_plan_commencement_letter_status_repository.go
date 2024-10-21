package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type AuditPlanCommencementLetterStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditPlanCommencementLetterStatus, error)
	GetByID(id uint) (*models.AuditPlanCommencementLetterStatus, error)
	GetByType(status string) (*models.AuditPlanCommencementLetterStatus, error)
}

// Implementation
type AuditPlanCommencementLetterStatusRepository struct {
	db *gorm.DB
}

func NewAuditPlanCommencementLetterStatusRepository(db *gorm.DB) *AuditPlanCommencementLetterStatusRepository {
	return &AuditPlanCommencementLetterStatusRepository{
		db: db,
	}
}

// GetAll retrieves all statuses, with optional sorting and filtering
func (repo *AuditPlanCommencementLetterStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AuditPlanCommencementLetterStatus, error) {
	var statuses []*models.AuditPlanCommencementLetterStatus

	query := repo.db

	// Apply keyword filtering
	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	// Apply sorting, with default to ascending on title_en if not provided
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Filter by status (active, inactive)
	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	// Apply ordering
	query = query.Order(sortBy + " " + sort)

	// Execute the query
	result := query.Find(&statuses)
	return statuses, result.Error
}

// GetByID retrieves a specific status by its ID
func (repo *AuditPlanCommencementLetterStatusRepository) GetByID(id uint) (*models.AuditPlanCommencementLetterStatus, error) {
	var status *models.AuditPlanCommencementLetterStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

// GetByType retrieves a status by its type
func (repo *AuditPlanCommencementLetterStatusRepository) GetByType(status string) (*models.AuditPlanCommencementLetterStatus, error) {
	var auditStatus *models.AuditPlanCommencementLetterStatus

	result := repo.db.Where("type = ?", status).First(&auditStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit Plan Commencement Letter Status with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditStatus, nil
}
