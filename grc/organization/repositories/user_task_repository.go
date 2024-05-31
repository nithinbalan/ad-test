package repositories

import (
	"jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

type UserTaskRepositoryInteface interface {
	GetUserTaskCountByStatus(userID uint) ([]dtos.UserTaskRawCountDTO, error)
}

type UserTaskRepository struct {
	db *gorm.DB
}

func NewUserTaskRepository(db *gorm.DB) *UserTaskRepository {
	return &UserTaskRepository{
		db: db,
	}
}

func (repo *UserTaskRepository) GetUserTaskCountByStatus(userID uint) ([]dtos.UserTaskRawCountDTO, error) {
	var taskStatusCounts []dtos.UserTaskRawCountDTO

	query := `
		SELECT task_statuses.type AS Type,
		       COUNT(tasks.id) AS Count
		FROM tasks
		LEFT JOIN task_statuses ON task_statuses.id = tasks.task_status_id
		WHERE tasks.responsible_user_id = ?
		GROUP BY task_statuses.type
	`
	err := repo.db.Raw(query, userID).Scan(&taskStatusCounts).Error
	if err != nil {
		return nil, err
	}

	return taskStatusCounts, nil
}
