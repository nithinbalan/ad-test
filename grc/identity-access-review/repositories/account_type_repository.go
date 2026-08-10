package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AccountTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AccountType, error)
	GetByID(id uint) (*models.AccountType, error)
}

type AccountTypeRepository struct {
	db *gorm.DB
}

func NewAccountTypeRepository(db *gorm.DB) *AccountTypeRepository {
	return &AccountTypeRepository{db: db}
}

func (repo *AccountTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AccountType, error) {
	var accountTypes []*models.AccountType
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

	result := query.Order(sortBy + " " + sort).Find(&accountTypes)
	return accountTypes, result.Error
}

func (repo *AccountTypeRepository) GetByID(id uint) (*models.AccountType, error) {
	var accountType models.AccountType
	if err := repo.db.First(&accountType, id).Error; err != nil {
		return nil, err
	}
	return &accountType, nil
}
