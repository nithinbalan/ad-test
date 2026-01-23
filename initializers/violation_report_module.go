package initializers

import (
	"jethur/grc/violation-report/controllers"
	"jethur/grc/violation-report/repositories"
	"jethur/grc/violation-report/services"

	"gorm.io/gorm"
)

type ViolationReportModule struct {
	ViolationReportChannelController        *controllers.ViolationReportChannelController
	ViolationReportReasonController         *controllers.ViolationReportReasonController
	ViolationReportCategoryController       *controllers.ViolationReportCategoryController
	ViolationReportSubCategoryController    *controllers.ViolationReportSubCategoryController
	ViolationReportStatusController         *controllers.ViolationReportStatusController
	ViolationReportController               *controllers.ViolationReportController
	ViolationReportAssociatedItemController *controllers.ViolationReportAssociatedItemController
	ViolationReportInvestigationController  *controllers.ViolationReportInvestigationController
}

func NewViolationReportModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	complianceModule *ComplianceModule,
	riskModule *RiskModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	assetModule *AssetModule,
) *ViolationReportModule {

	violationReportChannelRepo := repositories.NewViolationReportChannelRepository(db)
	violationReportChannelSvc := services.NewViolationReportChannelService(violationReportChannelRepo, sharedModule.ActivityLogService)
	violationReportChannelController := controllers.NewViolationReportChannelController(violationReportChannelSvc)

	violationReportReasonRepo := repositories.NewViolationReportReasonRepository(db)
	violationReportReasonSvc := services.NewViolationReportReasonService(violationReportReasonRepo, sharedModule.ActivityLogService)
	violationReportReasonController := controllers.NewViolationReportReasonController(violationReportReasonSvc)

	violationReportCategoryRepo := repositories.NewViolationReportCategoryRepository(db)
	violationReportCategorySvc := services.NewViolationReportCategoryService(violationReportCategoryRepo, sharedModule.ActivityLogService, settingsModule.OrganizationUnitSvc)
	violationReportCategoryController := controllers.NewViolationReportCategoryController(violationReportCategorySvc)

	violationReportSubCategoryRepo := repositories.NewViolationReportSubCategoryRepository(db)
	violationReportSubCategorySvc := services.NewViolationReportSubCategoryService(violationReportSubCategoryRepo, sharedModule.ActivityLogService, violationReportCategorySvc)
	violationReportSubCategoryController := controllers.NewViolationReportSubCategoryController(violationReportSubCategorySvc)

	violationReportStatusRepo := repositories.NewViolationReportStatusRepository(db)
	violationReportStatusSvc := services.NewViolationReportStatusService(violationReportStatusRepo)
	violationReportStatusController := controllers.NewViolationReportStatusController(violationReportStatusSvc)

	violationReportRepo := repositories.NewViolationReportRepository(db)
	violationReportSvc := services.NewViolationReportService(violationReportRepo, sharedModule.ActivityLogService, violationReportChannelSvc, violationReportReasonSvc, violationReportCategorySvc, violationReportSubCategorySvc, violationReportStatusSvc, masterModule.SeverityLevelSvc, organizationModule.UserSvc, organizationModule.UserSettingSvc, commonModule.NotificationService, settingsModule.SystemPreferenceSvc)
	violationReportController := controllers.NewViolationReportController(violationReportSvc)

	violationReportAssociatedItemRepo := repositories.NewViolationReportAssociatedItemRepository(db)
	violationReportAssociatedItemSvc := services.NewViolationReportAssociatedItemService(violationReportAssociatedItemRepo, violationReportSvc, commonModule.AssociatedItemValidationService)
	violationReportAssociatedItemController := controllers.NewViolationReportAssociatedItemController(violationReportAssociatedItemSvc)

	violationReportInvestigationRepo := repositories.NewViolationReportInvestigationRepository(db)
	violationReportInvestigationSvc := services.NewViolationReportInvestigationService(
		violationReportInvestigationRepo,
		sharedModule.ActivityLogService,
		violationReportSvc,
		organizationModule.UserSvc,
	)
	violationReportInvestigationController := controllers.NewViolationReportInvestigationController(violationReportInvestigationSvc)

	return &ViolationReportModule{

		ViolationReportChannelController:        violationReportChannelController,
		ViolationReportReasonController:         violationReportReasonController,
		ViolationReportCategoryController:       violationReportCategoryController,
		ViolationReportSubCategoryController:    violationReportSubCategoryController,
		ViolationReportStatusController:         violationReportStatusController,
		ViolationReportController:               violationReportController,
		ViolationReportAssociatedItemController: violationReportAssociatedItemController,
		ViolationReportInvestigationController:  violationReportInvestigationController,
	}
}
