package routes

import (
	"jethur/grc/compliance/controllers"
	complianceRepositories "jethur/grc/compliance/repositories"
	complianceServices "jethur/grc/compliance/services"
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	complianceRoutes := router.Group("/")

	complianceRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	//Compliance Category
	complianceCategoryRepo := complianceRepositories.NewComplianceCategoryRepository(db)
	complianceCategorySvc := complianceServices.NewComplianceCategoryService(complianceCategoryRepo, activityLogSvc)
	complianceCategoryController := controllers.NewComplianceCategoryController(complianceCategorySvc)

	complianceRoutes.GET("/compliance-categories", complianceCategoryController.GetAllComplianceCategories)
	complianceRoutes.GET("/compliance-categories/:id", complianceCategoryController.GetComplianceCategoryDetails)
	complianceRoutes.POST("/compliance-categories", complianceCategoryController.CreateComplianceCategory)
	complianceRoutes.PUT("/compliance-categories/:id", complianceCategoryController.UpdateComplianceCategory)
	complianceRoutes.POST("/compliance-categories/delete", complianceCategoryController.DeleteComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/activate", complianceCategoryController.ActivateComplianceCategories)
	complianceRoutes.PUT("/compliance-categories/deactivate", complianceCategoryController.DeactivateComplianceCategories)

	//Compliance Status
	complianceStatusRepo := complianceRepositories.NewComplianceStatusRepository(db)
	complianceStatusSvc := complianceServices.NewComplianceStatusService(complianceStatusRepo)
	complianceStatusController := controllers.NewComplianceStatusController(complianceStatusSvc)

	complianceRoutes.GET("/compliance-statuses", complianceStatusController.GetAllComplianceStatuses)
	complianceRoutes.GET("/compliance-statuses/:id", complianceStatusController.GetComplianceStatusDetails)

	//SLA Category
	SLACategoryRepo := complianceRepositories.NewSLACategoryRepository(db)
	SLACategorySvc := complianceServices.NewSLACategoryService(SLACategoryRepo, activityLogSvc)
	SLACategoryController := controllers.NewSLACategoryController(SLACategorySvc)

	complianceRoutes.GET("/sla-categories/:id", SLACategoryController.GetSLACategoryDetails)
	complianceRoutes.POST("/sla-categories", SLACategoryController.CreateCategory)
	complianceRoutes.GET("/sla-categories", SLACategoryController.GetAllSLACategories)
	complianceRoutes.PUT("/sla-categories/:id", SLACategoryController.UpdateCategory)
	complianceRoutes.POST("/sla-categories/delete", SLACategoryController.DeleteSLACategories)
	complianceRoutes.POST("/sla-categories/activate", SLACategoryController.ActivateSLACategories)
	complianceRoutes.POST("/sla-categories/deactivate", SLACategoryController.DeactivateSLACategories)

	//Compliance Register
	complianceRegisterRepo := complianceRepositories.NewComplianceRegisterRepository(db)
	reviewFrequencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	severityLevelRepo := masterRepositories.NewSeverityLevelRepository(db)
	complianceSvc := complianceServices.NewComplianceRegisterService(complianceRegisterRepo,
		activityLogSvc,
		complianceCategoryRepo, //@todo add org unit and user repo
		complianceStatusRepo,
		reviewFrequencyRepo,
		severityLevelRepo)
	complianceRegController := controllers.NewComplianceRegisterController(complianceSvc)

	complianceRoutes.GET("/compliance-register", complianceRegController.GetAllComplianceRegisters)
	complianceRoutes.GET("/compliance-register/:id", complianceRegController.GetComplianceRegisterDetails)
	complianceRoutes.POST("/compliance-register", complianceRegController.CreateComplianceRegister)
	complianceRoutes.PUT("/compliance-register/:id", complianceRegController.UpdateComplianceRegister)
	complianceRoutes.POST("/compliance-register/delete", complianceRegController.DeleteComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/archive", complianceRegController.ArchiveComplianceRegisters)
	complianceRoutes.PUT("/compliance-register/restore", complianceRegController.RestoreComplianceRegisters)
}
