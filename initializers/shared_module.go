package initializers

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"

	"gorm.io/gorm"
)

type SharedModule struct {
	ActivityLogService logServices.ActivityLogServiceInterface
}

func NewSharedModule(db *gorm.DB) *SharedModule {
	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	return &SharedModule{
		ActivityLogService: activityLogSvc,
	}
}
