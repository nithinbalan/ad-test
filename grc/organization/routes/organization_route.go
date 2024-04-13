package routes

import (
	logRepositories "jethur/grc/logs/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/organization/controllers"
	organizationRepositories "jethur/grc/organization/repositories"
	organizationServices "jethur/grc/organization/services"
	settingRepositories "jethur/grc/settings/repositories"
	settingService "jethur/grc/settings/services"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB, router *gin.Engine) {

	organizationRoutes := router.Group("/")

	organizationRoutes.Use(middlewares.AuthMiddleware())

	activityLogRepo := logRepositories.NewActivityLogRepository(db)
	activityLogSvc := logServices.NewActivityLogService(activityLogRepo)

	roleRepo := settingRepositories.NewRoleRepository(db)
	roleSvc := settingService.NewRoleService(roleRepo, activityLogSvc)

	//Designations
	designationRepo := organizationRepositories.NewDesignationRepository(db)
	designationSvc := organizationServices.NewDesignationService(designationRepo, activityLogSvc)
	designationController := controllers.NewDesignationController(designationSvc)

	organizationRoutes.GET("/designations", designationController.GetAllDesignations)
	organizationRoutes.GET("/designations/:id", designationController.GetDesignationDetails)
	organizationRoutes.POST("/designations", designationController.CreateDesignation)
	organizationRoutes.PUT("/designations/:id", designationController.UpdateDesignation)
	organizationRoutes.POST("/designations/delete", designationController.DeleteDesignations)
	organizationRoutes.PUT("/designations/activate", designationController.ActivateDesignations)
	organizationRoutes.PUT("/designations/deactivate", designationController.DeactivateDesignations)

	//User Status
	userStatusRepo := organizationRepositories.NewUserStatusRepository(db)
	userStatusSvc := organizationServices.NewUserStatusService(userStatusRepo)
	userStatusController := controllers.NewUserStatusController(userStatusSvc)

	organizationRoutes.GET("/user-statuses", userStatusController.GetAllUserStatuses)
	organizationRoutes.GET("/user-statuses/:id", userStatusController.GetUserStatusDetails)

	//user
	userRepo := organizationRepositories.NewUserRepository(db)
	userSvc := organizationServices.NewUserService(userRepo, designationSvc, userStatusSvc, activityLogSvc, roleSvc)
	userController := controllers.NewUSerController(userSvc)

	organizationRoutes.GET("/users", userController.GetAllUsers)
	organizationRoutes.POST("/users", userController.CreateUser)
	organizationRoutes.PUT("/users/:id", userController.UpdateUser)
	organizationRoutes.POST("/users/delete", userController.DeleteUsers)

	//strategic profile
	strategicProfileRepo := organizationRepositories.NewStrategicProfileRepository(db)
	strategicProfileSvc := organizationServices.NewStrategicProfileService(strategicProfileRepo)
	strategicProfileController := controllers.NewStrategicProfileController(strategicProfileSvc)

	organizationRoutes.GET("/strategic-profiles", strategicProfileController.GetStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles", strategicProfileController.UpdateStrategicProfile)
	organizationRoutes.PUT("/strategic-profiles/vision", strategicProfileController.UpdateVision)
	organizationRoutes.PUT("/strategic-profiles/mission", strategicProfileController.UpdateMission)

}
