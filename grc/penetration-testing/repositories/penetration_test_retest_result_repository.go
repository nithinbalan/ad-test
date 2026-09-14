package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestRetestResultRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestRetestResult, error)
	GetByID(id uint) (*models.PenetrationTestRetestResult, error)
}

type PenetrationTestRetestResultRepository struct{ db *gorm.DB }

func NewPenetrationTestRetestResultRepository(db *gorm.DB) *PenetrationTestRetestResultRepository {
	return &PenetrationTestRetestResultRepository{db: db}
}

func (repo *PenetrationTestRetestResultRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestRetestResult, error) {
	var values []*models.PenetrationTestRetestResult
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

func (repo *PenetrationTestRetestResultRepository) GetByID(id uint) (*models.PenetrationTestRetestResult, error) {
	var value *models.PenetrationTestRetestResult
	result := repo.db.First(&value, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}
