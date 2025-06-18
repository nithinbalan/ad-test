package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ThirdPartyAssessmentStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyAssessmentStatus, error)
	GetByID(id uint) (*models.ThirdPartyAssessmentStatus, error)
	GetByType(status string) (*models.ThirdPartyAssessmentStatus, error)
}

// implementation
type ThirdPartyAssessmentStatusRepository struct {
	db *gorm.DB
}

func NewThirdPartyAssessmentStatusRepository(db *gorm.DB) *ThirdPartyAssessmentStatusRepository {
	return &ThirdPartyAssessmentStatusRepository{
		db: db,
	}
}

func (repo *ThirdPartyAssessmentStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyAssessmentStatus, error) {

	var thirdPartyAssessmentStatuses []*models.ThirdPartyAssessmentStatus

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
	result := query.Find(&thirdPartyAssessmentStatuses)

	return thirdPartyAssessmentStatuses, result.Error
}

func (repo *ThirdPartyAssessmentStatusRepository) GetByID(id uint) (*models.ThirdPartyAssessmentStatus, error) {
	var thirdPartyAssessmentStatus *models.ThirdPartyAssessmentStatus
	result := repo.db.First(&thirdPartyAssessmentStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return thirdPartyAssessmentStatus, result.Error
}

func (repo *ThirdPartyAssessmentStatusRepository) GetByType(status string) (*models.ThirdPartyAssessmentStatus, error) {
	var thirdPartyAssessmentStatus *models.ThirdPartyAssessmentStatus

	result := repo.db.Where("type = ?", status).First(&thirdPartyAssessmentStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "ThirdPartyAssessment Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return thirdPartyAssessmentStatus, nil
}
