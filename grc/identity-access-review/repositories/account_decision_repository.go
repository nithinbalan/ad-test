package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AccountDecisionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AccountDecision, error)
	GetByID(id uint) (*models.AccountDecision, error)
}

type AccountDecisionRepository struct {
	db *gorm.DB
}

func NewAccountDecisionRepository(db *gorm.DB) *AccountDecisionRepository {
	return &AccountDecisionRepository{db: db}
}

func (repo *AccountDecisionRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AccountDecision, error) {
	var accountDecisions []*models.AccountDecision
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	if sortBy != "id" && sortBy != "title_en" && sortBy != "title_ar" && sortBy != "type" && sortBy != "status" {
		sortBy = "title_en"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	result := query.Order(sortBy + " " + sort).Find(&accountDecisions)
	return accountDecisions, result.Error
}

func (repo *AccountDecisionRepository) GetByID(id uint) (*models.AccountDecision, error) {
	var accountDecision models.AccountDecision
	if err := repo.db.First(&accountDecision, id).Error; err != nil {
		return nil, err
	}
	return &accountDecision, nil
}
