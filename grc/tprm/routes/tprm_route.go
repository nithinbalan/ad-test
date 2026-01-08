package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	tprmRoutes := router.Group("/")

	tprmRoutes.Use(middlewares.AuthMiddleware(app.DB))

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	thirdPartyTypeController := app.TprmModule.ThirdPartyTypeController
	thirdPartyRegisterController := app.TprmModule.ThirdPartyRegisterController
	thirdPartyQuestionCategoryController := app.TprmModule.ThirdPartyQuestionCategoryController
	thirdPartyQuestionTypeController := app.TprmModule.ThirdPartyQuestionTypeController
	thirdPartyQuestionLibraryController := app.TprmModule.ThirdPartyQuestionLibraryController
	questionSetController := app.TprmModule.QuestionSetController
	thirdPartyAssessmentStatusController := app.TprmModule.ThirdPartyAssessmentStatusController
	thirdPartyAssessmentController := app.TprmModule.ThirdPartyAssessmentController
	thirdPartyAssociatedItemController := app.TprmModule.ThirdPartyAssociatedItemController
	tprmSettingController := app.TprmModule.TPRMSettingController
	thirdPartyAssessmentWorkflowController := app.TprmModule.ThirdPartyAssessmentWorkflowController
	thirdPartyAssessmentReviewController := app.TprmModule.ThirdPartyAssessmentReviewController
	tprmAssesmentSectionQuestionCommentController := app.TprmModule.TPRMAssesmentSectionQuestionCommentController
	tprmComplianceResponseController := app.TprmModule.TPRMComplianeResponseController
	tprmDashboardController := app.TprmModule.TPRMDashboardController

	tprmRoutes.GET("/third-party-types", thirdPartyTypeController.GetAllThirdPartyTypes)
	tprmRoutes.GET("/third-party-types/:id", thirdPartyTypeController.GetThirdPartyTypeDetails)
	tprmRoutes.POST("/third-party-types", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.CreateThirdPartyType)
	tprmRoutes.PUT("/third-party-types/:id", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.UpdateThirdPartyType)
	tprmRoutes.POST("/third-party-types/delete", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.DeleteThirdPartyTypes)
	tprmRoutes.PUT("/third-party-types/activate", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.ActivateThirdPartyTypes)
	tprmRoutes.PUT("/third-party-types/deactivate", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyTypeController.DeactivateThirdPartyTypes)

	// === Third Party Register ===
	tprmRoutes.GET("/third-party-register", thirdPartyRegisterController.GetAllThirdPartyRegister)
	tprmRoutes.GET("/third-party-register/export", thirdPartyRegisterController.GetThirdPartyRegisterExport)
	tprmRoutes.GET("/third-party-register/:id", thirdPartyRegisterController.GetThirdPartyRegisterDetails)
	tprmRoutes.POST("/third-party-register", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyRegisterController.CreateThirdPartyRegister)
	tprmRoutes.PUT("/third-party-register/:id", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyRegisterController.UpdateThirdPartyRegister)
	tprmRoutes.POST("/third-party-register/delete", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyRegisterController.DeleteThirdPartyRegister)

	tprmRoutes.GET("/third-party-question-categories", thirdPartyQuestionCategoryController.GetAllThirdPartyQuestionCategories)
	tprmRoutes.GET("/third-party-question-categories/:id", thirdPartyQuestionCategoryController.GetThirdPartyQuestionCategoryDetails)
	tprmRoutes.POST("/third-party-question-categories", thirdPartyQuestionCategoryController.CreateThirdPartyQuestionCategory)
	tprmRoutes.PUT("/third-party-question-categories/:id", thirdPartyQuestionCategoryController.UpdateThirdPartyQuestionCategory)
	tprmRoutes.POST("/third-party-question-categories/delete", thirdPartyQuestionCategoryController.DeleteThirdPartyQuestionCategories)
	tprmRoutes.PUT("/third-party-question-categories/activate", thirdPartyQuestionCategoryController.ActivateThirdPartyQuestionCategories)
	tprmRoutes.PUT("/third-party-question-categories/deactivate", thirdPartyQuestionCategoryController.DeactivateThirdPartyQuestionCategories)

	tprmRoutes.GET("/third-party-question-types", thirdPartyQuestionTypeController.GetAllThirdPartyQuestionTypes)
	tprmRoutes.GET("/third-party-question-types/:id", thirdPartyQuestionTypeController.GetThirdPartyQuestionTypeDetails)

	// === Third Party Question Library ===
	tprmRoutes.GET("/third-party-question-libraries", thirdPartyQuestionLibraryController.GetAllThirdPartyQuestionLibraries)
	tprmRoutes.GET("/third-party-question-libraries/:id", thirdPartyQuestionLibraryController.GetThirdPartyQuestionLibraryDetails)
	tprmRoutes.POST("/third-party-question-libraries", thirdPartyQuestionLibraryController.CreateThirdPartyQuestionLibrary)
	tprmRoutes.PUT("/third-party-question-libraries/:id", thirdPartyQuestionLibraryController.UpdateThirdPartyQuestionLibrary)
	tprmRoutes.POST("/third-party-question-libraries/delete", thirdPartyQuestionLibraryController.DeleteThirdPartyQuestionLibraries)
	tprmRoutes.GET("/third-party-question-libraries/export", thirdPartyQuestionLibraryController.GetThirdPartyQuestionLibraryExport)

	tprmRoutes.GET("/third-party-question-libraries/template", middlewares.RolePermissionMiddleware(db, "question_library", "view"), thirdPartyQuestionLibraryController.GenerateQuestionLibraryTemplate)
	tprmRoutes.POST("/third-party-question-libraries/template/validate", middlewares.RolePermissionMiddleware(db, "question_library", "modify"), thirdPartyQuestionLibraryController.ValidateQuestionLibraryTemplate)
	tprmRoutes.POST("/third-party-question-libraries/import/validate", middlewares.RolePermissionMiddleware(db, "question_library", "modify"), thirdPartyQuestionLibraryController.ValidateImportQuestionLibraries)
	tprmRoutes.POST("/third-party-question-libraries/import", middlewares.RolePermissionMiddleware(db, "question_library", "modify"), thirdPartyQuestionLibraryController.ImportQuestionLibraries)

	tprmRoutes.GET("/question-sets", questionSetController.GetAllQuestionSets)
	tprmRoutes.GET("/question-sets/:id", questionSetController.GetQuestionSetDetails)
	tprmRoutes.POST("/question-sets", questionSetController.CreateQuestionSet)
	tprmRoutes.PUT("/question-sets/:id", questionSetController.UpdateQuestionSet)
	tprmRoutes.POST("/question-sets/delete", questionSetController.DeleteQuestionSets)
	tprmRoutes.POST("/question-sets/:id/questions", questionSetController.AddQuestionsToSet)
	tprmRoutes.DELETE("/question-sets/:id/questions/:questionId", questionSetController.DeleteQuestionFromSet)
	tprmRoutes.PUT("/question-sets/:id/sections/reorder", questionSetController.ReorderSection)
	tprmRoutes.PUT("/question-sets/:id/sections/:sectionId/questions/reorder", questionSetController.ReorderQuestion)
	tprmRoutes.GET("/question-sets/:id/export", questionSetController.GetQuestionSetExport)

	// assessment status
	tprmRoutes.GET("/third-party-assessment-statuses", thirdPartyAssessmentStatusController.GetAllThirdPartyAssessmentStatuses)
	tprmRoutes.GET("/third-party-assessment-statuses/:id", thirdPartyAssessmentStatusController.GetThirdPartyAssessmentStatusDetails)

	// === Third Party Assessment ===
	tprmRoutes.GET("/third-party-assessments", thirdPartyAssessmentController.GetAllThirdPartyAssessment)
	tprmRoutes.GET("/third-party-assessments/:assessmentID", thirdPartyAssessmentController.GetThirdPartyAssessmentDetails)
	tprmRoutes.POST("/third-party-assessments", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.CreateThirdPartyAssessment)
	tprmRoutes.PUT("/third-party-assessments/:assessmentID", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.UpdateThirdPartyAssessment)
	tprmRoutes.POST("/third-party-assessments/delete", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.DeleteThirdPartyAssessment)
	tprmRoutes.POST("/third-party-assessments/:assessmentID/share-link", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.ShareLink)
	tprmRoutes.GET("/third-party-assessments/:assessmentID/workflow-list", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentWorkflowController.GetWorkflowList)
	tprmRoutes.PUT("/third-party-assessments/:assessmentID/submit", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentReviewController.SubmitAssessment)
	tprmRoutes.PUT("/third-party-assessments/:assessmentID/approve", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentReviewController.ApproveAssessment)
	tprmRoutes.PUT("/third-party-assessments/:assessmentID/reject", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentReviewController.RejectAssessment)
	tprmRoutes.PUT("/third-party-assessments/:assessmentID/reopen", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentReviewController.ReopenAssessment)
	tprmRoutes.GET("/third-party-assessments/:assessmentID/export", middlewares.RolePermissionMiddleware(db, "assessment", "view"), thirdPartyAssessmentController.GetTPRMAssessmentExport)
	tprmRoutes.POST("/third-party-assessments/:assessmentID/template/validate", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.ValidateTPRMAssessmentTemplate)
	tprmRoutes.POST("/third-party-assessments/:assessmentID/import/validate", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.ValidateImportTPRMAssessment)
	tprmRoutes.POST("/third-party-assessments/:assessmentID/import", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.ImportTPRMAssessment)
	tprmRoutes.PUT("/third-party-assessments/:assessmentID/mark-as-send", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.MarkAsSend)
	tprmRoutes.GET("/third-party-assessments/:assessmentID/answer-summary", thirdPartyAssessmentController.GetThirdPartyAssessmentAnswerSummary)
	tprmRoutes.GET("/third-party-assessments/:assessmentID/export-pdf", middlewares.RolePermissionMiddleware(db, "assessment", "view"), thirdPartyAssessmentController.GetTPRMAssessmentReportExport)
	tprmRoutes.PUT("/third-party-assessments/:assessmentID/section-questions/:questionId/update-answer", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.UpdateAnswer)
	tprmRoutes.GET("/third-party-assessments/:assessmentID/sections/:sectionId", middlewares.RolePermissionMiddleware(db, "assessment", "view"), thirdPartyAssessmentController.GetAssessmentSectionDetails)
	tprmRoutes.GET("/third-party-assessments/:assessmentID/sections/:sectionId/questions/:questionId", middlewares.RolePermissionMiddleware(db, "assessment", "view"), thirdPartyAssessmentController.GetAssessmentSectionQuestionDetails)

	tprmRoutes.GET("/third-party-assessments/:assessmentID/section-questions/:questionId/answer/:answerId/documents", middlewares.RolePermissionMiddleware(db, "assessment", "view"), thirdPartyAssessmentController.GetAllAnswerDocuments)
	tprmRoutes.POST("/third-party-assessments/:assessmentID/section-questions/:questionId/answer/:answerId/documents", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.CreateAnswerDocument)
	tprmRoutes.DELETE("/third-party-assessments/:assessmentID/section-questions/:questionId/answer/:answerId/documents/:documentId", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), thirdPartyAssessmentController.DeleteAnswerDocument)
	unauthenticatedRoutes.GET("/third-party-assessments/:assessmentID/section-questions/:questionId/answer/:answerId/document-download", thirdPartyAssessmentController.DownloadAnswerDocument)
	unauthenticatedRoutes.GET("/third-party-assessments/:assessmentID/section-questions/:questionId/answer/:answerId/document-preview", thirdPartyAssessmentController.GetAnswerDocumentPreview)

	tprmRoutes.GET("/third-party-register/:id/processes", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllThirdPartyProcesses)
	tprmRoutes.POST("/third-party-register/:id/processes", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.CreateThirdPartyProcess)
	tprmRoutes.DELETE("/third-party-register/:id/processes/:processId", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.DeleteThirdPartyProcess)

	tprmRoutes.GET("/third-party-register/:id/strategic-objectives", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllThirdPartyStrategicObjectives)
	tprmRoutes.POST("/third-party-register/:id/strategic-objectives", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.CreateThirdPartyStrategicObjective)
	tprmRoutes.DELETE("/third-party-register/:id/strategic-objectives/:strategicObjectiveId", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.DeleteThirdPartyStrategicObjective)

	tprmRoutes.GET("/third-party-register/:id/assets", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllThirdPartyAssets)
	tprmRoutes.POST("/third-party-register/:id/assets", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.CreateThirdPartyAsset)
	tprmRoutes.DELETE("/third-party-register/:id/assets/:assetId", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.DeleteThirdPartyAsset)

	tprmRoutes.GET("/third-party-register/:id/incidents", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllThirdPartyIncidents)
	tprmRoutes.POST("/third-party-register/:id/incidents", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.CreateThirdPartyIncident)
	tprmRoutes.DELETE("/third-party-register/:id/incidents/:incidentId", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.DeleteThirdPartyIncident)

	tprmRoutes.GET("/third-party-register/:id/documents", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllThirdPartyDocuments)
	tprmRoutes.POST("/third-party-register/:id/documents", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.CreateThirdPartyDocument)
	tprmRoutes.DELETE("/third-party-register/:id/documents/:documentId", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.DeleteThirdPartyDocument)

	tprmRoutes.GET("/third-party-register/:id/controls", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllThirdPartyControls)
	tprmRoutes.POST("/third-party-register/:id/controls", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.CreateThirdPartyControl)
	tprmRoutes.DELETE("/third-party-register/:id/controls/:controlId", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.DeleteThirdPartyControl)

	tprmRoutes.GET("/third-party-register/:id/risks", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllThirdPartyRisks)
	tprmRoutes.POST("/third-party-register/:id/risks", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.CreateThirdPartyRisk)
	tprmRoutes.DELETE("/third-party-register/:id/risks/:riskId", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.DeleteThirdPartyRisk)

	tprmRoutes.GET("/third-party-register/:id/kpis", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllThirdPartyKPI)
	tprmRoutes.POST("/third-party-register/:id/kpis", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.CreateThirdPartyKPI)
	tprmRoutes.DELETE("/third-party-register/:id/kpis/:kpiRegisterId", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.DeleteThirdPartyKPI)

	tprmRoutes.GET("/third-party-register/:id/clauses", middlewares.RolePermissionMiddleware(db, "third_party_register", "view"), thirdPartyAssociatedItemController.GetAllClauses)
	tprmRoutes.PUT("/third-party-register/:id/clauses", middlewares.RolePermissionMiddleware(db, "third_party_register", "modify"), thirdPartyAssociatedItemController.UpdateClauses)

	tprmRoutes.GET("/tprm-settings", middlewares.RolePermissionMiddleware(db, "tprm_setting", "view"), tprmSettingController.GetTPRMSetting)
	tprmRoutes.PUT("/tprm-settings", middlewares.RolePermissionMiddleware(db, "tprm_setting", "modify"), tprmSettingController.UpdateTPRMSetting)

	tprmRoutes.GET("/third-party-assessments/:assessmentID/section-questions/:questionId/comments", middlewares.RolePermissionMiddleware(db, "assessment", "view"), tprmAssesmentSectionQuestionCommentController.GetAllComments)
	tprmRoutes.GET("/third-party-assessments/:assessmentID/comments", middlewares.RolePermissionMiddleware(db, "assessment", "view"), tprmAssesmentSectionQuestionCommentController.GetAllCommentsByAssessment)
	tprmRoutes.POST("/third-party-assessments/:assessmentID/comments", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), tprmAssesmentSectionQuestionCommentController.CreateAssessmentComment)
	tprmRoutes.PUT("/third-party-assessments/:assessmentID/comments/:commentId", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), tprmAssesmentSectionQuestionCommentController.UpdateComment)
	tprmRoutes.POST("/third-party-assessments/:assessmentID/comments/delete", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), tprmAssesmentSectionQuestionCommentController.DeleteComments)
	tprmRoutes.GET("/third-party-assessments/:assessmentID/section-questions/:questionId/comments/:commentId", middlewares.RolePermissionMiddleware(db, "assessment", "view"), tprmAssesmentSectionQuestionCommentController.GetCommentDetails)
	tprmRoutes.POST("/third-party-assessments/:assessmentID/section-questions/:questionId/comments", middlewares.RolePermissionMiddleware(db, "assessment", "modify"), tprmAssesmentSectionQuestionCommentController.CreateComment)

	tprmRoutes.GET("/tprm-compliance-responses", tprmComplianceResponseController.GetAllTPRMComplianceResponses)
	tprmRoutes.GET("/tprm-compliance-responses/:id", tprmComplianceResponseController.GetTPRMComplianceResponseDetails)
	tprmRoutes.POST("/tprm-compliance-responses", middlewares.RolePermissionMiddleware(db, "tprm_setting", "modify"), tprmComplianceResponseController.CreateTPRMComplianceResponse)
	tprmRoutes.PUT("/tprm-compliance-responses/:id", middlewares.RolePermissionMiddleware(db, "tprm_setting", "modify"), tprmComplianceResponseController.UpdateTPRMComplianceResponse)
	tprmRoutes.POST("/tprm-compliance-responses/delete", middlewares.RolePermissionMiddleware(db, "tprm_setting", "modify"), tprmComplianceResponseController.DeleteTPRMComplianceResponses)
	tprmRoutes.PUT("/tprm-compliance-responses/activate", middlewares.RolePermissionMiddleware(db, "tprm_setting", "modify"), tprmComplianceResponseController.ActivateTPRMComplianceResponses)
	tprmRoutes.PUT("/tprm-compliance-responses/deactivate", middlewares.RolePermissionMiddleware(db, "tprm_setting", "modify"), tprmComplianceResponseController.DeactivateTPRMComplianceResponses)

	tprmRoutes.GET("/tprm/third-party-insights", tprmDashboardController.GetThirdPartyInsights)
	tprmRoutes.GET("/tprm/assessment-insights", tprmDashboardController.GetAssessmentInsights)
	tprmRoutes.GET("/tprm/system-statistics", tprmDashboardController.GetTPRMSystemStatistics)
	tprmRoutes.GET("/tprm/ongoing-assessments", tprmDashboardController.GetOngoingAssessments)
	tprmRoutes.GET("/tprm/third-party-service-expiry", tprmDashboardController.GetThirdPartyServiceExpiry)
	tprmRoutes.GET("/tprm/assessment-by-template-type", tprmDashboardController.GetAssessmentByTemplateType)
	tprmRoutes.GET("/tprm/third-party-by-type", tprmDashboardController.GetThirdPartyByType)
	tprmRoutes.GET("/tprm/assessment-risk-monthly-trend", tprmDashboardController.GetAssessmentRiskMonthlyTrend)
}
