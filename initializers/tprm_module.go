package initializers

import (
	"jethur/grc/tprm/controllers"
	"jethur/grc/tprm/repositories"
	"jethur/grc/tprm/services"

	"gorm.io/gorm"
)

type TprmModule struct {
	ThirdPartyTypeController                      *controllers.ThirdPartyTypeController
	ThirdPartyRegisterController                  *controllers.ThirdPartyRegisterController
	ThirdPartyQuestionCategoryController          *controllers.ThirdPartyQuestionCategoryController
	ThirdPartyQuestionTypeController              *controllers.ThirdPartyQuestionTypeController
	ThirdPartyQuestionLibraryController           *controllers.ThirdPartyQuestionLibraryController
	QuestionSetController                         *controllers.QuestionSetController
	ThirdPartyAssessmentStatusController          *controllers.ThirdPartyAssessmentStatusController
	ThirdPartyAssessmentController                *controllers.ThirdPartyAssessmentController
	ThirdPartyAssociatedItemController            *controllers.ThirdPartyAssociatedItemController
	TPRMSettingController                         *controllers.TPRMSettingController
	ThirdPartyAssessmentWorkflowController        *controllers.ThirdPartyAssessmentWorkflowController
	ThirdPartyAssessmentReviewController          *controllers.ThirdPartyAssessmentReviewController
	TPRMAssesmentSectionQuestionCommentController *controllers.TprmAssessmentSectionQuestionCommentController
	TPRMComplianeResponseController               *controllers.TPRMComplianceResponseController
	TPRMDashboardController                       *controllers.TPRMDashboardController

	ThirdPartyTypeSvc               services.ThirdPartyTypeServiceInterface
	ThirdPartyRegisterSvc           services.ThirdPartyRegisterServiceInterface
	ThirdPartyQuestionCategorySvc   services.ThirdPartyQuestionCategoryServiceInterface
	ThirdPartyQuestionTypeSvc       services.ThirdPartyQuestionTypeServiceInterface
	TPRMSettingSvc                  services.TPRMSettingServiceInterface
	ThirdPartyAssessmentWorkflowSvc services.ThirdPartyAssessmentWorkflowServiceInterface
	TPRMAssesmentSectionQuestionSvc services.TprmAssessmentSectionQuestionCommentServiceInterface
	TPRMComplianceResponseSvc       services.TPRMComplianceResponseServiceInterface
	ThirdPartyAssessmentSvc         services.ThirdPartyAssessmentServiceInterface
}

func NewTprmModule(
	db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	riskModule *RiskModule,
) *TprmModule {

	// === Third Party Type ===
	thirdPartyTypeRepo := repositories.NewThirdPartyTypeRepository(db)
	thirdPartyTypeSvc := services.NewThirdPartyTypeService(thirdPartyTypeRepo, sharedModule.ActivityLogService)
	thirdPartyTypeController := controllers.NewThirdPartyTypeController(thirdPartyTypeSvc)

	tprmSettingRepo := repositories.NewTPRMSettingRepository(db)
	tprmSettingSvc := services.NewTPRMSettingService(tprmSettingRepo)
	tprmSettingController := controllers.NewTPRMSettingController(tprmSettingSvc)

	// === Third Party Register ===
	thirdPartyRegisterRepo := repositories.NewThirdPartyRegisterRepository(db)
	thirdPartyRegisterSvc := services.NewThirdPartyRegisterService(
		thirdPartyRegisterRepo,
		sharedModule.ActivityLogService,
		masterModule.SeverityLevelSvc,
		organizationModule.UserSvc,
		settingsModule.OrganizationUnitSvc,
		thirdPartyTypeSvc,
		masterModule.CountrySvc,
		organizationModule.UserSettingSvc,
		commonModule.NotificationService,
	)
	thirdPartyRegisterController := controllers.NewThirdPartyRegisterController(thirdPartyRegisterSvc)

	thirdPartyQuestionCategoryRepo := repositories.NewThirdPartyQuestionCategoryRepository(db)
	thirdPartyQuestionCategorySvc := services.NewThirdPartyQuestionCategoryService(thirdPartyQuestionCategoryRepo, sharedModule.ActivityLogService)
	thirdPartyQuestionCategoryController := controllers.NewThirdPartyQuestionCategoryController(thirdPartyQuestionCategorySvc)

	thirdPartyQuestionTypeRepo := repositories.NewThirdPartyQuestionTypeRepository(db)
	thirdPartyQuestionTypeSvc := services.NewThirdPartyQuestionTypeService(thirdPartyQuestionTypeRepo)
	thirdPartyQuestionTypeController := controllers.NewThirdPartyQuestionTypeController(thirdPartyQuestionTypeSvc)

	tprmComplianceResponseRepository := repositories.NewTPRMComplianceResponseRepository(db)
	tprmComplianceResponseSvc := services.NewTPRMComplianceResponseService(tprmComplianceResponseRepository,
		sharedModule.ActivityLogService)
	tprmComplianceResponseController := controllers.NewTPRMComplianceResponseController(tprmComplianceResponseSvc)

	// === Third Party Question Library ===
	thirdPartyQuestionLibraryRepo := repositories.NewThirdPartyQuestionLibraryRepository(db)
	thirdPartyQuestionLibrarySvc := services.NewThirdPartyQuestionLibraryService(
		thirdPartyQuestionLibraryRepo,
		sharedModule.ActivityLogService,
		thirdPartyQuestionCategorySvc,
		thirdPartyQuestionTypeSvc,
		tprmSettingSvc,
		riskModule.RiskLevelSvc,
		tprmComplianceResponseSvc,
		organizationModule.UserSettingSvc,
	)
	thirdPartyQuestionLibraryController := controllers.NewThirdPartyQuestionLibraryController(thirdPartyQuestionLibrarySvc)

	questionSetRepo := repositories.NewQuestionSetRepository(db)
	questionSetSvc := services.NewQuestionSetService(
		questionSetRepo,
		sharedModule.ActivityLogService,
		thirdPartyQuestionLibrarySvc,
		organizationModule.UserSettingSvc,
	)
	questionSetController := controllers.NewQuestionSetController(questionSetSvc)

	thirdPartyAssessmentStatusRepo := repositories.NewThirdPartyAssessmentStatusRepository(db)
	thirdPartyAssessmentStatusSvc := services.NewThirdPartyAssessmentStatusService(thirdPartyAssessmentStatusRepo)
	thirdPartyAssessmentStatusController := controllers.NewThirdPartyAssessmentStatusController(thirdPartyAssessmentStatusSvc)

	tprmAssessmentWorkflowRepo := repositories.NewThirdPartyAssessmentWorkflowRepository(db)
	tprmAssessmentWorkflowSvc := services.NewThirdPartyAssessmentWorkflowService(tprmAssessmentWorkflowRepo, sharedModule.ActivityLogService, settingsModule.WorkflowSvc, masterModule.WorkflowStatusSvc, organizationModule.UserSvc, masterModule.WorkflowContextSvc)
	tprmAssessmentWorkflowController := controllers.NewThirdPartyAssessmentWorkflowController(tprmAssessmentWorkflowSvc)

	// === Third Party Register ===
	thirdPartyAssessmentRepo := repositories.NewThirdPartyAssessmentRepository(db)
	thirdPartyAssessmentSvc := services.NewThirdPartyAssessmentService(
		thirdPartyAssessmentRepo,
		sharedModule.ActivityLogService,
		masterModule.SeverityLevelSvc,
		organizationModule.UserSvc,
		settingsModule.OrganizationUnitSvc,
		thirdPartyRegisterSvc,
		questionSetSvc,
		thirdPartyAssessmentStatusSvc,
		organizationModule.UserSettingSvc,
		commonModule.NotificationService,
		tprmAssessmentWorkflowSvc,
		tprmComplianceResponseSvc,
	)
	thirdPartyAssessmentController := controllers.NewThirdPartyAssessmentController(thirdPartyAssessmentSvc)

	thirdPartyAssociatedItemRepo := repositories.NewThirdPartyAssociatedItemRepository(db)
	thirdPartyAssociatedItemSvc := services.NewThirdPartyAssociatedItemService(thirdPartyAssociatedItemRepo, thirdPartyRegisterSvc, commonModule.AssociatedItemValidationService)
	thirdPartyAssociatedItemController := controllers.NewThirdPartyAssociatedItemController(thirdPartyAssociatedItemSvc)

	tprmAssessmentReviewRepo := repositories.NewThirdPartyAssessmentReviewRepository(db)
	tprmAssessmentReviewSvc := services.NewThirdPartyAssessmentReviewService(tprmAssessmentReviewRepo, sharedModule.ActivityLogService, thirdPartyAssessmentSvc, tprmAssessmentWorkflowSvc, masterModule.WorkflowStatusSvc, thirdPartyAssessmentStatusSvc, organizationModule.UserSvc, commonModule.NotificationService, taskModule.TaskService, tprmSettingSvc, tprmComplianceResponseSvc, riskModule.RiskLevelSvc, thirdPartyRegisterSvc)
	tprmAssessmentReviewController := controllers.NewThirdPartyAssessmentReviewController(tprmAssessmentReviewSvc)

	tprmAssesmentSectionQuestionCommentRepo := repositories.NewTprmAssessmentSectionQuestionCommentRepository(db)

	tprmAssesmentSectionQuestionCommentSvc := services.NewTprmAssessmentSectionQuestionCommentService(tprmAssesmentSectionQuestionCommentRepo, thirdPartyAssessmentSvc)
	tprmAssesmentSectionQuestionCommentController := controllers.NewTprmAssessmentSectionQuestionCommentController(tprmAssesmentSectionQuestionCommentSvc)

	tprmAssessmentWorkflowSvc.SetThirdPartyAssessmentService(thirdPartyAssessmentSvc)

	tprmDashboardRepo := repositories.NewTPRMDashboardRepository(db)
	tprmDashboardSvc := services.NewTPRMDashboardService(tprmDashboardRepo)
	tprmDashboardController := controllers.NewTPRMDashboardController(tprmDashboardSvc)

	return &TprmModule{
		ThirdPartyTypeController:                      thirdPartyTypeController,
		ThirdPartyRegisterController:                  thirdPartyRegisterController,
		ThirdPartyQuestionCategoryController:          thirdPartyQuestionCategoryController,
		ThirdPartyQuestionTypeController:              thirdPartyQuestionTypeController,
		ThirdPartyQuestionLibraryController:           thirdPartyQuestionLibraryController,
		QuestionSetController:                         questionSetController,
		ThirdPartyAssessmentStatusController:          thirdPartyAssessmentStatusController,
		ThirdPartyAssessmentController:                thirdPartyAssessmentController,
		ThirdPartyAssociatedItemController:            thirdPartyAssociatedItemController,
		TPRMSettingController:                         tprmSettingController,
		ThirdPartyAssessmentWorkflowController:        tprmAssessmentWorkflowController,
		ThirdPartyAssessmentReviewController:          tprmAssessmentReviewController,
		TPRMAssesmentSectionQuestionCommentController: tprmAssesmentSectionQuestionCommentController,
		TPRMComplianeResponseController:               tprmComplianceResponseController,
		TPRMDashboardController:                       tprmDashboardController,

		ThirdPartyTypeSvc:               thirdPartyTypeSvc,
		ThirdPartyRegisterSvc:           thirdPartyRegisterSvc,
		ThirdPartyQuestionCategorySvc:   thirdPartyQuestionCategorySvc,
		ThirdPartyQuestionTypeSvc:       thirdPartyQuestionTypeSvc,
		TPRMSettingSvc:                  tprmSettingSvc,
		ThirdPartyAssessmentWorkflowSvc: tprmAssessmentWorkflowSvc,
		TPRMAssesmentSectionQuestionSvc: tprmAssesmentSectionQuestionCommentSvc,
		TPRMComplianceResponseSvc:       tprmComplianceResponseSvc,
		ThirdPartyAssessmentSvc:         thirdPartyAssessmentSvc,
	}
}
