package repositories

import (
	"jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type UserTaskRepositoryInteface interface {
	GetUserTaskCountByStatus(userID uint) ([]dtos.UserTaskStatusCountResponseDTO, error)
}

type UserTaskRepository struct {
	db *gorm.DB
}

func NewUserTaskRepository(db *gorm.DB) *UserTaskRepository {
	return &UserTaskRepository{
		db: db,
	}
}

func (repo *UserTaskRepository) GetUserTaskCountByStatus(userID uint) ([]dtos.UserTaskStatusCountResponseDTO, error) {
	var taskStatusCounts []dtos.UserTaskStatusCountResponseDTO

	query := `
        SELECT task_statuses.id, task_statuses.type, task_statuses.title_en, task_statuses.title_ar, task_statuses.color_code,
               COUNT(tasks.id) AS Count
        FROM tasks
        LEFT JOIN task_statuses ON task_statuses.id = tasks.task_status_id
        WHERE tasks.responsible_user_id = ?
        GROUP BY task_statuses.id, task_statuses.type, task_statuses.title_en, task_statuses.title_ar, task_statuses.color_code
    `
	err := repo.db.Raw(query, userID).Scan(&taskStatusCounts).Error
	if err != nil {
		return nil, err
	}

	return taskStatusCounts, nil
}
