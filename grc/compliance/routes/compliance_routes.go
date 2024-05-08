package routes

import (
	"jethur/grc/compliance/controllers"
	complianceRepositories "jethur/grc/compliance/repositories"
	complianceServices "jethur/grc/compliance/services"
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	lifecycleStatusRepo := masterRepositories.NewLifecycleStatusRepository(db)
	lifecycleStatusSvc := masterServices.NewLifecycleStatusService(lifecycleStatusRepo)

	actionPlanStatusRepo := masterRepositories.NewActionPlanStatusRepository(db)
	actionPlanStatusSvc := masterServices.NewActionPlanStatusService(actionPlanStatusRepo)

	reviewFrequencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	reviewFrequencySvc := masterServices.NewReviewFrequencyService(reviewFrequencyRepo)

	severityLevelRepo := masterRepositories.NewSeverityLevelRepository(db)
	severityLevelSvc := masterServices.NewSeverityLevelService(severityLevelRepo)

	questionGroupRepository := masterRepositories.NewQuestionGroupRepository(db)
	questionGroupService := masterServices.NewQuestionGroupService(questionGroupRepository, activityLogSvc)

	checklistRepo := masterRepositories.NewChecklistRepository(db)
	checklistSvc := masterServices.NewChecklistService(checklistRepo, activityLogSvc)

	checklistQuestionRepo := masterRepositories.NewChecklistQuestionRepository(db)
	checklistQuestionSvc := masterServices.NewChecklistQuestionService(checklistQuestionRepo, activityLogSvc, questionGroupService, checklistSvc)

	complianceCategoryRepo := complianceRepositories.NewComplianceCategoryRepository(db)
	complianceCategorySvc := complianceServices.NewComplianceCategoryService(complianceCategoryRepo, activityLogSvc)
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategorySvc)

	complianceStatusRepo := complianceRepositories.NewComplianceStatusRepository(db)
	complianceStatusSvc := complianceServices.NewComplianceStatusService(complianceStatusRepo)
	complianceStatusController := controllers.NewComplianceStatusController(complianceStatusSvc)

	complianceRegisterRepo := complianceRepositories.NewComplianceRegisterRepository(db)
	complianceRegisterSvc := complianceServices.NewComplianceRegisterService(
		complianceRegisterRepo,
		activityLogSvc,
		lifecycleStatusSvc,
		complianceCategorySvc,
		complianceStatusSvc,
		reviewFrequencySvc,
		severityLevelSvc,
	)
	complianceRegisterController := controllers.NewComplianceRegisterController(complianceRegisterSvc)

	complianceRequirementCategoryRepo := complianceRepositories.NewComplianceRequirementCategoryRepository(db)
	complianceRequirementCategorySvc := complianceServices.NewComplianceRequirementCategoryService(complianceRequirementCategoryRepo, activityLogSvc)
	complianceRequirementCategoryController := controllers.NewComplianceRequirementCategoryController(complianceRequirementCategorySvc)

	complianceRequirementRepo := complianceRepositories.NewComplianceRequirementRepository(db)
	complianceRequirementSvc := complianceServices.NewComplianceRequirementService(
		complianceRequirementRepo,
		activityLogSvc,
		complianceRegisterSvc,
		complianceRequirementCategorySvc,
		complianceStatusSvc,
		severityLevelSvc,
	)
	complianceRequirementController := controllers.NewComplianceRequirementController(complianceRequirementSvc)

	complianceActionPlanRepo := complianceRepositories.NewComplianceActionPlanRepository(db)
	complianceActionPlanSvc := complianceServices.NewComplianceActionPlanService(
		complianceActionPlanRepo,
		activityLogSvc,
		complianceRegisterSvc,
		complianceRequirementSvc,
		actionPlanStatusSvc,
		severityLevelSvc,
	)
	complianceActionPlanController := controllers.NewComplianceActionPlanController(complianceActionPlanSvc)

	complianceStatusUpdateRepo := complianceRepositories.NewComplianceStatusUpdateRepository(db)
	complianceStatusUpdateSvc := complianceServices.NewComplianceStatusUpdateService(complianceStatusUpdateRepo, activityLogSvc, complianceRegisterSvc, complianceStatusSvc)
	complianceStatusUpdateController := controllers.NewComplianceStatusUpdateController(complianceStatusUpdateSvc)

	complianceRequirementStatusUpdateRepo := complianceRepositories.NewComplianceRequirementStatusUpdateRepository(db)
	complianceRequirementStatusUpdateSvc := complianceServices.NewComplianceRequirementStatusUpdateService(complianceRequirementStatusUpdateRepo, activityLogSvc, complianceRequirementSvc, complianceStatusSvc)
	complianceRequirementStatusUpdateController := controllers.NewComplianceRequirementStatusUpdateController(complianceRequirementStatusUpdateSvc)

	complianceActionPlanStatusUpdateRepo := complianceRepositories.NewComplianceActionPlanStatusUpdateRepository(db)
	complianceActionPlanStatusUpdateSvc := complianceServices.NewComplianceActionPlanStatusUpdateService(complianceActionPlanStatusUpdateRepo, activityLogSvc, complianceActionPlanSvc, actionPlanStatusSvc)
	complianceActionPlanStatusUpdateController := controllers.NewComplianceActionPlanStatusUpdateController(complianceActionPlanStatusUpdateSvc)

	complianceReviewRepo := complianceRepositories.NewComplianceReviewRepository(db)
	complianceReviewSvc := complianceServices.NewComplianceReviewService(complianceReviewRepo, activityLogSvc, complianceRegisterSvc)
	complianceReviewController := controllers.NewComplianceReviewController(complianceReviewSvc)

	SLACategoryRepo := complianceRepositories.NewSLACategoryRepository(db)
	SLACategorySvc := complianceServices.NewSLACategoryService(SLACategoryRepo, activityLogSvc)
	SLACategoryController := controllers.NewSLACategoryController(SLACategorySvc)

	SLARegisterRepo := complianceRepositories.NewSLARegisterRepository(db)
	slaRegisterSvc := complianceServices.NewSLARegisterService(SLARegisterRepo,
		SLACategorySvc,
		severityLevelSvc,
		lifecycleStatusSvc,
		activityLogSvc)
	SLARegisterController := controllers.NewSLARegisterController(slaRegisterSvc)

	SLACheckQuestionRepo := complianceRepositories.NewSLACheckQuestionRepository(db)
	slaCheckQuestionSvc := complianceServices.NewSLACheckQuestionService(SLACheckQuestionRepo,
		activityLogSvc,
		slaRegisterSvc,
		checklistQuestionSvc,
		questionGroupService)
	SLACheckQuestionController := controllers.NewSLACheckQuestionController(slaCheckQuestionSvc)

	SLACheckAnswerRepo := complianceRepositories.NewSLACheckAnswerRepository(db)
	SLACheckAnswerSvc := complianceServices.NewSLACheckAnswerService(SLACheckAnswerRepo,
		activityLogSvc,
		slaRegisterSvc,
		complianceStatusSvc,
		slaCheckQuestionSvc)
	SLACheckAnswerController := controllers.NewSLACheckAnswerController(SLACheckAnswerSvc)

	SLARenewRepo := complianceRepositories.NewSLARenewRepository(db)
	SLARenewSvc := complianceServices.NewSLARenewService(SLARenewRepo, SLARegisterRepo, lifecycleStatusSvc,
		activityLogSvc)
	SLARenewController := controllers.NewSLARenewController(SLARenewSvc)

	slaActionPlanRepo := complianceRepositories.NewSLAActionPlanRepository(db)
	slaActionPlanSvc := complianceServices.NewSLAActionPlanService(
		slaActionPlanRepo,
		activityLogSvc,
		slaRegisterSvc,
		slaCheckQuestionSvc,
		actionPlanStatusSvc,
		severityLevelSvc,
	)
	slaActionPlanController := controllers.NewSLAActionPlanController(slaActionPlanSvc)

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

	//Compliance Register
	complianceRoutes.GET("/compliance-register", complianceRegisterController.GetAllComplianceRegisters)
	complianceRoutes.GET("/compliance-register/:id", complianceRegisterController.GetComplianceRegisterDetails)
	complianceRoutes.POST("/compliance-register", complianceRegisterController.CreateComplianceRegister)
	complianceRoutes.PUT("/compliance-register/:id", complianceRegisterController.UpdateComplianceRegister)
	complianceRoutes.POST("/compliance-register/delete", complianceRegisterController.DeleteComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/archive", complianceRegisterController.ArchiveComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/restore", complianceRegisterController.RestoreComplianceRegisters)

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
	complianceRoutes.GET("/compliances/:complianceId/requirements", complianceRequirementController.GetAllComplianceRequirements)
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

	//SLA Action plan
	complianceRoutes.GET("/slas/:slaId/action-plans", slaActionPlanController.GetAllSLAActionPlans)
	complianceRoutes.GET("/slas/:slaId/action-plans/:actionPlanId", slaActionPlanController.GetSLAActionPlanDetails)
	complianceRoutes.POST("/slas/:slaId/action-plans", slaActionPlanController.CreateSLAActionPlan)
	complianceRoutes.PUT("/slas/:slaId/action-plans/:actionPlanId", slaActionPlanController.UpdateSLAActionPlan)
	complianceRoutes.POST("/slas/:slaId/action-plans/delete", slaActionPlanController.DeleteSLAActionPlans)

}
