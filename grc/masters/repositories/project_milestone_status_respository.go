package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ProjectMilestoneStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectMilestoneStatus, error)
	GetByID(id uint) (*models.ProjectMilestoneStatus, error)
	GetByType(status string) (*models.ProjectMilestoneStatus, error)
}

// implementation
type ProjectMilestoneStatusRepository struct {
	db *gorm.DB
}

func NewProjectMilestoneStatusRepository(db *gorm.DB) *ProjectMilestoneStatusRepository {
	return &ProjectMilestoneStatusRepository{
		db: db,
	}
}

func (repo *ProjectMilestoneStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectMilestoneStatus, error) {

	var projectMilestoneStatuses []*models.ProjectMilestoneStatus

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
	result := query.Find(&projectMilestoneStatuses)

	return projectMilestoneStatuses, result.Error
}

func (repo *ProjectMilestoneStatusRepository) GetByID(id uint) (*models.ProjectMilestoneStatus, error) {
	var projectMilestoneStatus *models.ProjectMilestoneStatus
	result := repo.db.First(&projectMilestoneStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return projectMilestoneStatus, result.Error
}

func (repo *ProjectMilestoneStatusRepository) GetByType(status string) (*models.ProjectMilestoneStatus, error) {
	var projectMilestoneStatus *models.ProjectMilestoneStatus

	result := repo.db.Where("type = ?", status).First(&projectMilestoneStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Project Milestone Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return projectMilestoneStatus, nil
}
