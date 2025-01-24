package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type TestResultOptionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TestResultOption, error)
	GetByID(id uint) (*models.TestResultOption, error)
	GetByType(status string) (*models.TestResultOption, error)
}

// implementation
type TestResultOptionRepository struct {
	db *gorm.DB
}

func NewTestResultOptionRepository(db *gorm.DB) *TestResultOptionRepository {
	return &TestResultOptionRepository{
		db: db,
	}
}

func (repo *TestResultOptionRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TestResultOption, error) {
	var testResultOption []*models.TestResultOption

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
	result := query.Find(&testResultOption)

	return testResultOption, result.Error
}

func (repo *TestResultOptionRepository) GetByID(id uint) (*models.TestResultOption, error) {
	var status *models.TestResultOption
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *TestResultOptionRepository) GetByType(status string) (*models.TestResultOption, error) {
	var testResultOption *models.TestResultOption

	result := repo.db.Where("type = ?", status).First(&testResultOption)

	if result.Error != nil {
		return nil, result.Error
	}

	return testResultOption, nil
}
