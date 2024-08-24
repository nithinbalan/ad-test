package initializers

import (
	"jethur/grc/audit/controllers"
	"jethur/grc/audit/repositories"
	"jethur/grc/audit/services"

	"gorm.io/gorm"
)

type AuditModule struct {
	AuditStatusController                    *controllers.AuditStatusController
	AuditorRiskRatingController              *controllers.AuditorRiskRatingController
	AuditUniverseController                  *controllers.AuditUniverseController
	AuditableItemController                  *controllers.AuditableItemController
	AuditTeamController                      *controllers.AuditTeamController
	AuditCategoryController                  *controllers.AuditCategoryController
	AuditPlanController                      *controllers.AuditPlanController
	AuditPlanChecklistController             *controllers.AuditPlanChecklistController
	AuditPlanToolController                  *controllers.AuditPlanToolController
	AuditPlanAuditableItemController         *controllers.AuditPlanAuditableItemController
	AuditScheduleStatusController            *controllers.AuditScheduleStatusController
	AuditActivityTypeController              *controllers.AuditActivityTypeController
	AuditPlanBudgetController                *controllers.AuditPlanBudgetController
	AuditPlanScheduleController              *controllers.AuditPlanScheduleController
	AuditEngagementTypeController            *controllers.AuditEngagementTypeController
	AuditPlanStakeholderEngagementController *controllers.AuditPlanStakeholderEngagementController
	AuditController                          *controllers.AuditController
}

func NewAuditModule(db *gorm.DB, sharedModule *SharedModule, organizationModule *OrganizationModule, complianceModule *ComplianceModule, riskModule *RiskModule, incidentModule *IncidentModule, documentModule *DocumentModule, settingModule *SettingsModule, masterModule *MasterModule) *AuditModule {

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
		documentModule.DocumentSvc,
	)
	auditUniverseController := controllers.NewAuditUniverseController(auditUniverseSvc)

	auditableItemRepo := repositories.NewAuditableItemRepository(db)
	auditableItemSvc := services.NewAuditableItemService(
		auditableItemRepo,
		sharedModule.ActivityLogService,
		auditStatusSvc, auditorRiskRatingSvc,
		settingModule.OrganizationUnitSvc,
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

	auditPlanChecklistRepo := repositories.NewAuditPlanChecklistRepository(db)
	auditPlanChecklistSvc := services.NewAuditPlanChecklistService(
		auditPlanChecklistRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		auditableItemSvc,
	)
	auditPlanChecklistController := controllers.NewAuditPlanChecklistController(auditPlanChecklistSvc)

	auditPlanToolRepo := repositories.NewAuditPlanToolRepository(db)
	auditPlanToolSvc := services.NewAuditPlanToolService(
		auditPlanToolRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		masterModule.ToolTypeSvc,
	)
	auditPlanToolController := controllers.NewAuditPlanToolController(auditPlanToolSvc)

	auditPlanAuditableItemRepo := repositories.NewAuditPlanAuditableItemRepository(db)
	auditPlanAuditableItemSvc := services.NewAuditPlanAuditableItemService(auditPlanAuditableItemRepo, sharedModule.ActivityLogService, auditPlanSvc, auditableItemSvc)
	auditPlanAuditableItemController := controllers.NewAuditPlanAuditableItemController(auditPlanAuditableItemSvc)

	auditPlanBudgetRepo := repositories.NewAuditPlanBudgetRepository(db)
	auditPlanBudgetSvc := services.NewAuditPlanBudgetService(
		auditPlanBudgetRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		masterModule.ExpenseCategorySvc,
	)
	auditPlanBudgetController := controllers.NewAuditPlanBudgetController(auditPlanBudgetSvc)

	auditActivityTypeRepo := repositories.NewAuditActivityTypeRepository(db)
	auditActivityTypeSvc := services.NewAuditActivityTypeService(auditActivityTypeRepo, sharedModule.ActivityLogService)
	auditActivityTypeController := controllers.NewAuditActivityTypeController(auditActivityTypeSvc)

	auditScheduleStatusRepo := repositories.NewAuditScheduleStatusRepository(db)
	auditScheduleStatusSvc := services.NewAuditScheduleStatusService(auditScheduleStatusRepo)
	auditScheduleStatusController := controllers.NewAuditScheduleStatusController(auditScheduleStatusSvc)

	auditPlanScheduleRepo := repositories.NewAuditPlanScheduleRepository(db)
	auditPlanScheduleSvc := services.NewAuditPlanScheduleService(
		auditPlanScheduleRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		auditActivityTypeSvc,
		organizationModule.UserSvc,
		auditScheduleStatusSvc,
		auditableItemSvc,
	)
	auditPlanScheduleController := controllers.NewAuditPlanScheduleController(auditPlanScheduleSvc)

	auditEngagementTypeRepo := repositories.NewAuditEngagementTypeRepository(db)
	auditEngagementTypeSvc := services.NewAuditEngagementTypeService(auditEngagementTypeRepo, sharedModule.ActivityLogService)
	auditEngagementTypeController := controllers.NewAuditEngagementTypeController(auditEngagementTypeSvc)

	auditPlanStakeholderEngagementRepo := repositories.NewAuditPlanStakeholderEngagementRepository(db)
	auditPlanStakeholderEngagementSvc := services.NewAuditPlanStakeholderEngagementService(
		auditPlanStakeholderEngagementRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
		auditEngagementTypeSvc,
		organizationModule.UserSvc,
		auditScheduleStatusSvc,
	)
	auditPlanStakeholderEngagementController := controllers.NewAuditPlanStakeholderEngagementController(auditPlanStakeholderEngagementSvc)

	auditRepo := repositories.NewAuditRepository(db)
	auditSvc := services.NewAuditService(
		auditRepo,
		sharedModule.ActivityLogService,
		auditPlanSvc,
	)
	auditController := controllers.NewAuditController(auditSvc)

	return &AuditModule{
		AuditorRiskRatingController:              auditorRiskRatingController,
		AuditUniverseController:                  auditUniverseController,
		AuditableItemController:                  auditableItemController,
		AuditTeamController:                      auditTeamController,
		AuditCategoryController:                  auditCategoryController,
		AuditStatusController:                    auditStatusController,
		AuditPlanController:                      auditPlanController,
		AuditPlanChecklistController:             auditPlanChecklistController,
		AuditPlanToolController:                  auditPlanToolController,
		AuditPlanAuditableItemController:         auditPlanAuditableItemController,
		AuditActivityTypeController:              auditActivityTypeController,
		AuditScheduleStatusController:            auditScheduleStatusController,
		AuditPlanBudgetController:                auditPlanBudgetController,
		AuditPlanScheduleController:              auditPlanScheduleController,
		AuditEngagementTypeController:            auditEngagementTypeController,
		AuditPlanStakeholderEngagementController: auditPlanStakeholderEngagementController,
		AuditController:                          auditController,
	}
}
