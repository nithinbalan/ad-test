package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type RiskUpdatePlanTimelineTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.RiskUpdatePlanTimelineType, error)
	GetByID(id uint) (*models.RiskUpdatePlanTimelineType, error)
	GetByType(timelineType string) (*models.RiskUpdatePlanTimelineType, error)
}

// implementation
type RiskUpdatePlanTimelineTypeRepository struct {
	db *gorm.DB
}

func NewRiskUpdatePlanTimelineTypeRepository(db *gorm.DB) *RiskUpdatePlanTimelineTypeRepository {
	return &RiskUpdatePlanTimelineTypeRepository{
		db: db,
	}
}

func (repo *RiskUpdatePlanTimelineTypeRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.RiskUpdatePlanTimelineType, error) {

	var timelineTypes []*models.RiskUpdatePlanTimelineType

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

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&timelineTypes)

	return timelineTypes, result.Error
}

func (repo *RiskUpdatePlanTimelineTypeRepository) GetByID(id uint) (*models.RiskUpdatePlanTimelineType, error) {
	var timelineType *models.RiskUpdatePlanTimelineType
	result := repo.db.First(&timelineType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return timelineType, nil
}

func (repo *RiskUpdatePlanTimelineTypeRepository) GetByType(timelineType string) (*models.RiskUpdatePlanTimelineType, error) {
	var riskUpdatePlanTimelineType *models.RiskUpdatePlanTimelineType

	result := repo.db.Where("type = ?", timelineType).First(&riskUpdatePlanTimelineType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Risk Update Plan Timeline Type with type "+timelineType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskUpdatePlanTimelineType, nil
}
