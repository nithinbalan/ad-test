package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentImpactTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentImpactType, error)
	GetByID(id uint) (*models.ComplianceAssessmentImpactType, error)
	GetByType(status string) (*models.ComplianceAssessmentImpactType, error)
}

type ComplianceAssessmentImpactTypeRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentImpactTypeRepository(db *gorm.DB) *ComplianceAssessmentImpactTypeRepository {
	return &ComplianceAssessmentImpactTypeRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentImpactTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentImpactType, error) {
	var impactTypes []*models.ComplianceAssessmentImpactType
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ? OR description_en LIKE ? OR description_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
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

func (repo *ComplianceAssessmentImpactTypeRepository) GetByID(id uint) (*models.ComplianceAssessmentImpactType, error) {
	var impactType *models.ComplianceAssessmentImpactType
	result := repo.db.First(&impactType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return impactType, nil
}

func (repo *ComplianceAssessmentImpactTypeRepository) GetByType(status string) (*models.ComplianceAssessmentImpactType, error) {
	var assessmentStatus *models.ComplianceAssessmentImpactType
	result := repo.db.Where("type = ?", status).First(&assessmentStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Compliance Assessment Impact Type with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return assessmentStatus, nil
}
