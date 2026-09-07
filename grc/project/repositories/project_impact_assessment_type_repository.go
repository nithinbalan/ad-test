package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectImpactAssessmentTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectImpactAssessmentType, error)
	GetByID(id uint) (*models.ProjectImpactAssessmentType, error)
	GetByType(status string) (*models.ProjectImpactAssessmentType, error)
}

type ProjectImpactAssessmentTypeRepository struct {
	db *gorm.DB
}

func NewProjectImpactAssessmentTypeRepository(db *gorm.DB) *ProjectImpactAssessmentTypeRepository {
	return &ProjectImpactAssessmentTypeRepository{
		db: db,
	}
}

func (repo *ProjectImpactAssessmentTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectImpactAssessmentType, error) {
	var impactTypes []*models.ProjectImpactAssessmentType
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
	result := query.Find(&impactTypes)

	return impactTypes, result.Error
}

func (repo *ProjectImpactAssessmentTypeRepository) GetByID(id uint) (*models.ProjectImpactAssessmentType, error) {
	var impactType *models.ProjectImpactAssessmentType
	result := repo.db.First(&impactType, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return impactType, result.Error
}

func (repo *ProjectImpactAssessmentTypeRepository) GetByType(status string) (*models.ProjectImpactAssessmentType, error) {
	var impactType *models.ProjectImpactAssessmentType

	result := repo.db.Where("type = ?", status).First(&impactType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Project Impact Assessment Type with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return impactType, nil
}
