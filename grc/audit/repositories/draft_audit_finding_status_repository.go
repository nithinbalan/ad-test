package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type DraftAuditFindingStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.DraftAuditFindingStatus, error)
	GetByID(id uint) (*models.DraftAuditFindingStatus, error)
	GetByType(status string) (*models.DraftAuditFindingStatus, error)
}

// Implementation
type DraftAuditFindingStatusRepository struct {
	db *gorm.DB
}

func NewDraftAuditFindingStatusRepository(db *gorm.DB) *DraftAuditFindingStatusRepository {
	return &DraftAuditFindingStatusRepository{
		db: db,
	}
}

func (repo *DraftAuditFindingStatusRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.DraftAuditFindingStatus, error) {

	var draftAuditFindingStatuses []*models.DraftAuditFindingStatus

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
	result := query.Find(&draftAuditFindingStatuses)

	return draftAuditFindingStatuses, result.Error
}

func (repo *DraftAuditFindingStatusRepository) GetByID(id uint) (*models.DraftAuditFindingStatus, error) {
	var status *models.DraftAuditFindingStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *DraftAuditFindingStatusRepository) GetByType(status string) (*models.DraftAuditFindingStatus, error) {
	var draftAuditFindingStatus *models.DraftAuditFindingStatus

	result := repo.db.Where("type = ?", status).First(&draftAuditFindingStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Draft Audit Finding Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return draftAuditFindingStatus, nil
}
