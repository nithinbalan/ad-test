package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	documentRoutes := router.Group("/")

	documentRoutes.Use(middlewares.AuthMiddleware())

	documentCategoryController := app.DocumentModule.DocumentCategoryController
	documentTypeController := app.DocumentModule.DocumentTypeController
	documentTemplateController := app.DocumentModule.DocumentTemplateController
	documentStatusController := app.DocumentModule.DocumentStatusController
	documentAccessTypeController := app.DocumentModule.DocumentAccessTypeController
	documentController := app.DocumentModule.DocumentController

	documentRoutes.GET("/document-categories", documentCategoryController.GetAllDocumentCategories)
	documentRoutes.GET("/document-categories/:id", documentCategoryController.GetDocumentCategoryDetails)
	documentRoutes.POST("/document-categories", documentCategoryController.CreateDocumentCategory)
	documentRoutes.PUT("/document-categories/:id", documentCategoryController.UpdateDocumentCategory)
	documentRoutes.POST("/document-categories/delete", documentCategoryController.DeleteDocumentCategories)
	documentRoutes.PUT("/document-categories/activate", documentCategoryController.ActivateDocumentCategories)
	documentRoutes.PUT("/document-categories/deactivate", documentCategoryController.DeactivateDocumentCategories)

	documentRoutes.GET("/document-types", documentTypeController.GetAllDocumentTypes)
	documentRoutes.GET("/document-types/:id", documentTypeController.GetDocumentTypeDetails)
	documentRoutes.POST("/document-types", documentTypeController.CreateDocumentType)
	documentRoutes.PUT("/document-types/:id", documentTypeController.UpdateDocumentType)
	documentRoutes.POST("/document-types/delete", documentTypeController.DeleteDocumentTypes)
	documentRoutes.PUT("/document-types/activate", documentTypeController.ActivateDocumentTypes)
	documentRoutes.PUT("/document-types/deactivate", documentTypeController.DeactivateDocumentTypes)

	documentRoutes.GET("/document-templates", documentTemplateController.GetAllDocumentTemplates)
	documentRoutes.GET("/document-templates/:id", documentTemplateController.GetDocumentTemplateDetails)
	documentRoutes.POST("/document-templates", documentTemplateController.CreateDocumentTemplate)
	documentRoutes.PUT("/document-templates/:id", documentTemplateController.UpdateDocumentTemplate)
	documentRoutes.POST("/document-templates/delete", documentTemplateController.DeleteDocumentTemplates)
	documentRoutes.PUT("/document-templates/activate", documentTemplateController.ActivateDocumentTemplates)
	documentRoutes.PUT("/document-templates/deactivate", documentTemplateController.DeactivateDocumentTemplates)

	documentRoutes.GET("/document-statuses", documentStatusController.GetAllDocumentStatuses)
	documentRoutes.GET("/document-statuses/:id", documentStatusController.GetDocumentStatusDetails)

	documentRoutes.GET("/document-access-types", documentAccessTypeController.GetAllDocumentAccessTypes)
	documentRoutes.GET("/document-access-types/:id", documentAccessTypeController.GetDocumentAccessTypeDetails)

	//document
	documentRoutes.GET("/documents", documentController.GetAllDocuments)
	documentRoutes.GET("/documents/:id", documentController.GetDocumentDetails)
	documentRoutes.POST("/documents", documentController.CreateDocument)
	documentRoutes.POST("/documents/delete", documentController.DeleteDocuments)
	documentRoutes.GET("/documents/:id/check-out", documentController.CheckoutDocument)
	documentRoutes.POST("/documents/:id/check-in", documentController.CheckinDocument)
	documentRoutes.POST("/documents/:id/check-out-revert", documentController.RevertDocumentCheckout)

}
