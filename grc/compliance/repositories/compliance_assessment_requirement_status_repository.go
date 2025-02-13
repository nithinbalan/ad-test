package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentRequirementStatus, error)
	GetByID(id uint) (*models.ComplianceAssessmentRequirementStatus, error)
	GetByType(status string) (*models.ComplianceAssessmentRequirementStatus, error)
}

type ComplianceAssessmentRequirementStatusRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentRequirementStatusRepository(db *gorm.DB) *ComplianceAssessmentRequirementStatusRepository {
	return &ComplianceAssessmentRequirementStatusRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentRequirementStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentRequirementStatus, error) {
	var statuses []*models.ComplianceAssessmentRequirementStatus
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
	result := query.Find(&statuses)

	return statuses, result.Error
}

func (repo *ComplianceAssessmentRequirementStatusRepository) GetByID(id uint) (*models.ComplianceAssessmentRequirementStatus, error) {
	var status *models.ComplianceAssessmentRequirementStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *ComplianceAssessmentRequirementStatusRepository) GetByType(status string) (*models.ComplianceAssessmentRequirementStatus, error) {
	var requirementStatus *models.ComplianceAssessmentRequirementStatus
	result := repo.db.Where("type = ?", status).First(&requirementStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Compliance Assessment Requirement Status with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return requirementStatus, nil
}
