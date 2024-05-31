package repositories

import (
	taskDTOs "jethur/grc/task/dtos"

	"gorm.io/gorm"
)

type UserTaskRepositoryInteface interface {
	GetUserTaskCountByStatus(userID uint) ([]taskDTOs.UserTaskStatusCountDTO, error)
}

type UserTaskRepository struct {
	db *gorm.DB
}

func NewUserTaskRepository(db *gorm.DB) *UserTaskRepository {
	return &UserTaskRepository{
		db: db,
	}
}

func (repo *UserTaskRepository) GetUserTaskCountByStatus(userID uint) ([]taskDTOs.UserTaskStatusCountDTO, error) {
	var taskStatusCounts []taskDTOs.UserTaskStatusCountDTO

	query := `
		SELECT task_statuses.title_en AS status, COUNT(tasks.id) AS count
		FROM tasks
		LEFT JOIN task_statuses ON task_statuses.id = tasks.task_status_id
		WHERE tasks.responsible_user_id = ?
		GROUP BY task_statuses.title_en
	`
	err := repo.db.Raw(query, userID).Scan(&taskStatusCounts).Error
	if err != nil {
		return nil, err
	}

	return taskStatusCounts, nil
}
