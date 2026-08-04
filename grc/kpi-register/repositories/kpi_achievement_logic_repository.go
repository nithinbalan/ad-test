package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type KPIAchievementLogicRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIAchievementLogic, error)
	GetByID(id uint) (*models.KPIAchievementLogic, error)
	GetByType(logicType string) (*models.KPIAchievementLogic, error)
}

type KPIAchievementLogicRepository struct {
	db *gorm.DB
}

func NewKPIAchievementLogicRepository(db *gorm.DB) *KPIAchievementLogicRepository {
	return &KPIAchievementLogicRepository{
		db: db,
	}
}

func (repo *KPIAchievementLogicRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIAchievementLogic, error) {

	var kpiAchievementLogic []*models.KPIAchievementLogic

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
	result := query.Find(&kpiAchievementLogic)

	return kpiAchievementLogic, result.Error
}

func (repo *KPIAchievementLogicRepository) GetByID(id uint) (*models.KPIAchievementLogic, error) {
	var kpiAchievementLogic *models.KPIAchievementLogic
	result := repo.db.First(&kpiAchievementLogic, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return kpiAchievementLogic, result.Error
}

func (repo *KPIAchievementLogicRepository) GetByType(logicType string) (*models.KPIAchievementLogic, error) {
	var kpiAchievementLogic *models.KPIAchievementLogic

	result := repo.db.Where("type = ?", logicType).First(&kpiAchievementLogic)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "KPI Achievement Logic with type "+logicType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return kpiAchievementLogic, nil
}
