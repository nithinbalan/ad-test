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
	AuditableItemController     *controllers.AuditableItemController
	AuditTeamController         *controllers.AuditTeamController
	AuditCategoryController     *controllers.AuditCategoryController
	AuditPlanController         *controllers.AuditPlanController
}

func NewAuditModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, complianceModule *ComplianceModule, riskModule *RiskModule, incidentModule *IncidentModule, documentModule *DocumentModule, settingModule *SettingsModule) *AuditModule {

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

	auditableItemRepo := repositories.NewAuditableItemRepository(db)
	auditableItemSvc := services.NewAuditableItemService(
		auditableItemRepo,
		sharedModule.ActivityLogService,
		auditStatusSvc, auditorRiskRatingSvc,
	)
	auditableItemController := controllers.NewAuditableItemController(auditableItemSvc)

	auditTeamRepo := repositories.NewAuditTeamRepository(db)
	auditTeamSvc := services.NewAuditTeamService(
		auditTeamRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
	)
	auditTeamController := controllers.NewAuditTeamController(auditTeamSvc)

	auditCategoryRepo := repositories.NewAuditCategoryRepository(db)
	auditCategorySvc := services.NewAuditCategoryService(auditCategoryRepo, sharedModule.ActivityLogService)
	auditCategoryController := controllers.NewAuditCategoryController(auditCategorySvc)

	auditPlanRepo := repositories.NewAuditPlanRepository(db)
	auditPlanSvc := services.NewAuditPlanService(
		auditPlanRepo,
		sharedModule.ActivityLogService,
		auditCategorySvc,
		auditStatusSvc,
		auditTeamSvc,
		organizationModule.UserSvc,
		auditableItemSvc,
	)
	auditPlanController := controllers.NewAuditPlanController(auditPlanSvc)

	return &AuditModule{
		AuditStatusController:       auditStatusController,
		AuditorRiskRatingController: auditorRiskRatingController,
		AuditUniverseController:     auditUniverseController,
		AuditableItemController:     auditableItemController,
		AuditTeamController:         auditTeamController,
		AuditCategoryController:     auditCategoryController,
		AuditPlanController:         auditPlanController,
	}
}
