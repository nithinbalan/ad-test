package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestExecutionStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestExecutionStatus, error)
	GetByID(id uint) (*models.PenetrationTestExecutionStatus, error)
}

type PenetrationTestExecutionStatusRepository struct{ db *gorm.DB }

func NewPenetrationTestExecutionStatusRepository(db *gorm.DB) *PenetrationTestExecutionStatusRepository {
	return &PenetrationTestExecutionStatusRepository{db: db}
}

func (repo *PenetrationTestExecutionStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestExecutionStatus, error) {
	var values []*models.PenetrationTestExecutionStatus
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

func (repo *PenetrationTestExecutionStatusRepository) GetByID(id uint) (*models.PenetrationTestExecutionStatus, error) {
	var value *models.PenetrationTestExecutionStatus
	result := repo.db.First(&value, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}
