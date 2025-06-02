package initializers

import (
	"jethur/grc/logs/controllers"
	"jethur/grc/logs/repositories"
	"jethur/grc/logs/services"

	"gorm.io/gorm"
)

type SharedModule struct {
	ActivityLogController      *controllers.ActivityLogController
	TaskSchedularLogController *controllers.TaskSchedularLogController

	ActivityLogService services.ActivityLogServiceInterface
}

func NewSharedModule(db *gorm.DB) *SharedModule {
	activityLogRepo := repositories.NewActivityLogRepository(db)
	activityLogSvc := services.NewActivityLogService(activityLogRepo)
	activityLogController := controllers.NewActivityLogController(activityLogSvc)

	taskSchedularLogRepo := repositories.NewTaskSchedularLogRepository(db)
	taskSchedularLogController := controllers.NewTaskSchedularLogController(db, taskSchedularLogRepo)

	return &SharedModule{
		ActivityLogController:      activityLogController,
		TaskSchedularLogController: taskSchedularLogController,

		ActivityLogService: activityLogSvc,
	}
}
