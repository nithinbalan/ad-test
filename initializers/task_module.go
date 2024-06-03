package initializers

import (
	"jethur/grc/task/controllers"
	"jethur/grc/task/repositories"
	"jethur/grc/task/services"

	"gorm.io/gorm"
)

type TaskModule struct {
	TaskStatusController *controllers.TaskStatusController
	TaskController       *controllers.TaskController

	TaskService services.TaskServiceInterface
}

func NewTaskModule(db *gorm.DB) *TaskModule {
	taskStatusRepo := repositories.NewTaskStatusRepository(db)
	taskStatusSvc := services.NewTaskStatusService(taskStatusRepo)
	taskStatusController := controllers.NewTaskStatusController(taskStatusSvc)

	taskRepo := repositories.NewTaskRepository(db)
	taskSvc := services.NewTaskService(taskRepo, taskStatusSvc)
	taskController := controllers.NewTaskController(taskSvc)

	return &TaskModule{
		TaskStatusController: taskStatusController,
		TaskController:       taskController,

		TaskService: taskSvc,
	}
}
