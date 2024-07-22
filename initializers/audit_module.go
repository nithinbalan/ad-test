package initializers

import (
	"jethur/grc/audit/controllers"
	"jethur/grc/audit/repositories"
	"jethur/grc/audit/services"

	"gorm.io/gorm"
)

type AuditModule struct {
	AuditStatusController     *controllers.AuditStatusController
	AuditRiskRatingController *controllers.AuditRiskRatingController
}

func NewAuditModule(db *gorm.DB) *AuditModule {

	auditStatusRepo := repositories.NewAuditStatusRepository(db)
	auditStatusSvc := services.NewAuditStatusService(auditStatusRepo)
	auditStatusController := controllers.NewAuditStatusController(auditStatusSvc)

	auditRiskRatingRepo := repositories.NewAuditRiskRatingRepository(db)
	auditRiskRatingSvc := services.NewAuditRiskRatingService(auditRiskRatingRepo)
	auditRiskRatingController := controllers.NewAuditRiskRatingController(auditRiskRatingSvc)

	return &AuditModule{
		AuditStatusController:     auditStatusController,
		AuditRiskRatingController: auditRiskRatingController,
	}
}
