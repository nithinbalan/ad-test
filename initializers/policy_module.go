package initializers

import (
	"jethur/grc/policy/controllers"
	"jethur/grc/policy/repositories"
	"jethur/grc/policy/services"

	"gorm.io/gorm"
)

type PolicyModule struct {
	PolicyFolderController        *controllers.PolicyFolderController
	PolicyRegisterController      *controllers.PolicyRegisterController
	PolicyClauseController        *controllers.PolicyClauseController
	PolicyReviewController        *controllers.PolicyReviewController
	PolicyWorkflowController      *controllers.PolicyWorkflowController
	PolicyClauseCommentController *controllers.PolicyClauseCommentController
	PolicyReportController        *controllers.PolicyReportController

	PolicyFolderSvc   services.PolicyFolderServiceInterface
	PolicyRegisterSvc services.PolicyRegisterServiceInterface
	PolicyClauseSvc   services.PolicyClauseServiceInterface
}

func NewPolicyModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule,
	documentModule *DocumentModule) *PolicyModule {

	policyWorkflowController := controllers.NewPolicyWorkflowController(documentModule.DocumentWorkflowSvc)

	policyReviewController := controllers.NewPolicyReviewController(documentModule.DocumentReviewSvc)

	policyFolderRepo := repositories.NewPolicyFolderRepository(db)
	policyFolderSvc := services.NewPolicyFolderService(policyFolderRepo, sharedModule.ActivityLogService, documentModule.DocumentAccessTypeSvc, settingsModule.OrganizationLevelSvc, organizationModule.UserSvc)
	policyFolderController := controllers.NewPolicyFolderController(policyFolderSvc)

	policyClauseRepo := repositories.NewPolicyClauseRepository(db)
	policyClauseSvc := services.NewPolicyClauseService(policyClauseRepo, documentModule.DocumentSvc, documentModule.DocumentStatusSvc, sharedModule.ActivityLogService)
	policyClauseController := controllers.NewPolicyClauseController(policyClauseSvc)

	policyRegisterRepo := repositories.NewPolicyRegisterRepository(db)
	policyRegisterSvc := services.NewPolicyRegisterService(
		policyRegisterRepo,
		documentModule.DocumentCategorySvc,
		documentModule.DocumentTypeSvc,
		documentModule.DocumentAccessTypeSvc,
		documentModule.DocumentStatusSvc,
		organizationModule.UserSvc,
		sharedModule.ActivityLogService,
		organizationModule.UserSettingSvc,
		documentModule.DocumentAccessSvc,
		taskModule.TaskService,
		commonModule.NotificationService,
		settingsModule.OrganizationLevelSvc,
		masterModule.ReviewFrequencySvc,
		settingsModule.SystemPreferenceSvc,
		documentModule.DocumentWorkflowSvc,
		documentModule.DocumentWorkflowStatusSvc,
		masterModule.WorkflowContextSvc,
		settingsModule.WorkflowSvc,
		documentModule.DocumentAuditLogSvc,
		policyClauseSvc,
		commonModule.FileService,
		settingsModule.ReportTemplateSettingSvc,
	)
	policyRegisterController := controllers.NewPolicyRegisterController(policyRegisterSvc)

	PolicyClauseCommentRepo := repositories.NewPolicyClauseCommentRepository(db)
	PolicyClauseCommentSvc := services.NewPolicyClauseCommentService(PolicyClauseCommentRepo, organizationModule.UserSvc, sharedModule.ActivityLogService, policyRegisterSvc, masterModule.CommentStatusSvc)
	PolicyClauseCommentController := controllers.NewPolicyClauseCommentController(PolicyClauseCommentSvc)

	policyReportRepo := repositories.NewPolicyReportRepository(db)
	policyReportSvc := services.NewPolicyReportService(policyReportRepo, sharedModule.ActivityLogService, documentModule.DocumentSvc, documentModule.DocumentClauseSvc, settingsModule.OrganizationLevelSvc, settingsModule.ReportTemplateSettingSvc)
	policyReportController := controllers.NewPolicyReportController(policyReportSvc)

	return &PolicyModule{
		PolicyFolderController:        policyFolderController,
		PolicyRegisterController:      policyRegisterController,
		PolicyClauseController:        policyClauseController,
		PolicyWorkflowController:      policyWorkflowController,
		PolicyReviewController:        policyReviewController,
		PolicyClauseCommentController: PolicyClauseCommentController,
		PolicyReportController:        policyReportController,
	}
}
