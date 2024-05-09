package routes

import (
	complianceRepositories "jethur/grc/compliance/repositories"
	complianceServices "jethur/grc/compliance/services"
	"jethur/grc/incident/controllers"
	incidentRepositories "jethur/grc/incident/repositories"
	incidentServices "jethur/grc/incident/services"
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	organizationRepositories "jethur/grc/organization/repositories"
	organizationServices "jethur/grc/organization/services"
	settingRepositories "jethur/grc/settings/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	incidentRoutes := router.Group("/")

	incidentRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	incidentCategoryRepo := incidentRepositories.NewIncidentCategoryRepository(db)
	incidentCategorySvc := incidentServices.NewIncidentCategoryService(incidentCategoryRepo, activityLogSvc)
	incidentCategoryController := controllers.NewIncidentCategoryController(incidentCategorySvc)

	incidentSubCategoryRepo := incidentRepositories.NewIncidentSubCategoryRepository(db)
	incidentSubCategorySvc := incidentServices.NewIncidentSubCategoryService(incidentSubCategoryRepo, activityLogSvc, incidentCategorySvc)
	incidentSubCategoryController := controllers.NewIncidentSubCategoryController(incidentSubCategorySvc)

	incidentTypeRepo := incidentRepositories.NewIncidentTypeRepository(db)
	incidentTypeSvc := incidentServices.NewIncidentTypeService(incidentTypeRepo, activityLogSvc)
	incidentTypeController := controllers.NewIncidentTypeController(incidentTypeSvc)

	incidentStatusRepo := incidentRepositories.NewIncidentStatusRepository(db)
	incidentStatusSvc := incidentServices.NewIncidentStatusService(incidentStatusRepo)
	incidentStatusController := controllers.NewIncidentStatusController(incidentStatusSvc)

	//Control svc
	controlTypeRepo := organizationRepositories.NewControlTypeRepository(db)
	controlTypeSvc := organizationServices.NewControlTypeService(controlTypeRepo)

	controlImplementationStatusRepo := organizationRepositories.NewControlImplementationStatusRepository(db)
	controlImplementationStatusSvc := organizationServices.NewControlImplementationStatusService(controlImplementationStatusRepo)

	controlEfficiencyStatusRepo := organizationRepositories.NewControlEfficiencyStatusRepository(db)
	controlEfficiencyStatusSvc := organizationServices.NewControlEfficiencyStatusService(controlEfficiencyStatusRepo)

	designationRepo := organizationRepositories.NewDesignationRepository(db)
	designationSvc := organizationServices.NewDesignationService(designationRepo, activityLogSvc)

	userStatusRepo := organizationRepositories.NewUserStatusRepository(db)
	userStatusSvc := organizationServices.NewUserStatusService(userStatusRepo)

	roleRepo := settingRepositories.NewRoleRepository(db)
	roleSvc := settingServices.NewRoleService(roleRepo, activityLogSvc)

	userRepo := organizationRepositories.NewUserRepository(db)
	userSvc := organizationServices.NewUserService(userRepo, designationSvc, userStatusSvc, activityLogSvc, roleSvc)

	reviewFrequiencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	reviewFrequiencySvc := masterServices.NewReviewFrequencyService(reviewFrequiencyRepo)

	lifecycleStatusRepo := masterRepositories.NewLifecycleStatusRepository(db)
	lifecycleStatusSvc := masterServices.NewLifecycleStatusService(lifecycleStatusRepo)
	controlCategoryRepo := organizationRepositories.NewControlCategoryRepository(db)
	controlCategorySvc := organizationServices.NewControlCategoryService(controlCategoryRepo, activityLogSvc)

	controlRepo := organizationRepositories.NewControlRepository(db)
	controlSvc := organizationServices.NewControlService(
		controlRepo,
		activityLogSvc,
		controlTypeSvc,
		controlCategorySvc,
		controlImplementationStatusSvc,
		controlEfficiencyStatusSvc,
		userSvc,
		reviewFrequiencySvc,
		lifecycleStatusSvc)
	//Compliance svc
	complianceRequirementCategoryRepo := complianceRepositories.NewComplianceRequirementCategoryRepository(db)
	complianceRequirementCategorySvc := complianceServices.NewComplianceRequirementCategoryService(complianceRequirementCategoryRepo, activityLogSvc)

	severityLevelRepo := masterRepositories.NewSeverityLevelRepository(db)
	severityLevelSvc := masterServices.NewSeverityLevelService(severityLevelRepo)

	complianceStatusRepo := complianceRepositories.NewComplianceStatusRepository(db)
	complianceStatusSvc := complianceServices.NewComplianceStatusService(complianceStatusRepo)

	complianceCategoryRepo := complianceRepositories.NewComplianceCategoryRepository(db)
	complianceCategorySvc := complianceServices.NewComplianceCategoryService(complianceCategoryRepo, activityLogSvc)

	reviewFrequencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	reviewFrequencySvc := masterServices.NewReviewFrequencyService(reviewFrequencyRepo)

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

	complianceRequirementRepo := complianceRepositories.NewComplianceRequirementRepository(db)
	complianceRequirementSvc := complianceServices.NewComplianceRequirementService(
		complianceRequirementRepo,
		activityLogSvc,
		complianceRegisterSvc,
		complianceRequirementCategorySvc,
		complianceStatusSvc,
		severityLevelSvc,
	)

	incidentRegisterRepo := incidentRepositories.NewIncidentRegisterRepository(db)
	incidentRegisterSvc := incidentServices.NewIncidentRegisterService(incidentRegisterRepo, activityLogSvc, incidentCategorySvc, incidentSubCategorySvc, incidentTypeSvc, incidentStatusSvc, severityLevelSvc, controlSvc, complianceRequirementSvc)
	incidentRegisterController := controllers.NewIncidentRegisterController(incidentRegisterSvc)

	incidentInvestigationStatusRepo := incidentRepositories.NewIncidentInvestigationStatusRepository(db)
	incidentInvestigationStatusSvc := incidentServices.NewIncidentInvestigationStatusService(incidentInvestigationStatusRepo)
	incidentInvestigationStatusController := controllers.NewIncidentInvestigationStatusController(incidentInvestigationStatusSvc)

	incidentInvestigationRepo := incidentRepositories.NewIncidentInvestigationRepository(db)
	incidentInvestigationSvc := incidentServices.NewIncidentInvestigationService(
		incidentInvestigationRepo,
		activityLogSvc,
		userSvc,
		incidentInvestigationStatusSvc,
	)
	incidentInvestigationController := controllers.NewIncidentInvestigationController(incidentInvestigationSvc)

	//Incident Category
	incidentRoutes.GET("/incident-categories", incidentCategoryController.GetAllIncidentCategories)
	incidentRoutes.GET("/incident-categories/:id", incidentCategoryController.GetIncidentCategoryDetails)
	incidentRoutes.POST("/incident-categories", incidentCategoryController.CreateIncidentCategory)
	incidentRoutes.PUT("/incident-categories/:id", incidentCategoryController.UpdateIncidentCategory)
	incidentRoutes.POST("/incident-categories/delete", incidentCategoryController.DeleteIncidentCategories)
	incidentRoutes.PUT("/incident-categories/activate", incidentCategoryController.ActivateIncidentCategories)
	incidentRoutes.PUT("/incident-categories/deactivate", incidentCategoryController.DeactivateIncidentCategories)

	//Incident Sub Category
	incidentRoutes.GET("/incident-sub-categories", incidentSubCategoryController.GetAllIncidentSubCategories)
	incidentRoutes.GET("/incident-sub-categories/:id", incidentSubCategoryController.GetIncidentSubCategoryDetails)
	incidentRoutes.POST("/incident-sub-categories", incidentSubCategoryController.CreateIncidentSubCategory)
	incidentRoutes.PUT("/incident-sub-categories/:id", incidentSubCategoryController.UpdateIncidentSubCategory)
	incidentRoutes.POST("/incident-sub-categories/delete", incidentSubCategoryController.DeleteIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/activate", incidentSubCategoryController.ActivateIncidentSubCategories)
	incidentRoutes.PUT("/incident-sub-categories/deactivate", incidentSubCategoryController.DeactivateIncidentSubCategories)

	//Incident Type
	incidentRoutes.GET("/incident-types", incidentTypeController.GetAllIncidentTypes)
	incidentRoutes.GET("/incident-types/:id", incidentTypeController.GetIncidentTypeDetails)
	incidentRoutes.POST("/incident-types", incidentTypeController.CreateIncidentType)
	incidentRoutes.PUT("/incident-types/:id", incidentTypeController.UpdateIncidentType)
	incidentRoutes.POST("/incident-types/delete", incidentTypeController.DeleteIncidentTypes)
	incidentRoutes.PUT("/incident-types/activate", incidentTypeController.ActivateIncidentTypes)
	incidentRoutes.PUT("/incident-types/deactivate", incidentTypeController.DeactivateIncidentTypes)

	//Incident Status
	incidentRoutes.GET("/incident-statuses", incidentStatusController.GetAllIncidentStatuses)
	incidentRoutes.GET("/incident-statuses/:id", incidentStatusController.GetIncidentStatusDetails)

	//Incident Register
	incidentRoutes.GET("/incident-register", incidentRegisterController.GetAllIncidentRegisters)
	incidentRoutes.GET("/incident-register/:id", incidentRegisterController.GetIncidentRegisterDetails)
	incidentRoutes.POST("/incident-register", incidentRegisterController.CreateIncidentRegister)
	incidentRoutes.PUT("/incident-register/:id", incidentRegisterController.UpdateIncidentRegister)
	incidentRoutes.POST("/incident-register/delete", incidentRegisterController.DeleteIncidentRegisters)

	//Incident Investigation Status
	incidentRoutes.GET("/incident-investigation-statuses", incidentInvestigationStatusController.GetAllIncidentInvestigationStatuses)
	incidentRoutes.GET("/incident-investigation-statuses/:id", incidentInvestigationStatusController.GetIncidentInvestigationStatusDetails)

	//Incident Investigation
	incidentRoutes.GET("/incident-investigations", incidentInvestigationController.GetAllIncidentInvestigations)
	incidentRoutes.GET("/incident-investigations/:id", incidentInvestigationController.GetIncidentInvestigationDetails)
	incidentRoutes.POST("/incident-investigations", incidentInvestigationController.CreateIncidentInvestigation)
	incidentRoutes.PUT("/incident-investigations/:id", incidentInvestigationController.UpdateIncidentInvestigation)
	incidentRoutes.POST("/incident-investigations/delete", incidentInvestigationController.DeleteIncidentInvestigations)

}
