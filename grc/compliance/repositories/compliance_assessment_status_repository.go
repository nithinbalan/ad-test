package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentStatus, error)
	GetByID(id uint) (*models.ComplianceAssessmentStatus, error)
	GetByType(status string) (*models.ComplianceAssessmentStatus, error)
}

type ComplianceAssessmentStatusRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentStatusRepository(db *gorm.DB) *ComplianceAssessmentStatusRepository {
	return &ComplianceAssessmentStatusRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentStatus, error) {
	var statuses []*models.ComplianceAssessmentStatus
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

func (repo *ComplianceAssessmentStatusRepository) GetByID(id uint) (*models.ComplianceAssessmentStatus, error) {
	var status *models.ComplianceAssessmentStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *ComplianceAssessmentStatusRepository) GetByType(status string) (*models.ComplianceAssessmentStatus, error) {
	var assessmentStatus *models.ComplianceAssessmentStatus
	result := repo.db.Where("type = ?", status).First(&assessmentStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Compliance Assessment Status with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return assessmentStatus, nil
}
