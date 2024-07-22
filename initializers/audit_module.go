package initializers

import (
	"jethur/grc/audit/controllers"
	"jethur/grc/audit/repositories"
	"jethur/grc/audit/services"

	"gorm.io/gorm"
)

type AuditModule struct {
	AuditStatusController       *controllers.AuditStatusController
	AuditorRiskRatingController *controllers.AuditorRiskRatingController
	AuditUniverseController     *controllers.AuditUniverseController
}

func NewAuditModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, complianceModule *ComplianceModule, riskModule *RiskModule, incidentModule *IncidentModule, documentModule *DocumentModule) *AuditModule {

	auditStatusRepo := repositories.NewAuditStatusRepository(db)
	auditStatusSvc := services.NewAuditStatusService(auditStatusRepo)
	auditStatusController := controllers.NewAuditStatusController(auditStatusSvc)

	auditorRiskRatingRepo := repositories.NewAuditorRiskRatingRepository(db)
	auditorRiskRatingSvc := services.NewAuditorRiskRatingService(auditorRiskRatingRepo)
	auditorRiskRatingController := controllers.NewAuditorRiskRatingController(auditorRiskRatingSvc)

	auditUniverseRepo := repositories.NewAuditUniverseRepository(db)
	auditUniverseSvc := services.NewAuditUniverseService(
		auditUniverseRepo,
		sharedModule.ActivityLogService,
		auditStatusSvc, auditorRiskRatingSvc,
		organizationModule.ProcessSvc,
		organizationModule.ControlSvc,
		organizationModule.UserSvc,
		complianceModule.ComplianceRequirementSvc,
		riskModule.RiskRegisterSvc,
		incidentModule.IncidentRegisterSvc,
		// documentModule.DocumentSvc,
	)
	auditUniverseController := controllers.NewAuditUniverseController(auditUniverseSvc)

	return &AuditModule{
		AuditStatusController:       auditStatusController,
		AuditorRiskRatingController: auditorRiskRatingController,
		AuditUniverseController:     auditUniverseController,
	}
}
