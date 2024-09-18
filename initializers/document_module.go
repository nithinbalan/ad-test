package initializers

import (
	"jethur/grc/documents/controllers"
	"jethur/grc/documents/repositories"
	"jethur/grc/documents/services"

	"gorm.io/gorm"
)

type DocumentModule struct {
	DocumentCategoryController        *controllers.DocumentCategoryController
	DocumentTypeController            *controllers.DocumentTypeController
	DocumentTemplateController        *controllers.DocumentTemplateController
	DocumentStatusController          *controllers.DocumentStatusController
	DocumentAccessTypeController      *controllers.DocumentAccessTypeController
	DocumentWorkflowStatusController  *controllers.DocumentWorkflowStatusController
	DocumentController                *controllers.DocumentController
	DocumentAcknowledgeUserController *controllers.DocumentAcknowledgeUserController
	DocumentAuditLogController        *controllers.DocumentAuditLogController
	DocumentClauseController          *controllers.DocumentClauseController
	DocumentAssociatedItemController  *controllers.DocumentAssociatedItemController

	DocumentSvc services.DocumentServiceInterface
}

func NewDocumentModule(db *gorm.DB,
	sharedModule *SharedModule,
	settingsModule *SettingsModule,
	masterModule *MasterModule,
	organizationModule *OrganizationModule) *DocumentModule {

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
		documentWorkflowStatusSvc,
		sharedModule.ActivityLogService,
		documentAuditLogSvc,
		organizationModule.UserSettingSvc)
	documentController := controllers.NewDocumentController(documentSvc)

	documentAcknowledgeUserRepo := repositories.NewDocumentAcknowledgeUserRepository(db)
	documentAcknowledgeUserSvc := services.NewDocumentAcknowledgeUserService(documentAcknowledgeUserRepo, sharedModule.ActivityLogService, documentSvc)
	documentAcknowledgeUserController := controllers.NewDocumentAcknowledgeUserController(documentAcknowledgeUserSvc)

	documentClauseRepo := repositories.NewDocumentClauseRepository(db)
	documentClauseSvc := services.NewDocumentClauseService(documentClauseRepo, documentSvc, sharedModule.ActivityLogService)
	documentClauseController := controllers.NewDocumentClauseController(documentClauseSvc)

	documentAssociatedItemRepo := repositories.NewDocumentAssociatedItemRepository(db)
	documentAssociatedItemSvc := services.NewDocumentAssociatedItemService(documentAssociatedItemRepo, documentSvc)
	documentAssociatedItemController := controllers.NewDocumentAssociatedItemController(documentAssociatedItemSvc)

	return &DocumentModule{
		DocumentCategoryController:        documentCategoryController,
		DocumentTypeController:            documentTypeController,
		DocumentTemplateController:        documentTemplateController,
		DocumentStatusController:          documentStatusController,
		DocumentAccessTypeController:      documentAccessTypeController,
		DocumentWorkflowStatusController:  documentWorkflowStatusController,
		DocumentController:                documentController,
		DocumentAcknowledgeUserController: documentAcknowledgeUserController,
		DocumentAuditLogController:        documentAuditLogController,
		DocumentClauseController:          documentClauseController,
		DocumentAssociatedItemController:  documentAssociatedItemController,

		DocumentSvc: documentSvc,
	}
}
