package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	violationReportRoutes := router.Group("/")

	violationReportRoutes.Use(middlewares.AuthMiddleware(app.DB))

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	violationReportChannelController := app.ViolationReportModule.ViolationReportChannelController
	violationReportReasonController := app.ViolationReportModule.ViolationReportReasonController
	violationReportCategoryController := app.ViolationReportModule.ViolationReportCategoryController
	violationReportSubCategoryController := app.ViolationReportModule.ViolationReportSubCategoryController
	violationReportStatusController := app.ViolationReportModule.ViolationReportStatusController
	violationReportController := app.ViolationReportModule.ViolationReportController
	violationReportAssociatedItemController := app.ViolationReportModule.ViolationReportAssociatedItemController
	violationReportInvestigationController := app.ViolationReportModule.ViolationReportInvestigationController
	violationReportRootCauseAnalysisController := app.ViolationReportModule.ViolationReportRootCauseAnalysisController
	violationReportActionPlanController := app.ViolationReportModule.ViolationReportActionPlanController
	violationReportActionPlanStatusUpdateController := app.ViolationReportModule.ViolationReportActionPlanStatusUpdateController
	violationReportActionPlanCommentController := app.ViolationReportModule.ViolationReportActionPlanCommentController

	// Violation report channels
	violationReportRoutes.GET("/violation-report-channels", violationReportChannelController.GetAllViolationReportChannels)
	violationReportRoutes.GET("/violation-report-channels/:id", violationReportChannelController.GetViolationReportChannelDetails)
	violationReportRoutes.POST("/violation-report-channels", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportChannelController.CreateViolationReportChannel)
	violationReportRoutes.PUT("/violation-report-channels/:id", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportChannelController.UpdateViolationReportChannel)
	violationReportRoutes.POST("/violation-report-channels/delete", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportChannelController.DeleteViolationReportChannels)

	// Violation report reasons
	violationReportRoutes.GET("/violation-report-reasons", violationReportReasonController.GetAllViolationReportReasons)
	violationReportRoutes.GET("/violation-report-reasons/:id", violationReportReasonController.GetViolationReportReasonDetails)
	violationReportRoutes.POST("/violation-report-reasons", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportReasonController.CreateViolationReportReason)
	violationReportRoutes.PUT("/violation-report-reasons/:id", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportReasonController.UpdateViolationReportReason)
	violationReportRoutes.POST("/violation-report-reasons/delete", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportReasonController.DeleteViolationReportReasons)

	// Violation report categories
	violationReportRoutes.GET("/violation-report-categories", violationReportCategoryController.GetAllViolationReportCategories)
	violationReportRoutes.GET("/violation-report-categories/:id", violationReportCategoryController.GetViolationReportCategoryDetails)
	violationReportRoutes.POST("/violation-report-categories", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportCategoryController.CreateViolationReportCategory)
	violationReportRoutes.PUT("/violation-report-categories/:id", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportCategoryController.UpdateViolationReportCategory)
	violationReportRoutes.POST("/violation-report-categories/delete", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportCategoryController.DeleteViolationReportCategories)

	// Violation report sub categories
	violationReportRoutes.GET("/violation-report-sub-categories", violationReportSubCategoryController.GetAllViolationReportSubCategories)
	violationReportRoutes.GET("/violation-report-sub-categories/:id", violationReportSubCategoryController.GetViolationReportSubCategoryDetails)
	violationReportRoutes.POST("/violation-report-sub-categories", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportSubCategoryController.CreateViolationReportSubCategory)
	violationReportRoutes.PUT("/violation-report-sub-categories/:id", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportSubCategoryController.UpdateViolationReportSubCategory)
	violationReportRoutes.POST("/violation-report-sub-categories/delete", middlewares.RolePermissionMiddleware(db, "violation_report_setting", "modify"), violationReportSubCategoryController.DeleteViolationReportSubCategories)

	//Violation report statuses
	violationReportRoutes.GET("/violation-report-statuses", violationReportStatusController.GetAllViolationReportStatuses)
	violationReportRoutes.GET("/violation-report-statuses/:id", violationReportStatusController.GetViolationReportStatusDetails)

	//Violation Report
	unauthenticatedRoutes.GET("/violation-reports/:id/document-preview", violationReportController.GetViolationReportDocumentPreview)
	unauthenticatedRoutes.GET("/violation-reports/:id/document-download", violationReportController.DownloadViolationReportDocument)
	violationReportRoutes.GET("/violation-reports/export", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportController.GetViolationReportExport)
	violationReportRoutes.GET("/violation-reports/:id/export-pdf", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportController.ExportViolationReportToPDF)
	violationReportRoutes.GET("/violation-reports", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportController.GetAllViolationReports)
	violationReportRoutes.GET("/violation-reports/:id", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportController.GetViolationReportDetails)
	violationReportRoutes.POST("/violation-reports", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportController.CreateViolationReport)
	violationReportRoutes.PUT("/violation-reports/:id", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportController.UpdateViolationReport)
	violationReportRoutes.POST("/violation-reports/delete", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportController.DeleteViolationReports)

	//Violation Report Associated Items
	violationReportRoutes.GET("/violation-reports/:id/controls", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportAssociatedItemController.GetAllViolationReportControls)
	violationReportRoutes.GET("/violation-reports/:id/compliance-registers", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportAssociatedItemController.GetAllViolationReportComplianceRegisters)
	violationReportRoutes.GET("/violation-reports/:id/risks", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportAssociatedItemController.GetAllViolationReportRisks)
	violationReportRoutes.GET("/violation-reports/:id/documents", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportAssociatedItemController.GetAllViolationReportDocuments)
	violationReportRoutes.GET("/violation-reports/:id/assets", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportAssociatedItemController.GetAllViolationReportAssets)

	violationReportRoutes.POST("/violation-reports/:id/assets", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.CreateViolationReportAsset)
	violationReportRoutes.DELETE("/violation-reports/:id/assets/:assetId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.DeleteViolationReportAsset)

	violationReportRoutes.POST("/violation-reports/:id/risks", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.CreateViolationReportRisk)
	violationReportRoutes.DELETE("/violation-reports/:id/risks/:riskId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.DeleteViolationReportRisk)

	violationReportRoutes.POST("/violation-reports/:id/documents", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.CreateViolationReportDocument)
	violationReportRoutes.DELETE("/violation-reports/:id/documents/:documentId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.DeleteViolationReportDocument)

	violationReportRoutes.POST("/violation-reports/:id/controls", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.CreateViolationReportControl)
	violationReportRoutes.DELETE("/violation-reports/:id/controls/:controlId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.DeleteViolationReportControl)

	violationReportRoutes.POST("/violation-reports/:id/compliance-registers", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportAssociatedItemController.CreateViolationReportComplianceRegister)
	violationReportRoutes.DELETE("/violation-reports/:id/compliance-registers/:complianceRegisterId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportAssociatedItemController.DeleteViolationReportComplianceRegister)

	// Violation Report Investigation
	unauthenticatedRoutes.GET("/violation-reports/:id/investigation/document-preview", violationReportInvestigationController.GetViolationReportInvestigationDocumentPreview)
	unauthenticatedRoutes.GET("/violation-reports/:id/investigation/document-download", violationReportInvestigationController.DownloadViolationReportInvestigationDocument)
	violationReportRoutes.GET("/violation-reports/:id/investigation", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportInvestigationController.GetViolationReportInvestigationDetails)
	violationReportRoutes.POST("/violation-reports/:id/investigation", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportInvestigationController.CreateViolationReportInvestigation)
	violationReportRoutes.PUT("/violation-reports/:id/investigation", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportInvestigationController.UpdateViolationReportInvestigation)
	violationReportRoutes.DELETE("/violation-reports/:id/investigation", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportInvestigationController.DeleteViolationReportInvestigations)
	violationReportRoutes.PUT("/violation-report-investigation/:investigationId/notes", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportInvestigationController.UpdateInvestigationNoteAndObservations)

	// Violation Report Root Cause Analyses
	violationReportRoutes.GET("/violation-reports/:id/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportRootCauseAnalysisController.GetAllViolationReportRootCauseAnalyses)
	violationReportRoutes.POST("/violation-reports/:id/root-cause-analysis", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportRootCauseAnalysisController.CreateViolationReportRootCauseAnalysis)
	violationReportRoutes.PUT("/violation-reports/:id/root-cause-analysis/:rootCauseAnalysisID", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportRootCauseAnalysisController.UpdateViolationReportRootCauseAnalysis)
	violationReportRoutes.DELETE("/violation-reports/:id/root-cause-analysis/:rootCauseAnalysisID", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportRootCauseAnalysisController.DeleteViolationReportRootCauseAnalyses)

	// Violation Report Action Plan
	violationReportRoutes.POST("/violation-report-action-plans/force-delete", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanController.DeleteForceViolationReportActionPlans)
	violationReportRoutes.GET("/violation-report-action-plans/export", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanController.GetViolationReportActionPlanExport)
	violationReportRoutes.GET("/violation-report-action-plans", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanController.GetAllViolationReportActionPlans)
	violationReportRoutes.GET("/violation-reports/:id/action-plans", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanController.GetAllViolationReportActionPlans)
	violationReportRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanController.GetViolationReportActionPlanDetails)
	violationReportRoutes.POST("/violation-reports/:id/action-plans", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportActionPlanController.CreateViolationReportActionPlan)
	violationReportRoutes.PUT("/violation-reports/:id/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportActionPlanController.UpdateViolationReportActionPlan)
	violationReportRoutes.POST("/violation-reports/:id/action-plans/delete", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportActionPlanController.DeleteViolationReportActionPlans)

	// ViolationReport Action Plan Status Updates
	unauthenticatedRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-preview", violationReportActionPlanStatusUpdateController.GetViolationReportActionPlanStatusUpdateDocumentPreview)
	unauthenticatedRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId/document-download", violationReportActionPlanStatusUpdateController.DownloadViolationReportActionPlanStatusUpdateDocument)
	violationReportRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanStatusUpdateController.GetAllViolationReportActionPlanStatusUpdates)
	violationReportRoutes.POST("/violation-reports/:id/action-plans/:actionPlanId/status-updates", violationReportActionPlanStatusUpdateController.CreateViolationReportActionPlanStatusUpdate)
	violationReportRoutes.PUT("/violation-reports/:id/action-plans/:actionPlanId/status-updates/:statusUpdateId", violationReportActionPlanStatusUpdateController.UpdateViolationReportActionPlanStatusUpdate)
	violationReportRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanStatusUpdateController.GetLatestViolationReportActionPlanStatusUpdate)
	violationReportRoutes.PUT("/violation-reports/:id/action-plans/:actionPlanId/approve", violationReportActionPlanStatusUpdateController.ApproveViolationReportActionPlan)
	violationReportRoutes.PUT("/violation-reports/:id/action-plans/:actionPlanId/revert", violationReportActionPlanStatusUpdateController.RevertViolationReportActionPlan)

	// Violation Report Action Plan Comments
	unauthenticatedRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId/comments/:commentId/document-download", violationReportActionPlanCommentController.DownloadViolationReportActionPlanCommentDocument)
	unauthenticatedRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId/comments/:commentId/document-preview", violationReportActionPlanCommentController.GetViolationReportActionPlanCommentDocumentPreview)
	violationReportRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanCommentController.GetAllComments)
	violationReportRoutes.GET("/violation-reports/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "view"), violationReportActionPlanCommentController.GetCommentDetails)
	violationReportRoutes.POST("/violation-reports/:id/action-plans/:actionPlanId/comments", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportActionPlanCommentController.CreateComment)
	violationReportRoutes.PUT("/violation-reports/:id/action-plans/:actionPlanId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportActionPlanCommentController.UpdateComment)
	violationReportRoutes.POST("/violation-reports/:id/action-plans/:actionPlanId/comments/delete", middlewares.RolePermissionMiddleware(db, "violation_report_register", "modify"), violationReportActionPlanCommentController.DeleteComment)

}
