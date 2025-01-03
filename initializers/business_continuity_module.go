package initializers

import (
	"jethur/grc/business-continuity/controllers"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/business-continuity/services"

	"gorm.io/gorm"
)

type BusinessContinuityModule struct {
	BusinessFunctionController         *controllers.BusinessFunctionController
	BusinessImpactAssessmentController *controllers.BusinessImpactAssessmentController
	BusinessContinuityPlanController   *controllers.BusinessContinuityPlanController

	BusinessFunctionSvc services.BusinessFunctionServiceInterface
}

func NewBusinessContinuityModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	assetModule *AssetModule,
	documentModule *DocumentModule,
	riskModule *RiskModule,
) *BusinessContinuityModule {

	businessFunctionRepo := repositories.NewBusinessFunctionRepository(db)

	businessFunctionSvc := services.NewBusinessFunctionService(businessFunctionRepo,
		sharedModule.ActivityLogService,
		organizationModule.UserSvc,
	)
	businessFunctionController := controllers.NewBusinessFunctionController(businessFunctionSvc)

	businessImpactAssessmentRepo := repositories.NewBusinessImpactAssessmentRepository(db)
	businessImpactAssessmentSvc := services.NewBusinessImpactAssessmentService(businessImpactAssessmentRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, masterModule.ImpactCategorySvc, organizationModule.ProcessSvc, assetModule.AssetRegisterSvc, documentModule.DocumentSvc, riskModule.RiskLevelSvc)
	businessImpactAssessmentController := controllers.NewBusinessImpactAssessmentController(businessImpactAssessmentSvc, organizationModule.ProcessAssociatedItemSvc, organizationModule.ProcessSvc)

	businessContinuityPlanRepo := repositories.NewBusinessContinuityPlanRepository(db)
	businessContinuityPlanSvc := services.NewBusinessContinuityPlanService(businessContinuityPlanRepo, sharedModule.ActivityLogService, organizationModule.UserSvc, businessFunctionSvc)
	businessContinuityPlanController := controllers.NewBusinessContinuityPlanController(businessContinuityPlanSvc)

	return &BusinessContinuityModule{
		BusinessFunctionController:         businessFunctionController,
		BusinessImpactAssessmentController: businessImpactAssessmentController,
		BusinessContinuityPlanController:   businessContinuityPlanController,

		BusinessFunctionSvc: businessFunctionSvc,
	}
}
