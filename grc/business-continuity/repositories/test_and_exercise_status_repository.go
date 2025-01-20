package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type TestAndExerciseStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TestAndExerciseStatus, error)
	GetByID(id uint) (*models.TestAndExerciseStatus, error)
	GetByType(status string) (*models.TestAndExerciseStatus, error)
}

// implementation
type TestAndExerciseStatusRepository struct {
	db *gorm.DB
}

func NewTestAndExerciseStatusRepository(db *gorm.DB) *TestAndExerciseStatusRepository {
	return &TestAndExerciseStatusRepository{
		db: db,
	}
}

func (repo *TestAndExerciseStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TestAndExerciseStatus, error) {
	var testAndExerciseStatuses []*models.TestAndExerciseStatus

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
	result := query.Find(&testAndExerciseStatuses)

	return testAndExerciseStatuses, result.Error
}

func (repo *TestAndExerciseStatusRepository) GetByID(id uint) (*models.TestAndExerciseStatus, error) {
	var status *models.TestAndExerciseStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *TestAndExerciseStatusRepository) GetByType(status string) (*models.TestAndExerciseStatus, error) {
	var testAndExerciseStatus *models.TestAndExerciseStatus

	result := repo.db.Where("type = ?", status).First(&testAndExerciseStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExerciseStatus, nil
}
