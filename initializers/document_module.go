package initializers

import (
	"jethur/grc/documents/controllers"
	"jethur/grc/documents/repositories"
	"jethur/grc/documents/services"

	"gorm.io/gorm"
)

type DocumentModule struct {
	DocumentCategoryController   *controllers.DocumentCategoryController
	DocumentTypeController       *controllers.DocumentTypeController
	DocumentTemplateController   *controllers.DocumentTemplateController
	DocumentStatusController     *controllers.DocumentStatusController
	DocumentAccessTypeController *controllers.DocumentAccessTypeController
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

	documentAccessTypeRepo := repositories.NewDocumentAccessTypeRepository(db)
	documentAccessTypeSvc := services.NewDocumentAccessTypeService(documentAccessTypeRepo)
	documentAccessTypeController := controllers.NewDocumentAccessTypeController(documentAccessTypeSvc)

	return &DocumentModule{
		DocumentCategoryController:   documentCategoryController,
		DocumentTypeController:       documentTypeController,
		DocumentTemplateController:   documentTemplateController,
		DocumentStatusController:     documentStatusController,
		DocumentAccessTypeController: documentAccessTypeController,
	}
}
