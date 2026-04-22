package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskRegisterUpdatePlanStatus, error)
	GetByID(id uint) (*models.RiskRegisterUpdatePlanStatus, error)
	GetByType(status string) (*models.RiskRegisterUpdatePlanStatus, error)
}

type RiskRegisterUpdatePlanStatusRepository struct {
	db *gorm.DB
}

func NewRiskRegisterUpdatePlanStatusRepository(db *gorm.DB) *RiskRegisterUpdatePlanStatusRepository {
	return &RiskRegisterUpdatePlanStatusRepository{
		db: db,
	}
}

func (repo *RiskRegisterUpdatePlanStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskRegisterUpdatePlanStatus, error) {
	var riskRegisterUpdatePlanStatuses []*models.RiskRegisterUpdatePlanStatus

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
	result := query.Find(&riskRegisterUpdatePlanStatuses)

	return riskRegisterUpdatePlanStatuses, result.Error
}

func (repo *RiskRegisterUpdatePlanStatusRepository) GetByID(id uint) (*models.RiskRegisterUpdatePlanStatus, error) {
	var riskRegisterUpdatePlanStatus *models.RiskRegisterUpdatePlanStatus

	result := repo.db.First(&riskRegisterUpdatePlanStatus, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskRegisterUpdatePlanStatus, nil
}

func (repo *RiskRegisterUpdatePlanStatusRepository) GetByType(status string) (*models.RiskRegisterUpdatePlanStatus, error) {
	var riskRegisterUpdatePlanStatus *models.RiskRegisterUpdatePlanStatus

	result := repo.db.Where("type = ?", status).First(&riskRegisterUpdatePlanStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Risk Register Update Plan Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskRegisterUpdatePlanStatus, nil
}
