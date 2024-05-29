package initializers

import (
	"jethur/grc/task/controllers"
	"jethur/grc/task/repositories"
	"jethur/grc/task/services"

	"gorm.io/gorm"
)

type TaskModule struct {
	TaskStatusController *controllers.TaskStatusController
	TaskService          services.TaskServiceInterface
}

func NewTaskModule(db *gorm.DB) *TaskModule {
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskStatusSvc := services.NewTaskStatusService(taskStatusRepo)
	taskStatusController := controllers.NewTaskStatusController(taskStatusSvc)

	taskRepo := repositories.NewTaskRepository(db)
	taskSvc := services.NewTaskService(taskRepo,taskStatusSvc)

	return &TaskModule{
		TaskStatusController: taskStatusController,
		TaskService:          taskSvc,
	}
}
