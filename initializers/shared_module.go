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
	EmailLogController         *controllers.EmailLogController
	APIRequestLogController    *controllers.APIRequestLogController

	ActivityLogService services.ActivityLogServiceInterface
}

func NewSharedModule(db *gorm.DB) *SharedModule {
	activityLogRepo := repositories.NewActivityLogRepository(db)
	activityLogSvc := services.NewActivityLogService(activityLogRepo)
	activityLogController := controllers.NewActivityLogController(activityLogSvc)

	taskSchedularLogRepo := repositories.NewTaskSchedularLogRepository(db)
	taskSchedularLogController := controllers.NewTaskSchedularLogController(db, taskSchedularLogRepo)

	emailLogRepo := repositories.NewEmailLogRepository(db)
	emailLogSvc := services.NewEmailLogService(emailLogRepo)
	emailLogController := controllers.NewEmailLogController(emailLogSvc)

	apiRequestLogRepo := repositories.NewAPIRequestLogRepository(db)
	apiRequestLogSvc := services.NewAPIRequestLogService(apiRequestLogRepo)
	apiRequestLogController := controllers.NewAPIRequestLogController(db, apiRequestLogSvc)

	return &SharedModule{
		ActivityLogController:      activityLogController,
		TaskSchedularLogController: taskSchedularLogController,
		EmailLogController:         emailLogController,
		APIRequestLogController:    apiRequestLogController,

		ActivityLogService: activityLogSvc,
	}
}
