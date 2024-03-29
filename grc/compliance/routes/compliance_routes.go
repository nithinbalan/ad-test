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

	reviewFrequencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	reviewFrequencySvc := masterServices.NewReviewFrequencyService(reviewFrequencyRepo)

	severityLevelRepo := masterRepositories.NewSeverityLevelRepository(db)
	severityLevelSvc := masterServices.NewSeverityLevelService(severityLevelRepo)

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
		complianceRequirementCategorySvc,
		complianceStatusSvc,
		severityLevelSvc,
	)
	complianceRequirementController := controllers.NewComplianceRequirementController(complianceRequirementSvc)

	SLACategoryRepo := complianceRepositories.NewSLACategoryRepository(db)
	SLACategorySvc := complianceServices.NewSLACategoryService(SLACategoryRepo, activityLogSvc)
	SLACategoryController := controllers.NewSLACategoryController(SLACategorySvc)

	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
	complianceRoutes.GET("/compliance-categories/:id", complianceCategoryController.GetComplianceCategoryDetails)
	complianceRoutes.POST("/compliance-categories", complianceCategoryController.CreateComplianceCategory)
	complianceRoutes.PUT("/compliance-categories/:id", complianceCategoryController.UpdateComplianceCategory)
	complianceRoutes.POST("/compliance-categories/delete", complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/activate", complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/deactivate", complianceCategoryController.DeactivateComplianceCategories)

	complianceRoutes.GET("/compliance-statuses", complianceStatusController.GetAllComplianceStatuses)
	complianceRoutes.GET("/compliance-statuses/:id", complianceStatusController.GetComplianceStatusDetails)

	complianceRoutes.GET("/compliance-requirement-categories", complianceRequirementCategoryController.GetAllComplianceRequirementCategories)
	complianceRoutes.GET("/compliance-requirement-categories/:id", complianceRequirementCategoryController.GetComplianceRequirementCategoryDetails)
	complianceRoutes.POST("/compliance-requirement-categories", complianceRequirementCategoryController.CreateComplianceRequirementCategory)
	complianceRoutes.PUT("/compliance-requirement-categories/:id", complianceRequirementCategoryController.UpdateComplianceRequirementCategory)
	complianceRoutes.POST("/compliance-requirement-categories/delete", complianceRequirementCategoryController.DeleteComplianceRequirementCategories)
	complianceRoutes.POST("/compliance-requirement-categories/activate", complianceRequirementCategoryController.ActivateComplianceRequirementCategories)
	complianceRoutes.POST("/compliance-requirement-categories/deactivate", complianceRequirementCategoryController.DeactivateComplianceRequirementCategories)

	complianceRoutes.GET("/sla-categories/:id", SLACategoryController.GetSLACategoryDetails)
	complianceRoutes.POST("/sla-categories", SLACategoryController.CreateCategory)
	complianceRoutes.GET("/sla-categories", SLACategoryController.GetAllSLACategories)
	complianceRoutes.PUT("/sla-categories/:id", SLACategoryController.UpdateCategory)
	complianceRoutes.POST("/sla-categories/delete", SLACategoryController.DeleteSLACategories)
	complianceRoutes.POST("/sla-categories/activate", SLACategoryController.ActivateSLACategories)
	complianceRoutes.POST("/sla-categories/deactivate", SLACategoryController.DeactivateSLACategories)

	complianceRoutes.GET("/compliance-register", complianceRegisterController.GetAllComplianceRegisters)
	complianceRoutes.GET("/compliance-register/:id", complianceRegisterController.GetComplianceRegisterDetails)
	complianceRoutes.POST("/compliance-register", complianceRegisterController.CreateComplianceRegister)
	complianceRoutes.PUT("/compliance-register/:id", complianceRegisterController.UpdateComplianceRegister)
	complianceRoutes.POST("/compliance-register/delete", complianceRegisterController.DeleteComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/archive", complianceRegisterController.ArchiveComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/restore", complianceRegisterController.RestoreComplianceRegisters)

	complianceRoutes.GET("/compliances/:complianceId/requirements", complianceRequirementController.GetAllComplianceRequirements)
	complianceRoutes.GET("/compliances/:complianceId/requirements/:id", complianceRequirementController.GetComplianceRequirementDetails)
	complianceRoutes.POST("/compliances/:complianceId/requirements", complianceRequirementController.CreateComplianceRequirement)
	complianceRoutes.PUT("/compliances/:complianceId/requirements/:id", complianceRequirementController.UpdateComplianceRequirement)
	complianceRoutes.POST("/compliances/:complianceId/requirements/delete", complianceRequirementController.DeleteComplianceRequirements)
}
