package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware())

	complianceCategoryController := app.ComplianceModule.ComplianceCategoryController
	complianceStatusController := app.ComplianceModule.ComplianceStatusController
	complianceReviewStatusController := app.ComplianceModule.ComplianceReviewStatusController
	complianceRegisterController := app.ComplianceModule.ComplianceRegisterController
	complianceAssociatedItemController := app.ComplianceModule.ComplianceAssociatedItemController
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

	//Compliance Category
	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
	complianceRoutes.GET("/compliance-categories/:id", complianceCategoryController.GetComplianceCategoryDetails)
	complianceRoutes.POST("/compliance-categories", complianceCategoryController.CreateComplianceCategory)
	complianceRoutes.PUT("/compliance-categories/:id", complianceCategoryController.UpdateComplianceCategory)
	complianceRoutes.POST("/compliance-categories/delete", complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/activate", complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/deactivate", complianceCategoryController.DeactivateComplianceCategories)

	//Compliance Status
	complianceRoutes.GET("/compliance-statuses", complianceStatusController.GetAllComplianceStatuses)
	complianceRoutes.GET("/compliance-statuses/:id", complianceStatusController.GetComplianceStatusDetails)

	//Compliance Requirement Category
	complianceRoutes.GET("/compliance-requirement-categories", complianceRequirementCategoryController.GetAllComplianceRequirementCategories)
	complianceRoutes.GET("/compliance-requirement-categories/:id", complianceRequirementCategoryController.GetComplianceRequirementCategoryDetails)
	complianceRoutes.POST("/compliance-requirement-categories", complianceRequirementCategoryController.CreateComplianceRequirementCategory)
	complianceRoutes.PUT("/compliance-requirement-categories/:id", complianceRequirementCategoryController.UpdateComplianceRequirementCategory)
	complianceRoutes.POST("/compliance-requirement-categories/delete", complianceRequirementCategoryController.DeleteComplianceRequirementCategories)
	complianceRoutes.PUT("/compliance-requirement-categories/activate", complianceRequirementCategoryController.ActivateComplianceRequirementCategories)
	complianceRoutes.PUT("/compliance-requirement-categories/deactivate", complianceRequirementCategoryController.DeactivateComplianceRequirementCategories)

	//SLA Category
	complianceRoutes.GET("/sla-categories/:id", SLACategoryController.GetSLACategoryDetails)
	complianceRoutes.POST("/sla-categories", SLACategoryController.CreateCategory)
	complianceRoutes.GET("/sla-categories", SLACategoryController.GetAllSLACategories)
	complianceRoutes.PUT("/sla-categories/:id", SLACategoryController.UpdateCategory)
	complianceRoutes.POST("/sla-categories/delete", SLACategoryController.DeleteSLACategories)
	complianceRoutes.PUT("/sla-categories/activate", SLACategoryController.ActivateSLACategories)
	complianceRoutes.PUT("/sla-categories/deactivate", SLACategoryController.DeactivateSLACategories)

	//Compliance Review Statuses
	complianceRoutes.GET("/compliance-review-statuses", complianceReviewStatusController.GetAllComplianceReviewStatuses)
	complianceRoutes.GET("/compliance-review-statuses/:id", complianceReviewStatusController.GetComplianceReviewStatusDetails)

	//Compliance Register
	complianceRoutes.GET("/compliance-register", complianceRegisterController.GetAllComplianceRegisters)
	complianceRoutes.GET("/compliance-register/:id", complianceRegisterController.GetComplianceRegisterDetails)
	complianceRoutes.POST("/compliance-register", complianceRegisterController.CreateComplianceRegister)
	complianceRoutes.PUT("/compliance-register/:id", complianceRegisterController.UpdateComplianceRegister)
	complianceRoutes.POST("/compliance-register/delete", complianceRegisterController.DeleteComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/archive", complianceRegisterController.ArchiveComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/restore", complianceRegisterController.RestoreComplianceRegisters)

	//Compliance Associated Items
	complianceRoutes.GET("compliances/:complianceId/controls", complianceAssociatedItemController.GetAllComplianceControls)

	//SLA Register
	complianceRoutes.GET("/sla-register", SLARegisterController.GetAllSLARegisters)
	complianceRoutes.GET("/sla-register/:id", SLARegisterController.GetSLARegisterDetails)
	complianceRoutes.POST("/sla-register", SLARegisterController.CreateSLARegister)
	complianceRoutes.PUT("/sla-register/:id", SLARegisterController.UpdateSLARegister)
	complianceRoutes.POST("/sla-register/delete", SLARegisterController.DeleteSLARegisters)
	complianceRoutes.PUT("/sla-register/archive", SLARegisterController.ArchiveSLARegisters)
	complianceRoutes.PUT("/sla-register/restore", SLARegisterController.RestoreSLARegisters)

	//SLA Renew
	complianceRoutes.GET("/slas/:slaId/renew-history", SLARenewController.GetSLARenewHistory)
	complianceRoutes.PUT("/slas/:slaId/renew", SLARenewController.Renew)
	complianceRoutes.PUT("/slas/:slaId/revert", SLARenewController.Revert)

	//Compliance Requirement
	complianceRoutes.GET("/compliance-requirements", complianceRequirementController.GetAllComplianceRequirements)
	complianceRoutes.GET("/compliances/:complianceId/requirements/:requirementId", complianceRequirementController.GetComplianceRequirementDetails)
	complianceRoutes.POST("/compliances/:complianceId/requirements", complianceRequirementController.CreateComplianceRequirement)
	complianceRoutes.PUT("/compliances/:complianceId/requirements/:requirementId", complianceRequirementController.UpdateComplianceRequirement)
	complianceRoutes.POST("/compliances/:complianceId/requirements/delete", complianceRequirementController.DeleteComplianceRequirements)

	//Compliance Action Plan
	complianceRoutes.GET("/compliances/:complianceId/action-plans", complianceActionPlanController.GetAllComplianceActionPlans)
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId", complianceActionPlanController.GetComplianceActionPlanDetails)
	complianceRoutes.POST("/compliances/:complianceId/action-plans", complianceActionPlanController.CreateComplianceActionPlan)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId", complianceActionPlanController.UpdateComplianceActionPlan)
	complianceRoutes.POST("/compliances/:complianceId/action-plans/delete", complianceActionPlanController.DeleteComplianceActionPlans)

	//Compliance Status Update
	complianceRoutes.GET("/compliances/:complianceId/status-updates", complianceStatusUpdateController.GetAllComplianceStatusUpdates)
	complianceRoutes.GET("/compliances/:complianceId/status-updates/:id", complianceStatusUpdateController.GetComplianceStatusUpdateDetails)
	complianceRoutes.POST("/compliances/:complianceId/status-updates", complianceStatusUpdateController.CreateComplianceStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/status-updates/:id", complianceStatusUpdateController.UpdateComplianceStatusUpdate)

	//Compliance Requirement Status Update
	complianceRoutes.GET("/compliances/:complianceId/requirements/:requirementId/status-updates", complianceRequirementStatusUpdateController.GetAllComplianceRequirementStatusUpdates)
	complianceRoutes.POST("/compliances/:complianceId/requirements/:requirementId/status-updates", complianceRequirementStatusUpdateController.CreateComplianceRequirementStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/requirements/:requirementId/status-updates/:id", complianceRequirementStatusUpdateController.UpdateComplianceRequirementStatusUpdate)

	//Compliance Action Plan
	complianceRoutes.GET("/compliances/:complianceId/action-plans/:actionPlanId/status-updates", complianceActionPlanStatusUpdateController.GetAllComplianceActionPlanStatusUpdates)
	complianceRoutes.POST("/compliances/:complianceId/action-plans/:actionPlanId/status-updates", complianceActionPlanStatusUpdateController.CreateComplianceActionPlanStatusUpdate)
	complianceRoutes.PUT("/compliances/:complianceId/action-plans/:actionPlanId/status-updates/:id", complianceActionPlanStatusUpdateController.UpdateComplianceActionPlanStatusUpdate)

	//Compliance Review
	complianceRoutes.GET("/compliances/:complianceId/reviews", complianceReviewController.GetAllComplianceReviews)
	complianceRoutes.POST("/compliances/:complianceId/reviews", complianceReviewController.CreateComplianceReview)
	complianceRoutes.PUT("/compliances/:complianceId/reviews/:id", complianceReviewController.UpdateComplianceReview)

	//SLA Check Question
	complianceRoutes.GET("/slas/:slaId/checks/:checkId/question-groups", SLACheckQuestionController.GetAllSLACheckQuestionGroups)
	complianceRoutes.PUT("/slas/:slaId/checks/:checkId/question-groups/:id", SLACheckQuestionController.UpdateSLACheckQuestionGroup)
	complianceRoutes.DELETE("/slas/:slaId/checks/:checkId/question-groups/:id", SLACheckQuestionController.DeleteSLACheckQuestionGroup)
	complianceRoutes.GET("/slas/:slaId/checks/:checkId/question-groups/:id/questions", SLACheckQuestionController.GetSLACheckQuestionsByGroup)
	complianceRoutes.POST("/slas/:slaId/checks/:checkId/choose-questions", SLACheckQuestionController.ChooseSLACheckQuestions)
	complianceRoutes.POST("/slas/:slaId/checks/:checkId/add-question", SLACheckQuestionController.CreateSLACheckQuestion)
	complianceRoutes.PUT("/slas/:slaId/checks/:checkId/questions/:id", SLACheckQuestionController.UpdateSLACheckQuestion)
	complianceRoutes.DELETE("/slas/:slaId/checks/:checkId/questions/:id", SLACheckQuestionController.DeleteSLACheckQuestion)

	//SLA Answer
	complianceRoutes.POST("/slas/:slaId/questions/:questionId/answers", SLACheckAnswerController.CreateSLACheckAnswer)
	complianceRoutes.PUT("/slas/:slaId/questions/:questionId/answers/:answerId", SLACheckAnswerController.UpdateSLACheckAnswer)
	complianceRoutes.DELETE("/slas/:slaId/questions/:questionId/answers/:answerId", SLACheckAnswerController.DeleteSLACheckAnswer)
	complianceRoutes.GET("/slas/:slaId/assessment-summary", SLACheckAnswerController.AssessmentSummary)

	//SLA Action plan
	complianceRoutes.GET("/slas/:slaId/action-plans", slaActionPlanController.GetAllSLAActionPlans)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId", slaActionPlanController.GetSLAActionPlanDetails)
	complianceRoutes.POST("/slas/:slaId/action-plans", slaActionPlanController.CreateSLAActionPlan)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId", slaActionPlanController.UpdateSLAActionPlan)
	complianceRoutes.POST("/slas/:slaId/action-plans/delete", slaActionPlanController.DeleteSLAActionPlans)

	// SLA Action Plan
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId/status-updates", slaActionPlanStatusUpdateController.GetAllSLAActionPlanStatusUpdates)
	complianceRoutes.POST("/slas/:slaId/action-plans/:actionPlanId/status-updates", slaActionPlanStatusUpdateController.CreateSLAActionPlanStatusUpdate)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId/status-updates/:id", slaActionPlanStatusUpdateController.UpdateSLAActionPlanStatusUpdate)

}
