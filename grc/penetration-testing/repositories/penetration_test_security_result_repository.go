package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestSecurityResultRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestSecurityResult, error)
	GetByID(id uint) (*models.PenetrationTestSecurityResult, error)
}

type PenetrationTestSecurityResultRepository struct{ db *gorm.DB }

func NewPenetrationTestSecurityResultRepository(db *gorm.DB) *PenetrationTestSecurityResultRepository {
	return &PenetrationTestSecurityResultRepository{db: db}
}

func (repo *PenetrationTestSecurityResultRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestSecurityResult, error) {
	var values []*models.PenetrationTestSecurityResult
	query := repo.db
	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	return values, query.Order(sortBy + " " + sort).Find(&values).Error
}

func (repo *PenetrationTestSecurityResultRepository) GetByID(id uint) (*models.PenetrationTestSecurityResult, error) {
	var value *models.PenetrationTestSecurityResult
	result := repo.db.First(&value, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}
