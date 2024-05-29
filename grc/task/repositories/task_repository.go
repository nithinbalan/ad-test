package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type TaskRepositoryInteface interface {
	CreateTask(task *models.Task, tx ...*gorm.DB) error
	GetModuleIDByIdentifier(identifier string) (uint, error)

}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (repo *TaskRepository) CreateTask(task *models.Task, tx ...*gorm.DB) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0]
	}
	return db.Create(&task).Error
}

func (repo *TaskRepository) GetModuleIDByIdentifier(identifier string) (uint, error) {
	var module models.Module
	if err := repo.db.Where("identifier = ?", identifier).First(&module).Error; err != nil {
		return 0, err
	}
	return module.ID, nil
}

