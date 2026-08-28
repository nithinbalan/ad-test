package repositories

import (
	"gorm.io/gorm"
	"jethur/grc/models"
)

type PenetrationTestExecutionStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestExecutionStatus, error)
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
