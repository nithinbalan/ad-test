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

	db := app.DB

	documentCategoryController := app.DocumentModule.DocumentCategoryController
	documentTypeController := app.DocumentModule.DocumentTypeController
	documentTemplateController := app.DocumentModule.DocumentTemplateController
	documentStatusController := app.DocumentModule.DocumentStatusController
	documentAccessTypeController := app.DocumentModule.DocumentAccessTypeController
	documentWorkflowStatusController := app.DocumentModule.DocumentWorkflowStatusController
	documentController := app.DocumentModule.DocumentController
	documentReviewController := app.DocumentModule.DocumentReviewController
	documentWorkflowController := app.DocumentModule.DocumentWorkflowController
	documentAccessController := app.DocumentModule.DocumentAccessController
	documentAcknowledgeUserController := app.DocumentModule.DocumentAcknowledgeUserController
	documentAuditLogController := app.DocumentModule.DocumentAuditLogController

	documentClauseController := app.DocumentModule.DocumentClauseController
	documentAssociatedItemController := app.DocumentModule.DocumentAssociatedItemController

	documentRoutes.GET("/document-categories", documentCategoryController.GetAllDocumentCategories)
	documentRoutes.GET("/document-categories/:id", documentCategoryController.GetDocumentCategoryDetails)
	documentRoutes.POST("/document-categories", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentCategoryController.CreateDocumentCategory)
	documentRoutes.PUT("/document-categories/:id", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentCategoryController.UpdateDocumentCategory)
	documentRoutes.POST("/document-categories/delete", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentCategoryController.DeleteDocumentCategories)
	documentRoutes.PUT("/document-categories/activate", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentCategoryController.ActivateDocumentCategories)
	documentRoutes.PUT("/document-categories/deactivate", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentCategoryController.DeactivateDocumentCategories)

	documentRoutes.GET("/document-types", documentTypeController.GetAllDocumentTypes)
	documentRoutes.GET("/document-types/:id", documentTypeController.GetDocumentTypeDetails)
	documentRoutes.POST("/document-types", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTypeController.CreateDocumentType)
	documentRoutes.PUT("/document-types/:id", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTypeController.UpdateDocumentType)
	documentRoutes.POST("/document-types/delete", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTypeController.DeleteDocumentTypes)
	documentRoutes.PUT("/document-types/activate", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTypeController.ActivateDocumentTypes)
	documentRoutes.PUT("/document-types/deactivate", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTypeController.DeactivateDocumentTypes)

	documentRoutes.GET("/document-templates", documentTemplateController.GetAllDocumentTemplates)
	documentRoutes.GET("/document-templates/:id", documentTemplateController.GetDocumentTemplateDetails)
	documentRoutes.POST("/document-templates", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTemplateController.CreateDocumentTemplate)
	documentRoutes.PUT("/document-templates/:id", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTemplateController.UpdateDocumentTemplate)
	documentRoutes.POST("/document-templates/delete", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTemplateController.DeleteDocumentTemplates)
	documentRoutes.PUT("/document-templates/activate", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTemplateController.ActivateDocumentTemplates)
	documentRoutes.PUT("/document-templates/deactivate", middlewares.RolePermissionMiddleware(db, "document_setting", "modify"), documentTemplateController.DeactivateDocumentTemplates)

	documentRoutes.GET("/document-statuses", documentStatusController.GetAllDocumentStatuses)
	documentRoutes.GET("/document-statuses/:id", documentStatusController.GetDocumentStatusDetails)

	documentRoutes.GET("/document-access-types", documentAccessTypeController.GetAllDocumentAccessTypes)
	documentRoutes.GET("/document-access-types/:id", documentAccessTypeController.GetDocumentAccessTypeDetails)

	documentRoutes.GET("/document-workflow-statuses", documentWorkflowStatusController.GetAllDocumentWorkflowStatuses)
	documentRoutes.GET("/document-workflow-statuses/:id", documentWorkflowStatusController.GetDocumentWorkflowStatusDetails)
	//document
	documentRoutes.GET("/documents/export", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentController.GetDocumentExport)
	documentRoutes.GET("/documents", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentController.GetAllDocuments)
	documentRoutes.GET("/documents/:id", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentController.GetDocumentDetails)
	documentRoutes.POST("/documents", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentController.CreateDocument)
	documentRoutes.PUT("/documents/:id", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentController.UpdateDocument)
	documentRoutes.POST("/documents/delete", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentController.DeleteDocuments)
	documentRoutes.GET("/documents/:id/check-out", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentController.CheckoutDocument)
	documentRoutes.POST("/documents/:id/check-in", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentController.CheckinDocument)
	documentRoutes.POST("/documents/:id/check-out-revert", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentController.RevertDocumentCheckout)
	documentRoutes.POST("/documents/:id/submit", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentController.SubmitDocument)
	documentRoutes.POST("/documents/:id/publish", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentController.PublishDocument)

	documentRoutes.GET("/documents/:id/shared-users", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAccessController.GetSharedUsers)
	documentRoutes.POST("/documents/:id/shared-users", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentAccessController.CreateSharedUsers)
	documentRoutes.POST("/documents/:id/shared-users/delete", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentAccessController.DeleteSharedUsers)
	documentRoutes.GET("/documents/:id/shared-units", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAccessController.GetSharedUnits)
	documentRoutes.POST("/documents/:id/shared-units", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentAccessController.CreateSharedUnits)
	documentRoutes.POST("/documents/:id/shared-units/delete", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentAccessController.DeleteSharedUnits)
	documentRoutes.GET("/documents/:id/review-users", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAccessController.GetReviewUsers)

	unauthenticatedRoutes.GET("/documents/:id/document-preview", documentController.GetDocumentPreview)
	unauthenticatedRoutes.GET("/documents/:id/document-download", documentController.DownloadDocument)

	documentRoutes.GET("/documents/:id/comments", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentReviewController.GetComments)
	documentRoutes.POST("/documents/:id/comments", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentReviewController.AddComment)
	documentRoutes.GET("/documents/:id/workflow-list", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentWorkflowController.GetWorkflowList)

	documentRoutes.GET("/documents/:id/last-edit-summary", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentController.GetLastEditSummary)

	documentRoutes.GET("/documents/:id/acknowledge-users", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAcknowledgeUserController.GetAllDocumentAcknowledgeUsers)
	documentRoutes.POST("/documents/:id/acknowledge", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentAcknowledgeUserController.CreateDocumentAcknowledgeUser)
	documentRoutes.DELETE("/documents/:id/acknowledge", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentAcknowledgeUserController.DeleteDocumentAcknowledgeUser)
	documentRoutes.GET("/documents/:id/acknowledge-user-summary", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAcknowledgeUserController.GetDocumentAcknowledgmentSummary)

	documentRoutes.GET("/documents/:id/audit-logs", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAuditLogController.GetAllDocumentAuditLogs)

	documentRoutes.GET("/documents/:id/clauses", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentClauseController.GetAllDocumentClauses)
	documentRoutes.GET("/documents/:id/clauses/:clauseId", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentClauseController.GetDocumentClauseDetails)
	documentRoutes.POST("/documents/:id/clauses", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentClauseController.CreateDocumentClause)
	documentRoutes.PUT("/documents/:id/clauses/:clauseId", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentClauseController.UpdateDocumentClause)
	documentRoutes.POST("/documents/:id/clauses/delete", middlewares.RolePermissionMiddleware(db, "document_register", "modify"), documentClauseController.DeleteDocumentClauses)

	//document associated item list
	documentRoutes.GET("documents/:id/controls", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAssociatedItemController.GetAllDocumentControls)
	documentRoutes.GET("documents/:id/compliance-requirements", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAssociatedItemController.GetAllDocumentComplianceRequirements)
	documentRoutes.GET("documents/:id/processes", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAssociatedItemController.GetAllDocumentProcesses)
	documentRoutes.GET("documents/:id/slas", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAssociatedItemController.GetAllDocumentSLAs)
	documentRoutes.GET("documents/:id/incidents", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAssociatedItemController.GetAllDocumentIncidents)
	documentRoutes.GET("documents/:id/risks", middlewares.RolePermissionMiddleware(db, "document_register", "view"), documentAssociatedItemController.GetAllDocumentRisks)

}
