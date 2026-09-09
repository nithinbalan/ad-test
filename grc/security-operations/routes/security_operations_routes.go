package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {
	securityOperationsRoutes := router.Group("/")

	securityOperationsRoutes.Use(middlewares.AuthMiddleware(app.DB))

	unauthenticatedRoutes := router.Group("/")

	db := app.DB

	csLogDataEntityTypeController := app.SecurityOperationsModule.CSLogDataEntityTypeController
	csLogDataReceptionMechanismController := app.SecurityOperationsModule.CSLogDataReceptionMechanismController
	csLogCommunicationChannelController := app.SecurityOperationsModule.CSLogCommunicationChannelController
	csLogSourceDirectoryController := app.SecurityOperationsModule.CSLogSourceDirectoryController
	csLogSocRoomAccessLogController := app.SecurityOperationsModule.CSLogSocRoomAccessLogController

	// CS Log Data Entity Types
	securityOperationsRoutes.GET("/security-operations/cs-log-data-entity-types", csLogDataEntityTypeController.GetAllCSLogDataEntityTypes)
	securityOperationsRoutes.GET("/security-operations/cs-log-data-entity-types/:id", csLogDataEntityTypeController.GetCSLogDataEntityTypeDetails)
	securityOperationsRoutes.POST("/security-operations/cs-log-data-entity-types", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.CreateCSLogDataEntityType)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-entity-types/:id", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.UpdateCSLogDataEntityType)
	securityOperationsRoutes.POST("/security-operations/cs-log-data-entity-types/delete", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.DeleteCSLogDataEntityTypes)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-entity-types/activate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.ActivateCSLogDataEntityTypes)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-entity-types/deactivate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataEntityTypeController.DeactivateCSLogDataEntityTypes)

	// CS Log Data Reception Mechanisms
	securityOperationsRoutes.GET("/security-operations/cs-log-data-reception-mechanisms", csLogDataReceptionMechanismController.GetAllCSLogDataReceptionMechanisms)
	securityOperationsRoutes.GET("/security-operations/cs-log-data-reception-mechanisms/:id", csLogDataReceptionMechanismController.GetCSLogDataReceptionMechanismDetails)
	securityOperationsRoutes.POST("/security-operations/cs-log-data-reception-mechanisms", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.CreateCSLogDataReceptionMechanism)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-reception-mechanisms/:id", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.UpdateCSLogDataReceptionMechanism)
	securityOperationsRoutes.POST("/security-operations/cs-log-data-reception-mechanisms/delete", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.DeleteCSLogDataReceptionMechanisms)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-reception-mechanisms/activate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.ActivateCSLogDataReceptionMechanisms)
	securityOperationsRoutes.PUT("/security-operations/cs-log-data-reception-mechanisms/deactivate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogDataReceptionMechanismController.DeactivateCSLogDataReceptionMechanisms)

	// CS Log Communication Channels
	securityOperationsRoutes.GET("/security-operations/cs-log-communication-channels", csLogCommunicationChannelController.GetAllCSLogCommunicationChannels)
	securityOperationsRoutes.GET("/security-operations/cs-log-communication-channels/:id", csLogCommunicationChannelController.GetCSLogCommunicationChannelDetails)
	securityOperationsRoutes.POST("/security-operations/cs-log-communication-channels", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.CreateCSLogCommunicationChannel)
	securityOperationsRoutes.PUT("/security-operations/cs-log-communication-channels/:id", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.UpdateCSLogCommunicationChannel)
	securityOperationsRoutes.POST("/security-operations/cs-log-communication-channels/delete", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.DeleteCSLogCommunicationChannels)
	securityOperationsRoutes.PUT("/security-operations/cs-log-communication-channels/activate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.ActivateCSLogCommunicationChannels)
	securityOperationsRoutes.PUT("/security-operations/cs-log-communication-channels/deactivate", middlewares.RolePermissionMiddleware(db, "security_operations_setting", "modify"), csLogCommunicationChannelController.DeactivateCSLogCommunicationChannels)

	// CS Log Source Directory
	securityOperationsRoutes.GET("/security-operations/source-directory", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "view"), csLogSourceDirectoryController.GetAllCSLogSourceDirectories)
	securityOperationsRoutes.GET("/security-operations/source-directory/export", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "view"), csLogSourceDirectoryController.GetCSLogSourceDirectoryExport)
	securityOperationsRoutes.GET("/security-operations/source-directory/:id", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "view"), csLogSourceDirectoryController.GetCSLogSourceDirectoryDetails)
	securityOperationsRoutes.POST("/security-operations/source-directory", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "modify"), csLogSourceDirectoryController.CreateCSLogSourceDirectory)
	securityOperationsRoutes.PUT("/security-operations/source-directory/:id", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "modify"), csLogSourceDirectoryController.UpdateCSLogSourceDirectory)
	securityOperationsRoutes.POST("/security-operations/source-directory/delete", middlewares.RolePermissionMiddleware(db, "security_operations_source_directory", "modify"), csLogSourceDirectoryController.DeleteCSLogSourceDirectories)

	// CS Log SOC Room Access Log
	securityOperationsRoutes.GET("/security-operations/soc-room-access-log", middlewares.RolePermissionMiddleware(db, "security_operations_soc_room_access_log", "view"), csLogSocRoomAccessLogController.GetAllCSLogSocRoomAccessLogs)
	securityOperationsRoutes.GET("/security-operations/soc-room-access-log/export", middlewares.RolePermissionMiddleware(db, "security_operations_soc_room_access_log", "view"), csLogSocRoomAccessLogController.GetCSLogSocRoomAccessLogExport)
	securityOperationsRoutes.GET("/security-operations/soc-room-access-log/:id", middlewares.RolePermissionMiddleware(db, "security_operations_soc_room_access_log", "view"), csLogSocRoomAccessLogController.GetCSLogSocRoomAccessLogDetails)
	securityOperationsRoutes.POST("/security-operations/soc-room-access-log", middlewares.RolePermissionMiddleware(db, "security_operations_soc_room_access_log", "modify"), csLogSocRoomAccessLogController.CreateCSLogSocRoomAccessLog)
	securityOperationsRoutes.PUT("/security-operations/soc-room-access-log/:id", middlewares.RolePermissionMiddleware(db, "security_operations_soc_room_access_log", "modify"), csLogSocRoomAccessLogController.UpdateCSLogSocRoomAccessLog)
	securityOperationsRoutes.POST("/security-operations/soc-room-access-log/delete", middlewares.RolePermissionMiddleware(db, "security_operations_soc_room_access_log", "modify"), csLogSocRoomAccessLogController.DeleteCSLogSocRoomAccessLogs)

	unauthenticatedRoutes.GET("/security-operations/soc-room-access-log/:id/document-preview", csLogSocRoomAccessLogController.GetCSLogSocRoomAccessLogDocumentPreview)
	unauthenticatedRoutes.GET("/security-operations/soc-room-access-log/:id/document-download", csLogSocRoomAccessLogController.DownloadCSLogSocRoomAccessLogDocument)
}
