package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentFrameworkTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentFrameworkType, error)
	GetByID(id uint) (*models.ComplianceAssessmentFrameworkType, error)
	GetByType(typeValue string) (*models.ComplianceAssessmentFrameworkType, error)
}

type ComplianceAssessmentFrameworkTypeRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentFrameworkTypeRepository(db *gorm.DB) *ComplianceAssessmentFrameworkTypeRepository {
	return &ComplianceAssessmentFrameworkTypeRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentFrameworkTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentFrameworkType, error) {
	var frameworkTypes []*models.ComplianceAssessmentFrameworkType
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
	result := query.Find(&frameworkTypes)

	return frameworkTypes, result.Error
}

func (repo *ComplianceAssessmentFrameworkTypeRepository) GetByID(id uint) (*models.ComplianceAssessmentFrameworkType, error) {
	var frameworkType *models.ComplianceAssessmentFrameworkType
	result := repo.db.First(&frameworkType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return frameworkType, nil
}

func (repo *ComplianceAssessmentFrameworkTypeRepository) GetByType(typeValue string) (*models.ComplianceAssessmentFrameworkType, error) {
	var result *models.ComplianceAssessmentFrameworkType
	err := repo.db.Where("type = ?", typeValue).First(&result).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(err, "Compliance Assessment Framework Type with type "+typeValue)
	} else if err != nil {
		return nil, err
	}

	return result, nil
}
