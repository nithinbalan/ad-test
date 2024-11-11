package initializers

import (
	"jethur/grc/logs/controllers"
	"jethur/grc/logs/repositories"
	"jethur/grc/logs/services"

	"gorm.io/gorm"
)

type SharedModule struct {
	ActivityLogController *controllers.ActivityLogController

	ActivityLogService services.ActivityLogServiceInterface
}

func NewSharedModule(db *gorm.DB) *SharedModule {
	activityLogRepo := repositories.NewActivityLogRepository(db)
	activityLogSvc := services.NewActivityLogService(activityLogRepo)
	activityLogController := controllers.NewActivityLogController(activityLogSvc)

	return &SharedModule{
		ActivityLogController: activityLogController,

		ActivityLogService: activityLogSvc,
	}
}
