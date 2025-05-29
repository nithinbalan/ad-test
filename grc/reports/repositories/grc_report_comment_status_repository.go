package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type GRCReportCommentStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.GRCReportCommentStatus, error)
	GetByID(id uint) (*models.GRCReportCommentStatus, error)
	GetByType(status string) (*models.GRCReportCommentStatus, error)
}

// implementation
type GRCReportCommentStatusRepository struct {
	db *gorm.DB
}

func NewGRCReportCommentStatusRepository(db *gorm.DB) *GRCReportCommentStatusRepository {
	return &GRCReportCommentStatusRepository{
		db: db,
	}
}

func (repo *GRCReportCommentStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.GRCReportCommentStatus, error) {
	var grcReportCommentStatuses []*models.GRCReportCommentStatus

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
	result := query.Find(&grcReportCommentStatuses)

	return grcReportCommentStatuses, result.Error
}

func (repo *GRCReportCommentStatusRepository) GetByID(id uint) (*models.GRCReportCommentStatus, error) {
	var status *models.GRCReportCommentStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *GRCReportCommentStatusRepository) GetByType(status string) (*models.GRCReportCommentStatus, error) {
	var grcReportCommentStatus *models.GRCReportCommentStatus

	result := repo.db.Where("type = ?", status).First(&grcReportCommentStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return grcReportCommentStatus, nil
}
