package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	securityOperationsRoutes := router.Group("/")

	securityOperationsRoutes.Use(middlewares.AuthMiddleware(app.DB))

	db := app.DB

	csLogDataEntityTypeController := app.SecurityOperationsModule.CSLogDataEntityTypeController
	csLogDataReceptionMechanismController := app.SecurityOperationsModule.CSLogDataReceptionMechanismController
	csLogCommunicationChannelController := app.SecurityOperationsModule.CSLogCommunicationChannelController
	csLogSourceDirectoryController := app.SecurityOperationsModule.CSLogSourceDirectoryController

	// CS Log Data Entity Types
	securityOperationsRoutes.GET("/security-operations/cs-log-data-entity-types", csLogDataEntityTypeController.GetAll)
	securityOperationsRoutes.GET("/security-operations/cs-log-data-entity-types/:id", csLogDataEntityTypeController.GetDetails)
	securityOperationsRoutes.POST("/security-operations/cs-log-data-entity-types", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.Create)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-entity-types/:id", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.Update)
	securityOperationsRoutes.POST("/security-operations/cs-log-data-entity-types/delete", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.Delete)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-entity-types/activate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.Activate)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-entity-types/deactivate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.Deactivate)

	// CS Log Data Reception Mechanisms
	securityOperationsRoutes.GET("/security-operations/cs-log-data-reception-mechanisms", csLogDataReceptionMechanismController.GetAll)
	securityOperationsRoutes.GET("/security-operations/cs-log-data-reception-mechanisms/:id", csLogDataReceptionMechanismController.GetDetails)
	securityOperationsRoutes.POST("/security-operations/cs-log-data-reception-mechanisms", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.Create)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-reception-mechanisms/:id", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.Update)
	securityOperationsRoutes.POST("/security-operations/cs-log-data-reception-mechanisms/delete", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.Delete)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-reception-mechanisms/activate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.Activate)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-reception-mechanisms/deactivate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.Deactivate)

	// CS Log Communication Channels
	securityOperationsRoutes.GET("/security-operations/cs-log-communication-channels", csLogCommunicationChannelController.GetAll)
	securityOperationsRoutes.GET("/security-operations/cs-log-communication-channels/:id", csLogCommunicationChannelController.GetDetails)
	securityOperationsRoutes.POST("/security-operations/cs-log-communication-channels", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.Create)
	securityOperationsRoutes.PUT("/security-operations/cs-log-communication-channels/:id", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.Update)
	securityOperationsRoutes.POST("/security-operations/cs-log-communication-channels/delete", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.Delete)
	securityOperationsRoutes.PUT("/security-operations/cs-log-communication-channels/activate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.Activate)
	securityOperationsRoutes.PUT("/security-operations/cs-log-communication-channels/deactivate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.Deactivate)

	// CS Log Source Directory
	securityOperationsRoutes.GET("/security-operations/source-directory", csLogSourceDirectoryController.GetAll)
	securityOperationsRoutes.GET("/security-operations/source-directory/:id", csLogSourceDirectoryController.GetDetails)
	securityOperationsRoutes.POST("/security-operations/source-directory", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "modify"), csLogSourceDirectoryController.Create)
	securityOperationsRoutes.PUT("/security-operations/source-directory/:id", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "modify"), csLogSourceDirectoryController.Update)
	securityOperationsRoutes.POST("/security-operations/source-directory/delete", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "modify"), csLogSourceDirectoryController.Delete)
}
