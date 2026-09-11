package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestStatus, error)
	GetByID(id uint) (*models.PenetrationTestStatus, error)
	GetByType(testType string) (*models.PenetrationTestStatus, error)
}

type PenetrationTestStatusRepository struct{ db *gorm.DB }

func NewPenetrationTestStatusRepository(db *gorm.DB) *PenetrationTestStatusRepository {
	return &PenetrationTestStatusRepository{db: db}
}

func (repo *PenetrationTestStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestStatus, error) {
	var values []*models.PenetrationTestStatus
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

func (repo *PenetrationTestStatusRepository) GetByID(id uint) (*models.PenetrationTestStatus, error) {
	var value *models.PenetrationTestStatus
	result := repo.db.First(&value, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}

func (repo *PenetrationTestStatusRepository) GetByType(testType string) (*models.PenetrationTestStatus, error) {
	var value *models.PenetrationTestStatus
	result := repo.db.Where("type = ?", testType).First(&value)
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}
