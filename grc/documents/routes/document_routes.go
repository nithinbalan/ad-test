package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	documentRoutes := router.Group("/")

	documentRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	documentCategoryController := app.DocumentModule.DocumentCategoryController
	documentTypeController := app.DocumentModule.DocumentTypeController
	documentTemplateController := app.DocumentModule.DocumentTemplateController
	documentStatusController := app.DocumentModule.DocumentStatusController
	documentAccessTypeController := app.DocumentModule.DocumentAccessTypeController
	documentWorkflowStatusController := app.DocumentModule.DocumentWorkflowStatusController
	documentController := app.DocumentModule.DocumentController
	documentAcknowledgeUserController := app.DocumentModule.DocumentAcknowledgeUserController
	documentAuditLogController := app.DocumentModule.DocumentAuditLogController

	documentClauseController := app.DocumentModule.DocumentClauseController

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

	documentRoutes.GET("/document-workflow-statuses", documentWorkflowStatusController.GetAllDocumentWorkflowStatuses)
	documentRoutes.GET("/document-workflow-statuses/:id", documentWorkflowStatusController.GetDocumentWorkflowStatusDetails)
	//document
	documentRoutes.GET("/documents", documentController.GetAllDocuments)
	documentRoutes.GET("/documents/:id", documentController.GetDocumentDetails)
	documentRoutes.POST("/documents", documentController.CreateDocument)
	documentRoutes.POST("/documents/delete", documentController.DeleteDocuments)
	documentRoutes.GET("/documents/:id/check-out", documentController.CheckoutDocument)
	documentRoutes.POST("/documents/:id/check-in", documentController.CheckinDocument)
	documentRoutes.POST("/documents/:id/check-out-revert", documentController.RevertDocumentCheckout)
	documentRoutes.POST("/documents/:id/submit", documentController.SubmitDocument)
	documentRoutes.POST("/documents/:id/publish", documentController.PublishDocument)
	documentRoutes.GET("/documents/:id/shared-users", documentController.GetSharedUsers)
	documentRoutes.GET("/documents/:id/review-users", documentController.GetReviewUsers)
	documentRoutes.GET("/documents/:id/shared-units", documentController.GetSharedUnits)

	unauthenticatedRoutes.GET("/documents/:id/document-preview", documentController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/documents/:id/document-download", documentController.DownloadDocument)

	documentRoutes.GET("/documents/:id/comments", documentController.GetComments)
	documentRoutes.POST("/documents/:id/comments", documentController.AddComment)
	documentRoutes.GET("/documents/:id/workflow-list", documentController.GetWorkflowList)

	documentRoutes.POST("/documents/:id/acknowledge", documentAcknowledgeUserController.CreateDocumentAcknowledgeUser)
	documentRoutes.DELETE("/documents/:id/acknowledge", documentAcknowledgeUserController.DeleteDocumentAcknowledgeUser)

	documentRoutes.GET("/documents/:id/audit-logs", documentAuditLogController.GetAllDocumentAuditLogs)

	documentRoutes.GET("/documents/:id/clauses", documentClauseController.GetAllDocumentClauses)
	documentRoutes.GET("/documents/:id/clauses/:clauseId", documentClauseController.GetDocumentClauseDetails)
	documentRoutes.POST("/documents/:id/clauses", documentClauseController.CreateDocumentClause)
	documentRoutes.PUT("/documents/:id/clauses/:clauseId", documentClauseController.UpdateDocumentClause)
	documentRoutes.POST("/documents/:id/clauses/delete", documentClauseController.DeleteDocumentClauses)

}
