package initializers

import (
	"jethur/grc/documents/controllers"
	"jethur/grc/documents/repositories"
	"jethur/grc/documents/services"

	"gorm.io/gorm"
)

type DocumentModule struct {
	DocumentCategoryController                    *controllers.DocumentCategoryController
	DocumentTypeController                        *controllers.DocumentTypeController
	DocumentTemplateController                    *controllers.DocumentTemplateController
	DocumentStatusController                      *controllers.DocumentStatusController
	DocumentAccessTypeController                  *controllers.DocumentAccessTypeController
	DocumentWorkflowStatusController              *controllers.DocumentWorkflowStatusController
	DocumentController                            *controllers.DocumentController
	DocumentAcknowledgeUserController             *controllers.DocumentAcknowledgeUserController
	DocumentAuditLogController                    *controllers.DocumentAuditLogController
	DocumentClauseController                      *controllers.DocumentClauseController
	DocumentAssociatedItemController              *controllers.DocumentAssociatedItemController
	DocumentAccessController                      *controllers.DocumentAccessController
	DocumentReviewController                      *controllers.DocumentReviewController
	DocumentWorkflowController                    *controllers.DocumentWorkflowController
	DocumentChangeRequestInitiateTypeController   *controllers.DocumentChangeRequestInitiateTypeController
	DocumentChangeRequestInitiateStatusController *controllers.DocumentChangeRequestInitiateStatusController
	DocumentChangeRequestInitiateController       *controllers.DocumentChangeRequestInitiateController
	DocumentChangeRequestTypeController           *controllers.DocumentChangeRequestTypeController
	DocumentChangeRequestStatusController         *controllers.DocumentChangeRequestStatusController
	DocumentChangeRequestController               *controllers.DocumentChangeRequestController
	DocumentChangeRequestWorkflowStatusController *controllers.DocumentChangeRequestWorkflowStatusController

	DocumentSvc services.DocumentServiceInterface
}

func NewDocumentModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule,
	taskModule *TaskModule,
	commonModule *CommonModule) *DocumentModule {

	documentCategoryRepo := repositories.NewDocumentCategoryRepository(db)
	documentCategorySvc := services.NewDocumentCategoryService(documentCategoryRepo, sharedModule.ActivityLogService)
	documentCategoryController := controllers.NewDocumentCategoryController(documentCategorySvc)

	documentTypeRepo := repositories.NewDocumentTypeRepository(db)
	documentTypeSvc := services.NewDocumentTypeService(documentTypeRepo, sharedModule.ActivityLogService)
	documentTypeController := controllers.NewDocumentTypeController(documentTypeSvc)

	documentTemplateRepo := repositories.NewDocumentTemplateRepository(db)
	documentTemplateSvc := services.NewDocumentTemplateService(documentTemplateRepo, sharedModule.ActivityLogService, documentTypeSvc)
	documentTemplateController := controllers.NewDocumentTemplateController(documentTemplateSvc)

	documentStatusRepo := repositories.NewDocumentStatusRepository(db)
	documentStatusSvc := services.NewDocumentStatusService(documentStatusRepo)
	documentStatusController := controllers.NewDocumentStatusController(documentStatusSvc)

	documentWorkflowStatusRepo := repositories.NewDocumentWorkflowStatusRepository(db)
	documentWorkflowStatusSvc := services.NewDocumentWorkflowStatusService(documentWorkflowStatusRepo)
	documentWorkflowStatusController := controllers.NewDocumentWorkflowStatusController(documentWorkflowStatusSvc)

	documentAccessTypeRepo := repositories.NewDocumentAccessTypeRepository(db)
	documentAccessTypeSvc := services.NewDocumentAccessTypeService(documentAccessTypeRepo)
	documentAccessTypeController := controllers.NewDocumentAccessTypeController(documentAccessTypeSvc)

	documentRepo := repositories.NewDocumentRepository(db)

	documentAuditLogRepo := repositories.NewDocumentAuditLogRepository(db)
	documentAuditLogSvc := services.NewDocumentAuditLogService(documentAuditLogRepo, documentRepo)
	documentAuditLogController := controllers.NewDocumentAuditLogController(documentAuditLogSvc)

	documentWorkflowRepo := repositories.NewDocumentWorkflowRepository(db)
	documentWorkflowSvc := services.NewDocumentWorkflowService(
		documentWorkflowRepo,
		sharedModule.ActivityLogService,
		settingsModule.WorkflowSvc,
		documentWorkflowStatusSvc,
		organizationModule.UserSvc)

	documentSvc := services.NewDocumentService(
		documentRepo,
		documentCategorySvc,
		documentTypeSvc,
		documentTemplateSvc,
		documentAccessTypeSvc,
		documentStatusSvc,
		organizationModule.UserSvc,
		settingsModule.WorkflowSvc,
		documentWorkflowStatusSvc,
		sharedModule.ActivityLogService,
		documentAuditLogSvc,
		organizationModule.UserSettingSvc,
		taskModule.TaskService,
		commonModule.NotificationService,
		settingsModule.OrganizationLevelSvc)

	documentSvc.SetDocumentWorkflowService(documentWorkflowSvc)

	documentAccessRepo := repositories.NewDocumentAccessRepository(db)
	documentAccessSvc := services.NewDocumentAccessService(documentAccessRepo, sharedModule.ActivityLogService, documentSvc)

	documentSvc.SetDocumentAccessService(documentAccessSvc)
	documentWorkflowSvc.SetDocumentAccessService(documentAccessSvc)
	documentWorkflowSvc.SetDocumentService(documentSvc)

	documentController := controllers.NewDocumentController(documentSvc)

	documentAccessController := controllers.NewDocumentAccessController(documentAccessSvc)

	documentWorkflowController := controllers.NewDocumentWorkflowController(documentWorkflowSvc)

	documentReviewRepo := repositories.NewDocumentReviewRepository(db)
	documentReviewSvc := services.NewDocumentReviewService(
		documentReviewRepo,
		sharedModule.ActivityLogService,
		documentSvc,
		documentWorkflowSvc,
		documentWorkflowStatusSvc,
		documentStatusSvc,
		organizationModule.UserSvc,
		commonModule.NotificationService,
		taskModule.TaskService)
	documentReviewController := controllers.NewDocumentReviewController(documentReviewSvc)

	documentAcknowledgeUserRepo := repositories.NewDocumentAcknowledgeUserRepository(db)
	documentAcknowledgeUserSvc := services.NewDocumentAcknowledgeUserService(documentAcknowledgeUserRepo, sharedModule.ActivityLogService, documentSvc, organizationModule.UserSvc, settingsModule.OrganizationLevelSvc)
	documentAcknowledgeUserController := controllers.NewDocumentAcknowledgeUserController(documentAcknowledgeUserSvc)

	documentClauseRepo := repositories.NewDocumentClauseRepository(db)
	documentClauseSvc := services.NewDocumentClauseService(documentClauseRepo, documentSvc, sharedModule.ActivityLogService)
	documentClauseController := controllers.NewDocumentClauseController(documentClauseSvc)

	documentAssociatedItemRepo := repositories.NewDocumentAssociatedItemRepository(db)
	documentAssociatedItemSvc := services.NewDocumentAssociatedItemService(documentAssociatedItemRepo, documentSvc)
	documentAssociatedItemController := controllers.NewDocumentAssociatedItemController(documentAssociatedItemSvc)

	documentChangeRequestInitiateTypeRepo := repositories.NewDocumentChangeRequestInitiateTypeRepository(db)
	documentChangeRequestInitiateTypeSvc := services.NewDocumentChangeRequestInitiateTypeService(documentChangeRequestInitiateTypeRepo)
	documentChangeRequestInitiateTypeController := controllers.NewDocumentChangeRequestInitiateTypeController(documentChangeRequestInitiateTypeSvc)

	documentChangeRequestInitiateStatusRepo := repositories.NewDocumentChangeRequestInitiateStatusRepository(db)
	documentChangeRequestInitiateStatusSvc := services.NewDocumentChangeRequestInitiateStatusService(documentChangeRequestInitiateStatusRepo)
	documentChangeRequestInitiateStatusController := controllers.NewDocumentChangeRequestInitiateStatusController(documentChangeRequestInitiateStatusSvc)

	documentChangeRequestInitiateRepo := repositories.NewDocumentChangeRequestInitiateRepository(db)
	documentChangeRequestInitiateSvc := services.NewDocumentChangeRequestInitiateService(documentChangeRequestInitiateRepo, sharedModule.ActivityLogService, documentChangeRequestInitiateTypeSvc, documentChangeRequestInitiateStatusSvc)
	documentChangeRequestInitiateController := controllers.NewDocumentChangeRequestInitiateController(documentChangeRequestInitiateSvc)

	documentChangeRequestTypeRepo := repositories.NewDocumentChangeRequestTypeRepository(db)
	documentChangeRequestTypeSvc := services.NewDocumentChangeRequestTypeService(documentChangeRequestTypeRepo)
	documentChangeRequestTypeController := controllers.NewDocumentChangeRequestTypeController(documentChangeRequestTypeSvc)

	documentChangeRequestStatusRepo := repositories.NewDocumentChangeRequestStatusRepository(db)
	documentChangeRequestStatusSvc := services.NewDocumentChangeRequestStatusService(documentChangeRequestStatusRepo)
	documentChangeRequestStatusController := controllers.NewDocumentChangeRequestStatusController(documentChangeRequestStatusSvc)

	documentChangeRequestWorkflowStatusRepo := repositories.NewDocumentChangeRequestWorkflowStatusRepository(db)
	documentChangeRequestWorkflowStatusSvc := services.NewDocumentChangeRequestWorkflowStatusService(documentChangeRequestWorkflowStatusRepo)
	documentChangeRequestWorkflowStatusController := controllers.NewDocumentChangeRequestWorkflowStatusController(documentChangeRequestWorkflowStatusSvc)

	documentChangeRequestRepo := repositories.NewDocumentChangeRequestRepository(db)
	documentChangeRequestSvc := services.NewDocumentChangeRequestService(
		documentChangeRequestRepo,
		sharedModule.ActivityLogService,
		documentChangeRequestTypeSvc,
		documentChangeRequestStatusSvc,
		documentChangeRequestWorkflowStatusSvc,
	)
	documentChangeRequestController := controllers.NewDocumentChangeRequestController(documentChangeRequestSvc)

	documentChangeRequestWorkflowRepo := repositories.NewDocumentChangeRequestWorkflowRepository(db)
	documentChangeRequestWorkflowSvc := services.NewDocumentChangeRequestWorkflowService(
		documentChangeRequestWorkflowRepo,
		sharedModule.ActivityLogService,
		documentSvc,
		settingsModule.WorkflowSvc,
		documentChangeRequestWorkflowStatusSvc,
		organizationModule.UserSvc,
	)

	documentChangeRequestSvc.SetDocumentChangeRequestWorkflowService(documentChangeRequestWorkflowSvc)

	return &DocumentModule{
		DocumentCategoryController:                    documentCategoryController,
		DocumentTypeController:                        documentTypeController,
		DocumentTemplateController:                    documentTemplateController,
		DocumentStatusController:                      documentStatusController,
		DocumentAccessTypeController:                  documentAccessTypeController,
		DocumentWorkflowStatusController:              documentWorkflowStatusController,
		DocumentController:                            documentController,
		DocumentAcknowledgeUserController:             documentAcknowledgeUserController,
		DocumentAuditLogController:                    documentAuditLogController,
		DocumentClauseController:                      documentClauseController,
		DocumentAssociatedItemController:              documentAssociatedItemController,
		DocumentAccessController:                      documentAccessController,
		DocumentReviewController:                      documentReviewController,
		DocumentWorkflowController:                    documentWorkflowController,
		DocumentChangeRequestInitiateTypeController:   documentChangeRequestInitiateTypeController,
		DocumentChangeRequestInitiateStatusController: documentChangeRequestInitiateStatusController,
		DocumentChangeRequestInitiateController:       documentChangeRequestInitiateController,
		DocumentChangeRequestTypeController:           documentChangeRequestTypeController,
		DocumentChangeRequestStatusController:         documentChangeRequestStatusController,
		DocumentChangeRequestController:               documentChangeRequestController,
		DocumentChangeRequestWorkflowStatusController: documentChangeRequestWorkflowStatusController,

		DocumentSvc: documentSvc,
	}
}
