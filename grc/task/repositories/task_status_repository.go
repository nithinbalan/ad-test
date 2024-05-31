package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type TaskStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TaskStatus, error)
	GetByID(id uint) (*models.TaskStatus, error)
	GetByType(status string) (*models.TaskStatus, error)
}

// implementation
type TaskStatusRepository struct {
	db *gorm.DB
}

func NewTaskStatusRepository(db *gorm.DB) *TaskStatusRepository {
	return &TaskStatusRepository{
		db: db,
	}
}

func (repo *TaskStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TaskStatus, error) {

	var taskStatuses []*models.TaskStatus

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
	result := query.Find(&taskStatuses)

	return taskStatuses, result.Error
}

func (repo *TaskStatusRepository) GetByID(id uint) (*models.TaskStatus, error) {
	var status *models.TaskStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *TaskStatusRepository) GetByType(status string) (*models.TaskStatus, error) {
	var taskStatus *models.TaskStatus

	result := repo.db.Where("type = ?", status).First(&taskStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Task Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return taskStatus, nil
}
