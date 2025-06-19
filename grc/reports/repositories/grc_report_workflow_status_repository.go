package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type GRCReportWorkflowStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.GRCReportWorkflowStatus, error)
	GetByID(id uint) (*models.GRCReportWorkflowStatus, error)
	GetByType(status string) (*models.GRCReportWorkflowStatus, error)
}

// implementation
type GRCReportWorkflowStatusRepository struct {
	db *gorm.DB
}

func NewGRCReportWorkflowStatusRepository(db *gorm.DB) *GRCReportWorkflowStatusRepository {
	return &GRCReportWorkflowStatusRepository{
		db: db,
	}
}

func (repo *GRCReportWorkflowStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.GRCReportWorkflowStatus, error) {
	var grcReportWorkflowStatuses []*models.GRCReportWorkflowStatus

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
	result := query.Find(&grcReportWorkflowStatuses)

	return grcReportWorkflowStatuses, result.Error
}

func (repo *GRCReportWorkflowStatusRepository) GetByID(id uint) (*models.GRCReportWorkflowStatus, error) {
	var status *models.GRCReportWorkflowStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *GRCReportWorkflowStatusRepository) GetByType(status string) (*models.GRCReportWorkflowStatus, error) {
	var grcReportWorkflowStatus *models.GRCReportWorkflowStatus

	result := repo.db.Where("type = ?", status).First(&grcReportWorkflowStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return grcReportWorkflowStatus, nil
}
