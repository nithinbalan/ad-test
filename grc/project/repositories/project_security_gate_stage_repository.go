package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateStageRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateStage, error)
	GetByID(id uint) (*models.ProjectSecurityGateStage, error)
	GetByType(stageType string) (*models.ProjectSecurityGateStage, error)
}

type ProjectSecurityGateStageRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateStageRepository(db *gorm.DB) *ProjectSecurityGateStageRepository {
	return &ProjectSecurityGateStageRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateStageRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectSecurityGateStage, error) {
	var stages []*models.ProjectSecurityGateStage
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
	result := query.Find(&stages)

	return stages, result.Error
}

func (repo *ProjectSecurityGateStageRepository) GetByID(id uint) (*models.ProjectSecurityGateStage, error) {
	var stage *models.ProjectSecurityGateStage
	result := repo.db.First(&stage, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return stage, result.Error
}

func (repo *ProjectSecurityGateStageRepository) GetByType(stageType string) (*models.ProjectSecurityGateStage, error) {
	var stage *models.ProjectSecurityGateStage

	result := repo.db.Where("type = ?", stageType).First(&stage)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Project Security Gate Stage with type "+stageType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return stage, nil
}
