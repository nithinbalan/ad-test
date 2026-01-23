package routes

import (
	"jethur/initializers"
	"jethur/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(app *initializers.App, router *gin.Engine) {

	consultationRequestRoutes := router.Group("/")
	consultationRequestRoutes.Use(middlewares.AuthMiddleware(app.DB))
	// unauthenticatedRoutes := router.Group("/")
	db := app.DB

	consultationRequestStatusController := app.ConsultationRequestModule.ConsultationRequestStatusController
	consultationRequestController := app.ConsultationRequestModule.ConsultationRequestController

	consultationRequestRoutes.GET("/consultation-request-statuses", consultationRequestStatusController.GetAllConsultationRequestStatuses)
	consultationRequestRoutes.GET("/consultation-request-statuses/:id", consultationRequestStatusController.GetConsultationRequestStatusDetails)

	consultationRequestRoutes.GET("/consultation-request", consultationRequestController.GetAllConsultationRequests)
	consultationRequestRoutes.GET("/consultation-request/:id", consultationRequestController.GetConsultationRequestDetails)
	consultationRequestRoutes.POST("/consultation-request", middlewares.RolePermissionMiddleware(db, "consultation", "modify"), consultationRequestController.CreateConsultationRequest)
	consultationRequestRoutes.PUT("/consultation-request/:id", middlewares.RolePermissionMiddleware(db, "consultation", "modify"), consultationRequestController.UpdateConsultationRequest)
	consultationRequestRoutes.POST("/consultation-request/delete", middlewares.RolePermissionMiddleware(db, "consultationt", "modify"), consultationRequestController.DeleteConsultationRequests)
	consultationRequestRoutes.PUT("/consultation-request/:id/responsible-user", middlewares.RolePermissionMiddleware(db, "consultation", "modify"), consultationRequestController.UpdateConsultationRequestResponsibleUser)
	consultationRequestRoutes.PUT("/consultation-request/:id/notes", middlewares.RolePermissionMiddleware(db, "consultation", "modify"), consultationRequestController.UpdateConsultationRequestNotes)
	consultationRequestRoutes.POST("/consultation-request/:id/comments", middlewares.RolePermissionMiddleware(db, "consultation", "modify"), consultationRequestController.CreateConsultationRequestComment)
	consultationRequestRoutes.GET("/consultation-request/:id/comments", middlewares.RolePermissionMiddleware(db, "consultation", "modify"), consultationRequestController.GetAllConsultationRequestComment)

}
