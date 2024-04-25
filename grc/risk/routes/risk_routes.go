package routes

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/risk/controllers"
	riskRepositories "jethur/grc/risk/repositories"
	riskServices "jethur/grc/risk/services"

	organizationRepositories "jethur/grc/organization/repositories"
	organizationServices "jethur/grc/organization/services"

	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"

	settingRepositories "jethur/grc/settings/repositories"
	settingService "jethur/grc/settings/services"

	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	riskRoutes := router.Group("/")

	riskRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	severityLevelRepo := masterRepositories.NewSeverityLevelRepository(db)
	severityLevelSvc := masterServices.NewSeverityLevelService(severityLevelRepo)

	kriCategoryRepo := masterRepositories.NewKRICategoryRepository(db)
	kriCategorySvc := masterServices.NewKRICategoryService(kriCategoryRepo, activityLogSvc)

	designationRepo := organizationRepositories.NewDesignationRepository(db)
	designationSvc := organizationServices.NewDesignationService(designationRepo, activityLogSvc)

	userStatusRepo := organizationRepositories.NewUserStatusRepository(db)
	userStatusSvc := organizationServices.NewUserStatusService(userStatusRepo)

	roleRepo := settingRepositories.NewRoleRepository(db)
	roleSvc := settingService.NewRoleService(roleRepo, activityLogSvc)

	userRepo := organizationRepositories.NewUserRepository(db)
	userSvc := organizationServices.NewUserService(userRepo, designationSvc, userStatusSvc, activityLogSvc, roleSvc)

	reviewFrequencyRepo := masterRepositories.NewReviewFrequencyRepository(db)
	reviewFrequencySvc := masterServices.NewReviewFrequencyService(reviewFrequencyRepo)

	riskCategoryRepo := riskRepositories.NewRiskCategoryRepository(db)
	riskCategorySvc := riskServices.NewRiskCategoryService(riskCategoryRepo, activityLogSvc)
	riskCategoryController := controllers.NewRiskCategoryController(riskCategorySvc)

	potentialRiskStatusRepo := riskRepositories.NewPotentialRiskStatusRepository(db)
	potentialRiskStatusSvc := riskServices.NewPotentialRiskStatusService(potentialRiskStatusRepo)
	potentialRiskStatusController := controllers.NewPotentialRiskStatusController(potentialRiskStatusSvc)

	riskLevelRepo := riskRepositories.NewRiskLevelRepository(db)
	riskLevelSvc := riskServices.NewRiskLevelService(riskLevelRepo)
	riskLevelController := controllers.NewRiskLevelController(riskLevelSvc)

	potentialRiskRepo := riskRepositories.NewPotentialRiskRepository(db)
	potentialRiskSvc := riskServices.NewPotentialRiskService(
		potentialRiskRepo,
		activityLogSvc,
		severityLevelSvc,
		riskCategorySvc,
		potentialRiskStatusSvc,
	)
	potentialRiskController := controllers.NewPotentialRiskController(potentialRiskSvc)

	potentialRiskStatusUpdateRepo := riskRepositories.NewPotentialRiskStatusUpdateRepository(db)
	potentialRiskStatusUpdateSvc := riskServices.NewPotentialRiskStatusUpdateService(potentialRiskStatusUpdateRepo, activityLogSvc, potentialRiskSvc, potentialRiskStatusSvc)
	potentialRiskStatusUpdateController := controllers.NewPotentialRiskStatusUpdateController(potentialRiskStatusUpdateSvc)

	kriRegisterRepo := riskRepositories.NewKRIRegisterRepository(db)
	kriRegisterSvc := riskServices.NewKRIRegisterService(
		kriRegisterRepo,
		activityLogSvc,
		kriCategorySvc,
		userSvc,
		riskLevelSvc,
		reviewFrequencySvc,
	)
	kriRegisterController := controllers.NewKRIRegisterController(kriRegisterSvc)

	riskAppetiteStatementRepo := riskRepositories.NewRiskAppetiteStatementRepository(db)
	riskAppetiteStatementSvc := riskServices.NewRiskAppetiteStatementService(
		riskAppetiteStatementRepo,
		activityLogSvc,
		riskCategorySvc,
	)
	riskAppetiteStatementController := controllers.NewRiskAppetiteStatementController(riskAppetiteStatementSvc)

	riskImpactCategoryRepo := riskRepositories.NewRiskImpactCategoryRepository(db)
	riskImpactCategorySvc := riskServices.NewRiskImpactCategoryService(riskImpactCategoryRepo, activityLogSvc)
	riskImpactCategoryController := controllers.NewRiskImpactCategoryController(riskImpactCategorySvc)

	riskMatrixRepo := riskRepositories.NewRiskMatrixRepository(db)
	riskMatrixSvc := riskServices.NewRiskMatrixService(riskMatrixRepo)

	riskMatrixImpactRepo := riskRepositories.NewRiskMatrixImpactRepository(db)
	riskMatrixImpactSvc := riskServices.NewRiskMatrixImpactService(
		riskMatrixImpactRepo,
		activityLogSvc,
		riskMatrixSvc,
	)
	riskMatrixImpactController := controllers.NewRiskMatrixImpactController(riskMatrixImpactSvc)

	riskMatrixLikelihoodRepo := riskRepositories.NewRiskMatrixLikelihoodRepository(db)
	riskMatrixLikelihoodSvc := riskServices.NewRiskMatrixLikelihoodService(
		riskMatrixLikelihoodRepo,
		activityLogSvc,
		riskMatrixSvc,
	)
	riskMatrixLikelihoodController := controllers.NewRiskMatrixLikelihoodController(riskMatrixLikelihoodSvc)

	riskMatrixRiskRatingRepo := riskRepositories.NewRiskMatrixRiskRatingRepository(db)
	riskMatrixRiskRatingSvc := riskServices.NewRiskMatrixRiskRatingService(
		riskMatrixRiskRatingRepo,
		activityLogSvc,
		riskMatrixSvc,
	)
	riskMatrixRiskRatingController := controllers.NewRiskMatrixRiskRatingController(riskMatrixRiskRatingSvc)

	riskRoutes.GET("/risk-categories", riskCategoryController.GetAllRiskCategories)
	riskRoutes.GET("/risk-categories/:id", riskCategoryController.GetRiskCategoryDetails)
	riskRoutes.POST("/risk-categories", riskCategoryController.CreateRiskCategory)
	riskRoutes.PUT("/risk-categories/:id", riskCategoryController.UpdateRiskCategory)
	riskRoutes.POST("/risk-categories/delete", riskCategoryController.DeleteRiskCategories)
	riskRoutes.PUT("/risk-categories/activate", riskCategoryController.ActivateRiskCategories)
	riskRoutes.PUT("/risk-categories/deactivate", riskCategoryController.DeactivateRiskCategories)

	riskRoutes.GET("/potential-risk-statuses", potentialRiskStatusController.GetAllPotentialRiskStatuses)
	riskRoutes.GET("/potential-risk-statuses/:id", potentialRiskStatusController.GetPotentialRiskStatusDetails)

	riskRoutes.GET("/potential-risks", potentialRiskController.GetAllPotentialRisks)
	riskRoutes.GET("/potential-risks/:id", potentialRiskController.GetPotentialRiskDetails)
	riskRoutes.POST("/potential-risks", potentialRiskController.CreatePotentialRisk)
	riskRoutes.PUT("/potential-risks/:id", potentialRiskController.UpdatePotentialRisk)
	riskRoutes.POST("/potential-risks/delete", potentialRiskController.DeletePotentialRisks)

	riskRoutes.GET("/risk-levels", riskLevelController.GetAllRiskLevels)
	riskRoutes.GET("/risk-levels/:id", riskLevelController.GetRiskLevelDetails)

	riskRoutes.GET("/potential-risks/:id/status-updates", potentialRiskStatusUpdateController.GetAllPotentialRiskStatusUpdates)
	riskRoutes.POST("/potential-risks/:id/status-updates", potentialRiskStatusUpdateController.CreatePotentialRiskStatusUpdate)

	riskRoutes.GET("/risk-appetite-statements", riskAppetiteStatementController.GetAllRiskAppetiteStatements)
	riskRoutes.GET("/risk-appetite-statements/:id", riskAppetiteStatementController.GetRiskAppetiteStatementDetails)
	riskRoutes.POST("/risk-appetite-statements", riskAppetiteStatementController.CreateRiskAppetiteStatement)
	riskRoutes.PUT("/risk-appetite-statements/:id", riskAppetiteStatementController.UpdateRiskAppetiteStatement)
	riskRoutes.POST("/risk-appetite-statements/delete", riskAppetiteStatementController.DeleteRiskAppetiteStatements)
	riskRoutes.PUT("/risk-appetite-statements/activate", riskAppetiteStatementController.ActivateRiskAppetiteStatements)
	riskRoutes.PUT("/risk-appetite-statements/deactivate", riskAppetiteStatementController.DeactivateRiskAppetiteStatements)

	riskRoutes.GET("/kri-register", kriRegisterController.GetAllKRIRegisters)
	riskRoutes.GET("/kri-register/:id", kriRegisterController.GetKRIRegisterDetails)
	riskRoutes.POST("/kri-register", kriRegisterController.CreateKRIRegister)
	riskRoutes.PUT("/kri-register/:id", kriRegisterController.UpdateKRIRegister)
	riskRoutes.POST("/kri-register/delete", kriRegisterController.DeleteKRIRegisters)

	riskRoutes.GET("/risk-matrix-impacts", riskMatrixImpactController.GetAllRiskMatrixImpacts)
	riskRoutes.GET("/risk-matrix-impacts/:id", riskMatrixImpactController.GetRiskMatrixImpactDetails)
	riskRoutes.POST("/risk-matrix-impacts", riskMatrixImpactController.CreateRiskMatrixImpact)
	riskRoutes.PUT("/risk-matrix-impacts/:id", riskMatrixImpactController.UpdateRiskMatrixImpact)
	riskRoutes.POST("/risk-matrix-impacts/delete", riskMatrixImpactController.DeleteRiskMatrixImpacts)
	riskRoutes.PUT("/risk-matrix-impacts/activate", riskMatrixImpactController.ActivateRiskMatrixImpact)
	riskRoutes.PUT("/risk-matrix-impacts/deactivate", riskMatrixImpactController.DeactivateRiskMatrixImpact)

	riskRoutes.GET("/risk-matrix-likelihood", riskMatrixLikelihoodController.GetAllRiskMatrixLikelihoods)
	riskRoutes.GET("/risk-matrix-likelihood/:id", riskMatrixLikelihoodController.GetRiskMatrixLikelihoodDetails)
	riskRoutes.POST("/risk-matrix-likelihood", riskMatrixLikelihoodController.CreateRiskMatrixLikelihood)
	riskRoutes.PUT("/risk-matrix-likelihood/:id", riskMatrixLikelihoodController.UpdateRiskMatrixLikelihood)
	riskRoutes.POST("/risk-matrix-likelihood/delete", riskMatrixLikelihoodController.DeleteRiskMatrixLikelihoods)
	riskRoutes.PUT("/risk-matrix-likelihood/activate", riskMatrixLikelihoodController.ActivateRiskMatrixLikelihood)
	riskRoutes.PUT("/risk-matrix-likelihood/deactivate", riskMatrixLikelihoodController.DeactivateRiskMatrixLikelihood)

	riskRoutes.GET("/risk-matrix-risk-ratings", riskMatrixRiskRatingController.GetAllRiskMatrixRiskRatings)
	riskRoutes.GET("/risk-matrix-risk-ratings/:id", riskMatrixRiskRatingController.GetRiskMatrixRiskRatingDetails)
	riskRoutes.POST("/risk-matrix-risk-ratings", riskMatrixRiskRatingController.CreateRiskMatrixRiskRating)
	riskRoutes.PUT("/risk-matrix-risk-ratings/:id", riskMatrixRiskRatingController.UpdateRiskMatrixRiskRating)
	riskRoutes.POST("/risk-matrix-risk-ratings/delete", riskMatrixRiskRatingController.DeleteRiskMatrixRiskRatings)
	riskRoutes.PUT("/risk-matrix-risk-ratings/activate", riskMatrixRiskRatingController.ActivateRiskMatrixRiskRating)
	riskRoutes.PUT("/risk-matrix-risk-ratings/deactivate", riskMatrixRiskRatingController.DeactivateRiskMatrixRiskRating)

	riskRoutes.GET("/risk-impact-categories", riskImpactCategoryController.GetAllRiskImpactCategories)
	riskRoutes.GET("/risk-impact-categories/:id", riskImpactCategoryController.GetRiskImpactCategoryDetails)
	riskRoutes.POST("/risk-impact-categories", riskImpactCategoryController.CreateRiskImpactCategory)
	riskRoutes.PUT("/risk-impact-categories/:id", riskImpactCategoryController.UpdateRiskImpactCategory)
	riskRoutes.POST("/risk-impact-categories/delete", riskImpactCategoryController.DeleteRiskImpactCategories)
	riskRoutes.PUT("/risk-impact-categories/activate", riskImpactCategoryController.ActivateRiskImpactCategories)
	riskRoutes.PUT("/risk-impact-categories/deactivate", riskImpactCategoryController.DeactivateRiskImpactCategories)

}
