package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware())

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	complianceCategoryController := app.ComplianceModule.ComplianceCategoryController
	complianceStatusController := app.ComplianceModule.ComplianceStatusController
	complianceReviewStatusController := app.ComplianceModule.ComplianceReviewStatusController
	complianceRegisterController := app.ComplianceModule.ComplianceRegisterController
	complianceAssociatedItemController := app.ComplianceModule.ComplianceAssociatedItemController
	complianceRequirementAssociatedItemController := app.ComplianceModule.ComplianceRequirementAssociatedItemController
	complianceRequirementCategoryController := app.ComplianceModule.ComplianceRequirementCategoryController
	complianceRequirementController := app.ComplianceModule.ComplianceRequirementController
	complianceActionPlanController := app.ComplianceModule.ComplianceActionPlanController
	complianceStatusUpdateController := app.ComplianceModule.ComplianceStatusUpdateController
	complianceRequirementStatusUpdateController := app.ComplianceModule.ComplianceRequirementStatusUpdateController
	complianceActionPlanStatusUpdateController := app.ComplianceModule.ComplianceActionPlanStatusUpdateController
	complianceReviewController := app.ComplianceModule.ComplianceReviewController
	SLACategoryController := app.ComplianceModule.SLACategoryController
	SLARegisterController := app.ComplianceModule.SLARegisterController
	SLACheckQuestionController := app.ComplianceModule.SLACheckQuestionController
	SLACheckAnswerController := app.ComplianceModule.SLACheckAnswerController
	SLARenewController := app.ComplianceModule.SLARenewController
	slaActionPlanController := app.ComplianceModule.SLAActionPlanController
	slaActionPlanStatusUpdateController := app.ComplianceModule.SLAActionPlanStatusUpdateController
	slaAssociatedItemController := app.ComplianceModule.SLAAssociatedItemController

	//Compliance Category
	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
	complianceRoutes.GET("/compliance-categories/:id", complianceCategoryController.GetComplianceCategoryDetails)
	complianceRoutes.POST("/compliance-categories", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.CreateComplianceCategory)
	complianceRoutes.PUT("/compliance-categories/:id", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.UpdateComplianceCategory)
	complianceRoutes.POST("/compliance-categories/delete", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/activate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/deactivate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceCategoryController.DeactivateComplianceCategories)

	//Compliance Status
	complianceRoutes.GET("/compliance-statuses", complianceStatusController.GetAllComplianceStatuses)
	complianceRoutes.GET("/compliance-statuses/:id", complianceStatusController.GetComplianceStatusDetails)

	//Compliance Requirement Category
	complianceRoutes.GET("/compliance-requirement-categories", complianceRequirementCategoryController.GetAllComplianceRequirementCategories)
	complianceRoutes.GET("/compliance-requirement-categories/:id", complianceRequirementCategoryController.GetComplianceRequirementCategoryDetails)
	complianceRoutes.POST("/compliance-requirement-categories", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.CreateComplianceRequirementCategory)
	complianceRoutes.PUT("/compliance-requirement-categories/:id", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.UpdateComplianceRequirementCategory)
	complianceRoutes.POST("/compliance-requirement-categories/delete", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.DeleteComplianceRequirementCategories)
	complianceRoutes.PUT("/compliance-requirement-categories/activate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.ActivateComplianceRequirementCategories)
	complianceRoutes.PUT("/compliance-requirement-categories/deactivate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), complianceRequirementCategoryController.DeactivateComplianceRequirementCategories)

	//SLA Category
	complianceRoutes.GET("/sla-categories/:id", SLACategoryController.GetSLACategoryDetails)
	complianceRoutes.POST("/sla-categories", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.CreateCategory)
	complianceRoutes.GET("/sla-categories", SLACategoryController.GetAllSLACategories)
	complianceRoutes.PUT("/sla-categories/:id", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.UpdateCategory)
	complianceRoutes.POST("/sla-categories/delete", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.DeleteSLACategories)
	complianceRoutes.PUT("/sla-categories/activate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.ActivateSLACategories)
	complianceRoutes.PUT("/sla-categories/deactivate", middlewares.RolePermissionMiddleware(db, "compliance_setting", "modify"), SLACategoryController.DeactivateSLACategories)

	//Compliance Review Statuses
	complianceRoutes.GET("/compliance-review-statuses", complianceReviewStatusController.GetAllComplianceReviewStatuses)
	complianceRoutes.GET("/compliance-review-statuses/:id", complianceReviewStatusController.GetComplianceReviewStatusDetails)

	//Compliance Register
	unauthenticatedRoutes.GET("/compliance-register/:id/document-preview", complianceRegisterController.GetComplianceRegisterDocumentPreview)
	unauthenticatedRoutes.GET("/compliance-register/:id/document-download", complianceRegisterController.DownloadComplianceRegisterDocument)
	complianceRoutes.GET("/compliance-register/export", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRegisterController.GetComplianceRegisterExport)
	complianceRoutes.GET("/compliance-register", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRegisterController.GetAllComplianceRegisters)
	complianceRoutes.GET("/compliance-register/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRegisterController.GetComplianceRegisterDetails)
	complianceRoutes.POST("/compliance-register", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.CreateComplianceRegister)
	complianceRoutes.PUT("/compliance-register/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.UpdateComplianceRegister)
	complianceRoutes.POST("/compliance-register/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.DeleteComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/archive", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.ArchiveComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/restore", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRegisterController.RestoreComplianceRegisters)

	//Compliance Associated Items
	complianceRoutes.GET("/compliances/:complianceId/controls", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceControls)
	complianceRoutes.GET("/compliances/:complianceId/risks", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceRisks)
	complianceRoutes.GET("/compliances/:complianceId/incidents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceIncidents)
	complianceRoutes.GET("/compliances/:complianceId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllComplianceStrategicObjectives)
	complianceRoutes.GET("/compliances/:complianceId/associated-item-summary", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetAllAssociatedCounts)
	complianceRoutes.GET("/compliances/:complianceId/requirements-status-summary", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceAssociatedItemController.GetComplianceRequirementCountByStatus)

	//SLA Register
	unauthenticatedRoutes.GET("/sla-register/:id/document-download", SLARegisterController.DownloadSLARegisterDocument)
	unauthenticatedRoutes.GET("/sla-register/:id/document-preview", SLARegisterController.GetSLARegisterDocumentPreview)
	complianceRoutes.GET("/sla-register/export", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARegisterController.GetSLARegisterExport)
	complianceRoutes.GET("/sla-register", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARegisterController.GetAllSLARegisters)
	complianceRoutes.GET("/sla-register/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARegisterController.GetSLARegisterDetails)
	complianceRoutes.POST("/sla-register", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.CreateSLARegister)
	complianceRoutes.PUT("/sla-register/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.UpdateSLARegister)
	complianceRoutes.POST("/sla-register/delete", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.DeleteSLARegisters)
	complianceRoutes.PUT("/sla-register/archive", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.ArchiveSLARegisters)
	complianceRoutes.PUT("/sla-register/restore", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARegisterController.RestoreSLARegisters)
	complianceRoutes.GET("/slas/:slaId/assessment-status-summary", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARegisterController.AssessmentStatusSummary)

	//SLA Renew
	complianceRoutes.GET("/slas/:slaId/renew-history", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLARenewController.GetSLARenewHistory)
	complianceRoutes.PUT("/slas/:slaId/renew", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARenewController.Renew)
	complianceRoutes.PUT("/slas/:slaId/revert", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLARenewController.Revert)

	//Compliance Requirement
	complianceRoutes.GET("/compliance-requirements", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementController.GetAllComplianceRequirements)
	complianceRoutes.GET("/compliances/:complianceId/requirements/export", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementController.GetComplianceRequirementsExport)
	complianceRoutes.GET("/compliances/:complianceId/requirements", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementController.GetComplianceRequirementsByID)
	complianceRoutes.GET("/compliances/:complianceId/requirements/:requirementId", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementController.GetComplianceRequirementDetails)
	complianceRoutes.POST("/compliances/:complianceId/requirements", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.CreateComplianceRequirement)
	complianceRoutes.PUT("/compliances/:complianceId/requirements/:requirementId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.UpdateComplianceRequirement)
	complianceRoutes.POST("/compliances/:complianceId/requirements/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementController.DeleteComplianceRequirements)

	//Compliance Requirements Associated Items
	complianceRoutes.GET("/compliance-requirements/:requirementId/controls", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementControls)
	complianceRoutes.GET("/compliance-requirements/:requirementId/strategic-objectives", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementStrategicObjectives)
	complianceRoutes.GET("/compliance-requirements/:requirementId/risks", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementRisks)
	complianceRoutes.GET("/compliance-requirements/:requirementId/incidents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementIncidents)
	complianceRoutes.GET("/compliance-requirements/:requirementId/documents", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementAssociatedItemController.GetAllComplianceRequirementDocuments)

	//Compliance Action Plan
	complianceRoutes.GET("/compliance-action-plans/export", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanController.GetComplianceActionPlanExport)
	complianceRoutes.GET("/compliance-action-plans", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanController.GetAllComplianceActionPlans)
	complianceRoutes.GET("/compliances/:complianceId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanController.GetAllComplianceActionPlans)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanController.GetComplianceActionPlanDetails)
	complianceRoutes.POST("/compliances/:complianceId/action-plans", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceActionPlanController.CreateComplianceActionPlan)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceActionPlanController.UpdateComplianceActionPlan)
	complianceRoutes.POST("/compliances/:complianceId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceActionPlanController.DeleteComplianceActionPlans)

	//Compliance Status Update
	unauthenticatedRoutes.GET("/compliances/:complianceId/status-updates/:id/document-download", complianceStatusUpdateController.DownloadComplianceStatusUpdateDocument)
	unauthenticatedRoutes.GET("/compliances/:complianceId/status-updates/:id/document-preview", complianceStatusUpdateController.GetComplianceStatusUpdateDocumentPreview)
	complianceRoutes.GET("/compliances/:complianceId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceStatusUpdateController.GetAllComplianceStatusUpdates)
	complianceRoutes.GET("/compliances/:complianceId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceStatusUpdateController.GetComplianceStatusUpdateDetails)
	complianceRoutes.POST("/compliances/:complianceId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceStatusUpdateController.CreateComplianceStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceStatusUpdateController.UpdateComplianceStatusUpdate)
	complianceRoutes.GET("/compliances/:complianceId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceStatusUpdateController.GetLatestComplianceStatusUpdate)

	//Compliance Requirement Status Update
	unauthenticatedRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates/:id/document-download", complianceRequirementStatusUpdateController.DownloadComplianceRequirementStatusUpdateDocument)
	unauthenticatedRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates/:id/document-preview", complianceRequirementStatusUpdateController.GetComplianceRequirementStatusUpdateDocumentPreview)
	complianceRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementStatusUpdateController.GetAllComplianceRequirementStatusUpdates)
	complianceRoutes.POST("/compliances/:complianceId/requirements/:requirementId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementStatusUpdateController.CreateComplianceRequirementStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/requirements/:requirementId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceRequirementStatusUpdateController.UpdateComplianceRequirementStatusUpdate)
	complianceRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceRequirementStatusUpdateController.GetLatestComplianceRequirementStatusUpdate)

	//Compliance Action Plan
	unauthenticatedRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/:id/document-download", complianceActionPlanStatusUpdateController.DownloadComplianceActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/:id/document-preview", complianceActionPlanStatusUpdateController.GetComplianceActionPlanStatusUpdateDocumentPreview)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanStatusUpdateController.GetAllComplianceActionPlanStatusUpdates)
	complianceRoutes.POST("/compliances/:complianceId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceActionPlanStatusUpdateController.CreateComplianceActionPlanStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceActionPlanStatusUpdateController.UpdateComplianceActionPlanStatusUpdate)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceActionPlanStatusUpdateController.GetLatestComplianceActionPlanStatusUpdate)

	//Compliance Review
	complianceRoutes.GET("/compliances/:complianceId/reviews", middlewares.RolePermissionMiddleware(db, "compliance_register", "view"), complianceReviewController.GetAllComplianceReviews)
	complianceRoutes.POST("/compliances/:complianceId/reviews", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceReviewController.CreateComplianceReview)
	complianceRoutes.PUT("/compliances/:complianceId/reviews/:id", middlewares.RolePermissionMiddleware(db, "compliance_register", "modify"), complianceReviewController.UpdateComplianceReview)

	//SLA Check Question
	complianceRoutes.GET("/slas/:slaId/checks/:checkId/question-groups", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckQuestionController.GetAllSLACheckQuestionGroups)
	complianceRoutes.GET("/slas/:slaId/checks/:checkId/question-groups/:groupId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckQuestionController.GetSLACheckQuestionGroupDetails)
	complianceRoutes.POST("/slas/:slaId/checks/:checkId/question-groups", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.CreateSLACheckQuestionGroup)
	complianceRoutes.PUT("/slas/:slaId/checks/:checkId/question-groups/:groupId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.UpdateSLACheckQuestionGroup)
	complianceRoutes.DELETE("/slas/:slaId/checks/:checkId/question-groups/:groupId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.DeleteSLACheckQuestionGroup)

	complianceRoutes.POST("/slas/:slaId/checks/:checkId/choose-questions", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.ChooseSLACheckQuestions)
	complianceRoutes.GET("/slas/:slaId/checks/:checkId/questions/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckQuestionController.GetSLACheckQuestionDetails)

	complianceRoutes.GET("/slas/:slaId/checks/:checkId/question-groups/:groupId/questions", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckQuestionController.GetSLACheckQuestionsByGroup)
	complianceRoutes.POST("/slas/:slaId/checks/:checkId/question-groups/:groupId/question", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.CreateSLACheckQuestion)
	complianceRoutes.PUT("/slas/:slaId/checks/:checkId/question-groups/:groupId/questions/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.UpdateSLACheckQuestion)
	complianceRoutes.DELETE("/slas/:slaId/checks/:checkId/question-groups/:groupId/questions/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckQuestionController.DeleteSLACheckQuestion)

	//SLA Answer
	complianceRoutes.POST("/slas/:slaId/questions/:questionId/answers", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckAnswerController.CreateSLACheckAnswer)
	complianceRoutes.PUT("/slas/:slaId/questions/:questionId/answers/:answerId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckAnswerController.UpdateSLACheckAnswer)
	complianceRoutes.DELETE("/slas/:slaId/questions/:questionId/answers/:answerId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), SLACheckAnswerController.DeleteSLACheckAnswer)
	complianceRoutes.GET("/slas/:slaId/assessment-summary", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), SLACheckAnswerController.AssessmentSummary)

	//SLA Action plan
	complianceRoutes.GET("/slas/:slaId/action-plans", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanController.GetAllSLAActionPlans)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanController.GetSLAActionPlanDetails)
	complianceRoutes.POST("/slas/:slaId/action-plans", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanController.CreateSLAActionPlan)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanController.UpdateSLAActionPlan)
	complianceRoutes.POST("/slas/:slaId/action-plans/delete", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanController.DeleteSLAActionPlans)
	complianceRoutes.GET("/slas/:slaId/action-plan-summary", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanController.GetSLAActionPlanCountByStatus)

	// SLA Action Plan
	unauthenticatedRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates/:id/document-download", slaActionPlanStatusUpdateController.DownloadSLAActionPlanStatusUpdateDocument)
	unauthenticatedRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates/:id/document-preview", slaActionPlanStatusUpdateController.GetSLAActionPlanStatusUpdateDocumentPreview)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanStatusUpdateController.GetAllSLAActionPlanStatusUpdates)
	complianceRoutes.POST("/slas/:slaId/action-plans/:actionPlanId/status-updates", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanStatusUpdateController.CreateSLAActionPlanStatusUpdate)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId/status-updates/:id", middlewares.RolePermissionMiddleware(db, "sla_contracts", "modify"), slaActionPlanStatusUpdateController.UpdateSLAActionPlanStatusUpdate)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates/latest", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaActionPlanStatusUpdateController.GetLatestSLAActionPlanStatusUpdate)

	complianceRoutes.GET("/slas/:slaId/processes", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllSLAProcesses)
	complianceRoutes.GET("/slas/:slaId/documents", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllSLADocuments)
	complianceRoutes.GET("/slas/:slaId/risks", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllSLARisks)
	complianceRoutes.GET("/slas/:slaId/incidents", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllSLAIncidents)
	complianceRoutes.GET("/slas/:slaId/associated-item-summary", middlewares.RolePermissionMiddleware(db, "sla_contracts", "view"), slaAssociatedItemController.GetAllAssociatedCounts)

}
