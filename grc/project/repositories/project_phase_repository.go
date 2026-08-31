package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ProjectPhaseRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectPhase, error)
	GetByID(id uint) (*models.ProjectPhase, error)
	GetByType(status string) (*models.ProjectPhase, error)
}

// implementation
type ProjectPhaseRepository struct {
	db *gorm.DB
}

func NewProjectPhaseRepository(db *gorm.DB) *ProjectPhaseRepository {
	return &ProjectPhaseRepository{
		db: db,
	}
}

func (repo *ProjectPhaseRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectPhase, error) {
	var projectPhases []*models.ProjectPhase
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
	result := query.Find(&projectPhases)

	return projectPhases, result.Error
}

func (repo *ProjectPhaseRepository) GetByID(id uint) (*models.ProjectPhase, error) {
	var phase *models.ProjectPhase
	result := repo.db.First(&phase, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return phase, result.Error
}

func (repo *ProjectPhaseRepository) GetByType(status string) (*models.ProjectPhase, error) {
	var projectPhase *models.ProjectPhase

	result := repo.db.Where("type = ?", status).First(&projectPhase)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Project Phase with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return projectPhase, nil
}