package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ProjectStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectStatus, error)
	GetByID(id uint) (*models.ProjectStatus, error)
	GetByType(status string) (*models.ProjectStatus, error)
}

// implementation
type ProjectStatusRepository struct {
	db *gorm.DB
}

func NewProjectStatusRepository(db *gorm.DB) *ProjectStatusRepository {
	return &ProjectStatusRepository{
		db: db,
	}
}

func (repo *ProjectStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectStatus, error) {
	var projectStatuses []*models.ProjectStatus
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
	result := query.Find(&projectStatuses)

	return projectStatuses, result.Error
}

func (repo *ProjectStatusRepository) GetByID(id uint) (*models.ProjectStatus, error) {
	var status *models.ProjectStatus
	result := repo.db.First(&status, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *ProjectStatusRepository) GetByType(status string) (*models.ProjectStatus, error) {
	var projectStatus *models.ProjectStatus

	result := repo.db.Where("type = ?", status).First(&projectStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Project Status with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return projectStatus, nil
}
